package healthcheckmanager

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"uapregistry/config"
	"uapregistry/leadermanager"
	"uapregistry/logger"
	"uapregistry/servicemanager"
	sc "uapregistry/storage/consulagent/cache"
	"uapregistry/types"
	"uapregistry/utils"
)

type HealthCheckInfo struct {
	ID                    string
	Name                  string
	ServiceIP             string
	ServicePort           string
	CheckType             string
	CheckInterval         time.Duration
	CheckTimeOut          time.Duration
	CheckUnhealthyTimeout time.Duration
	CheckHTTPMethod       string
	CheckHTTPURL          string
	CheckExpectedCodes    string

	RenewalInterval         time.Duration
	RenewalDeleteTimeout    time.Duration
	RenewalUnhealthyTimeout time.Duration
}

type HealthCheckInfoStore struct {
	sync.Mutex
	id     string
	HCInfo *HealthCheckInfo

	log logger.Logger

	httpClient     *http.Client
	dialer         *net.Dialer
	unhealthyTimer *time.Timer
	deleteTimer    *time.Timer

	healthCheckRetryEnabled bool
	isHealthy               bool
	stop                    bool
	stopCh                  chan struct{}
}

type HealthCheckInfoCache struct {
	sync.RWMutex
	stopped bool
	stopCh  chan struct{}
	index   uint64
	Data    map[string]*HealthCheckInfoStore
	log     logger.Logger
}

const (
	interval        = 5
	concurrency     = 500
	healthyStatus   = "healthy"
	unhealthyStatus = "unhealthy"
)

var (
	cache = &HealthCheckInfoCache{
		stopped: true,
		stopCh:  make(chan struct{}),
		index:   0,
		Data:    make(map[string]*HealthCheckInfoStore, 0),
	}
	sem = make(chan struct{}, concurrency)
)

func GetCache() *HealthCheckInfoCache {
	return cache
}

// 如果开启了健康检查重试，则选择一个uapregistry实例作为leader，leader做健康检查,非leader不做健康检查
// 如果不进行健康检查重试，则所有的uapregistry实例都做健康检查
func StartHealthCheckManager() {
	cache.log = logger.GetLogger()
	lm, err := leadermanager.GetLeaderManager()
	if err != nil {
		logger.GetLogger().Errorf("failed to get GetLeaderManager()")
		return
	}

	if lm.IsLeader() {
		logger.GetLogger().Infof("starting health check manager")
		cache.startHealthCheck()
	}

	ch := make(chan bool)
	lm.WaitLeaderCh(ch)

	for isLeader := range ch {
		if isLeader {
			logger.GetLogger().Infof("health check leader status changed,now isLeader is %t,starting health check manager", isLeader)
			cache.startHealthCheck()
		} else {
			logger.GetLogger().Infof("health check leader status changed,now isLeader is %t,stopping health check manager", isLeader)
			cache.stopHealthCheck()
		}
	}
}

func (c *HealthCheckInfoCache) startHealthCheck() {
	if c.stopped {
		c.stopped = false
		cache.initData()
		go cache.watchServices()
	}
}

func (c *HealthCheckInfoCache) stopHealthCheck() {
	if c.stopped {
		return
	}

	if c.stopCh != nil {
		close(c.stopCh)
	}

	for key := range c.Data {
		c.delete(key)
	}

	c.stopped = true
}

type BlockingQueryResp struct {
	sus       []*types.Service
	lastIndex uint64
}

func blockingQueryAllServices(waitTime time.Duration, waitIndex uint64, respCh chan BlockingQueryResp) {
	sus, index := sc.GetServicesCache().BlockingQueryAllServices(waitTime, waitIndex)

	timeoutCh := time.After(5 * time.Second)
	select {
	case respCh <- BlockingQueryResp{sus, index}:
	case <-timeoutCh:
		logger.GetLogger().Warn("timeout occurred while notify all services in health check manager")
	}
}

func (c *HealthCheckInfoCache) watchServices() {
	var (
		waitTime  = 10 * time.Minute
		waitIndex = c.getLastIndex()
	)

	c.log.Infof("start to watch all services")
	c.stopCh = make(chan struct{})

	for {
		respCh := make(chan BlockingQueryResp)
		go blockingQueryAllServices(waitTime, waitIndex, respCh)
		select {
		case <-c.stopCh:
			logger.GetLogger().Infof("Recieve message on the stopCh, return from watch all services")
			return
		case resp := <-respCh:
			if resp.lastIndex == waitIndex {
				continue
			}

			logger.GetLogger().Infof("services changed in health check manager")
			dataNew := services2healthCheckInfoStores(resp.sus)
			deleteKeys := getDeletedHealthCheckInfoStoresKeys(dataNew, c.Data)
			updatedHCIS := getUpdatedHealthCheckInfoStores(dataNew, c.Data)

			for _, key := range deleteKeys {
				c.delete(key)
			}
			for key, s := range updatedHCIS {
				c.update(key, s)
			}

			waitIndex = resp.lastIndex

			// 防止服务频繁变化，频繁处理
			time.Sleep(time.Second)
		}
	}
}

func (c *HealthCheckInfoCache) initData() {
	c.Lock()
	defer c.Unlock()

	c.log.Infof("start to init health check info cache")
	c.Data = make(map[string]*HealthCheckInfoStore, 0)

	allServices, lastIndex := sc.GetServicesCache().BlockingQueryAllServices(0, 0)

	c.index = lastIndex
	c.Data = services2healthCheckInfoStores(allServices)
	for _, s := range c.Data {
		go s.Start()
	}
	c.log.Infof("init health check info cache successfully")
}

func (c *HealthCheckInfoCache) getLastIndex() uint64 {
	c.RLock()
	defer c.RUnlock()

	return c.index
}

func (c *HealthCheckInfoCache) delete(key string) {
	if s, ok := c.Data[key]; ok {
		s.Stop()
		delete(c.Data, key)
		c.log.Infof("delete %s successfully", key)
	}
}

func (c *HealthCheckInfoCache) update(key string, s *HealthCheckInfoStore) {
	if s, ok := c.Data[key]; ok {
		s.Stop()
	}
	c.Data[key] = s
	go c.Data[key].Start()
	c.log.Infof("update %s successfully", key)
}

func (c *HealthCheckInfoCache) UpdateServiceTTL(id string) (statusCode int, err error) {
	c.Lock()
	defer c.Unlock()

	s, ok := c.Data[id]
	if !ok {
		return http.StatusNotFound, errors.New("microservice not found")
	}
	if s.HCInfo.CheckType != "TTL" {
		return http.StatusBadRequest, errors.New("checkType is not TTL")
	}

	err = servicemanager.NewServiceManager().UpdateServiceStatus(id, healthyStatus)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	s.updateStatus(true)
	s.unhealthyTimer.Reset(s.HCInfo.RenewalUnhealthyTimeout)
	s.deleteTimer.Reset(s.HCInfo.RenewalDeleteTimeout)
	return
}

func (s *HealthCheckInfoStore) Stop() {
	s.Lock()
	defer s.Unlock()

	if !s.stop {
		s.log.Infof("stop health check(id:%s)", s.id)
		s.stop = true
		close(s.stopCh)
	}
}

func (s *HealthCheckInfoStore) Start() {
	s.Lock()
	defer s.Unlock()

	if s.HCInfo.CheckType == "" {
		return
	}

	s.log.Infof("start health check(id:%s)", s.id)

	s.stop = false
	s.stopCh = make(chan struct{})

	checkType := strings.ToUpper(s.HCInfo.CheckType)
	switch checkType {
	case "HTTP", "HTTPS":
		s.createHTTPClient()
		s.unhealthyTimer = time.NewTimer(s.HCInfo.CheckUnhealthyTimeout)
		go s.runHTTPHealthCheck()
	case "TTL":
		s.unhealthyTimer = time.NewTimer(s.HCInfo.RenewalUnhealthyTimeout)
		s.deleteTimer = time.NewTimer(s.HCInfo.RenewalDeleteTimeout)
		go s.runTTLHealthCheck()
	case "TCP":
		s.initSockerDialer()
		s.unhealthyTimer = time.NewTimer(s.HCInfo.CheckUnhealthyTimeout)
		go s.runTCPHealthCheck()
	case "UDP-CONNECT":
		s.unhealthyTimer = time.NewTimer(s.HCInfo.CheckUnhealthyTimeout)
		go s.runUDPHealthCheck()
	}
}

func (s *HealthCheckInfoStore) doHealthCheckForOnce(c chan *types.HealthCheckResult) {
	var (
		err    error
		result = &types.HealthCheckResult{}
	)

	checkType := strings.ToUpper(s.HCInfo.CheckType)
	switch checkType {
	case "HTTP", "HTTPS":
		s.createHTTPClient()
		err = s.doHTTPRequest()
	case "TCP":
		s.initSockerDialer()
		err = s.doTCPHealthCheck()
	case "UDP-CONNECT":
		err = s.doUDPHealthCheck()
	default:
		err = fmt.Errorf("unsupported Health Check Type:%s", checkType)
	}

	if err != nil {
		result.Status = unhealthyStatus
	} else {
		result.Status = healthyStatus
	}

	timeoutCh := time.After(time.Second)
	select {
	case c <- result:
		s.log.Debugf("[%s]send health check result successfully", s.id)
	case <-timeoutCh:
		s.log.Warnf("[%s]send health check result timeout", s.id)
	}
}

func DoServiceHealthCheckForOnce(svc *types.Service) *types.HealthCheckResult {
	var (
		timeout  time.Duration
		hcResult = &types.HealthCheckResult{Status: unhealthyStatus}
		c        = make(chan *types.HealthCheckResult)
	)

	if svc == nil {
		return hcResult
	}

	s := service2HealthCheckInfoStore(svc)
	if s == nil {
		return hcResult
	}

	go s.doHealthCheckForOnce(c)
	timeout = s.HCInfo.CheckTimeOut

	timeout = timeout / 10
	timeoutCh := time.After(timeout)
	select {
	case hcResult = <-c:
		return hcResult
	case <-timeoutCh:
		logger.GetLogger().Warnf("[%s]timeout in receive all instance health check result", svc.ID)
		return hcResult
	}
}

func (s *HealthCheckInfoStore) runTTLHealthCheck() {
	for {
		select {
		case <-s.unhealthyTimer.C:
			s.log.Warnf("service(%s) missed TTL,set it as unhealthy", s.id)
			s.updateStatus(false)
		case <-s.deleteTimer.C:
			s.log.Warnf("service(%s) missed TTL,delete it", s.id)
			s.deleteService()
		case <-s.stopCh:
			return
		}
	}
}

func (s *HealthCheckInfoStore) runTCPHealthCheck() {
	next := time.After(randomStaggerTime(s.HCInfo.CheckInterval))
	for {
		select {
		case <-next:
			next = time.After(s.HCInfo.CheckInterval)
			s.tcpHealthCheck()
		case <-s.unhealthyTimer.C:
			s.log.Warnf("service(name:%s,id:%s) unhealthy timeout,set it as unhealthy", s.HCInfo.Name, s.id)
			s.updateStatus(false)
		case <-s.stopCh:
			s.log.Infof("stop health check(%s) successfully", s.id)
			return
		}
	}
}

func (s *HealthCheckInfoStore) runUDPHealthCheck() {
	next := time.After(randomStaggerTime(s.HCInfo.CheckInterval))
	for {
		select {
		case <-next:
			next = time.After(s.HCInfo.CheckInterval)
			s.udpHealthCheck()
		case <-s.unhealthyTimer.C:
			s.log.Warnf("service(name:%s,id:%s) unhealthy timeout,set it as unhealthy", s.HCInfo.Name, s.id)
			s.updateStatus(false)
		case <-s.stopCh:
			s.log.Infof("stop health check(%s) successfully", s.id)
			return
		}
	}
}

func (s *HealthCheckInfoStore) udpHealthCheck() {
	err := s.doUDPHealthCheck()
	if err != nil {
		s.log.Warnf("health check(id:%s,serviceName:%s) failed: %v,recheck service status from other members", s.id, s.HCInfo.Name, err)
		hcr := s.recheckService()
		if hcr.Status == unhealthyStatus {
			s.log.Warnf("health check(id:%s,serviceName:%s) failed from other members: %v", s.id, s.HCInfo.Name, err)
			return
		}
	}

	s.dealSuccessCheck()
}

func (s *HealthCheckInfoStore) doUDPHealthCheck() error {
	ip := net.ParseIP(s.HCInfo.ServiceIP)
	port, _ := strconv.Atoi(s.HCInfo.ServicePort)
	dstAddr := &net.UDPAddr{IP: ip, Port: port}

	conn, err := net.DialUDP("udp", nil, dstAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(s.HCInfo.CheckTimeOut))
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		return err
	}

	data := make([]byte, 0)
	_, _, err = conn.ReadFromUDP(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *HealthCheckInfoStore) tcpHealthCheck() {
	err := s.doTCPHealthCheck()
	if err != nil {
		s.log.Warnf("health check(id:%s,serviceName:%s) failed: %v,recheck service status from other members", s.id, s.HCInfo.Name, err)
		hcr := s.recheckService()
		if hcr.Status == unhealthyStatus {
			s.log.Warnf("health check(id:%s,serviceName:%s) failed from other members: %v", s.id, s.HCInfo.Name, err)
			return
		}
	}

	s.dealSuccessCheck()
}

func (s *HealthCheckInfoStore) doTCPHealthCheck() error {
	var (
		network = "tcp"
		addr    = s.HCInfo.ServiceIP + `:` + s.HCInfo.ServicePort
	)

	if strings.Contains(s.HCInfo.ServiceIP, ":") {
		addr = `[` + s.HCInfo.ServiceIP + `]:` + s.HCInfo.ServicePort
	}

	s.log.Debugf("network:%s,addr:%s", network, addr)
	conn, err := s.dialer.Dial(network, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func (s *HealthCheckInfoStore) initSockerDialer() {
	s.dialer = &net.Dialer{
		Timeout:   s.HCInfo.CheckTimeOut,
		DualStack: true,
	}
}

func (s *HealthCheckInfoStore) createHTTPClient() {
	tr := &http.Transport{
		/* #nosec */
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: config.GetInsecureSkipVerify()},
		DisableKeepAlives: true,
	}
	s.httpClient = &http.Client{
		Transport: tr,
		Timeout:   s.HCInfo.CheckTimeOut,
	}
}

func (s *HealthCheckInfoStore) runHTTPHealthCheck() {
	next := time.After(randomStaggerTime(s.HCInfo.CheckInterval))
	for {
		select {
		case <-next:
			next = time.After(s.HCInfo.CheckInterval)
			s.httpHealthCheck()
		case <-s.unhealthyTimer.C:
			s.log.Warnf("service(name:%s,id:%s) unhealthy timeout,set it as unhealthy", s.HCInfo.Name, s.id)
			s.updateStatus(false)
		case <-s.stopCh:
			s.log.Infof("stop health check(%s) successfully", s.id)
			return
		}
	}
}

func (s *HealthCheckInfoStore) httpHealthCheck() {
	err := s.doHTTPRequest()
	if err != nil {
		s.log.Warnf("health check(id:%s,serviceName:%s) failed: %v,recheck service status from other members", s.id, s.HCInfo.Name, err)
		hcr := s.recheckService()
		if hcr.Status == unhealthyStatus {
			s.log.Warnf("health check(id:%s,serviceName:%s) failed from other members: %v", s.id, s.HCInfo.Name, err)
			return
		}
	}

	s.dealSuccessCheck()
}

func (s *HealthCheckInfoStore) recheckService() *types.HealthCheckResult {
	defaultHcr := &types.HealthCheckResult{
		Status: unhealthyStatus,
	}

	if !s.healthCheckRetryEnabled {
		return defaultHcr
	}

	addresses := getUapregistryAddress()

	c := make(chan *types.HealthCheckResult)
	for _, addr := range addresses {
		go getInstanceHealthCheckResult(addr, s.id, s.HCInfo.CheckTimeOut/10, c)
	}

	timeoutCh := time.After(s.HCInfo.CheckTimeOut / 10)
xx:
	for range addresses {
		select {
		case result := <-c:
			if result != nil && result.Status == healthyStatus {
				return result
			}
		case <-timeoutCh:
			break xx
		}
	}

	return defaultHcr
}

func (s *HealthCheckInfoStore) doHTTPRequest() error {
	ip := s.HCInfo.ServiceIP
	if strings.Contains(ip, ":") {
		ip = "[" + ip + "]"
	}
	url := s.HCInfo.CheckType + `://` + ip + `:` + s.HCInfo.ServicePort + s.HCInfo.CheckHTTPURL
	req, err := http.NewRequest(strings.ToUpper(s.HCInfo.CheckHTTPMethod), url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", utils.ModuleName)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	ok := isExpectedStatusCode(statusCode, s.HCInfo.CheckExpectedCodes)
	if ok {
		return nil
	}

	return fmt.Errorf("(url:%s)statu code(%d) is not in CheckExpectedCodes(%s)", url, statusCode, s.HCInfo.CheckExpectedCodes)
}

func (s *HealthCheckInfoStore) dealSuccessCheck() {
	s.unhealthyTimer.Reset(s.HCInfo.CheckUnhealthyTimeout)

	if s.isHealthy {
		return
	}

	s.updateStatus(true)
	s.log.Warnf("service(name:%s,id:%s) health check successfully,reset unhealthy timeout to %s", s.HCInfo.Name, s.id, s.HCInfo.CheckUnhealthyTimeout.String())
}

func (s *HealthCheckInfoStore) updateStatus(isHealthy bool) {
	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	var status string
	if isHealthy {
		status = healthyStatus
	} else {
		status = unhealthyStatus
	}

	err := servicemanager.NewServiceManager().UpdateServiceStatus(s.id, status)
	if err != nil {
		s.log.Errorf("failed to update service(%s) status:%v", s.id, err)
		return
	}

	s.isHealthy = isHealthy
}

func (s *HealthCheckInfoStore) deleteService() {
	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	err := servicemanager.NewServiceManager().DeleteServiceByID(s.id)
	if err != nil {
		s.log.Errorf("failed to delete service(id:%s) status:%v", s.id, err)
		return
	}
}

func service2HealthCheckInfoStore(svc *types.Service) *HealthCheckInfoStore {
	if svc.EphemeralCheck == nil && svc.PersistentCheck == nil {
		return nil
	}

	hi := HealthCheckInfo{}

	if svc.EphemeralCheck != nil {
		hi = HealthCheckInfo{
			ID:                      svc.ID,
			Name:                    svc.Name,
			ServiceIP:               svc.Host,
			ServicePort:             strconv.Itoa(svc.Port),
			CheckType:               svc.EphemeralCheck.CheckType,
			RenewalInterval:         10 * time.Second,
			RenewalUnhealthyTimeout: 30 * time.Second,
			RenewalDeleteTimeout:    60 * time.Second,
		}
		if svc.EphemeralCheck.RenewalInterval != "" {
			t, _ := time.ParseDuration(svc.EphemeralCheck.RenewalInterval)
			hi.RenewalInterval = t
		}
		if svc.EphemeralCheck.RenewalDeleteTimeout != "" {
			t, _ := time.ParseDuration(svc.EphemeralCheck.RenewalDeleteTimeout)
			hi.RenewalDeleteTimeout = t
		}
		if svc.EphemeralCheck.RenewalUnhealthyTimeout != "" {
			t, _ := time.ParseDuration(svc.EphemeralCheck.RenewalUnhealthyTimeout)
			hi.RenewalUnhealthyTimeout = t
		}
	}
	if svc.PersistentCheck != nil {
		hi = HealthCheckInfo{
			ID:                 svc.ID,
			Name:               svc.Name,
			ServiceIP:          svc.Host,
			ServicePort:        strconv.Itoa(svc.Port),
			CheckType:          svc.PersistentCheck.CheckType,
			CheckInterval:      15 * time.Second,
			CheckTimeOut:       5 * time.Second,
			CheckHTTPMethod:    svc.PersistentCheck.CheckHTTPMethod,
			CheckHTTPURL:       svc.PersistentCheck.CheckHTTPURL,
			CheckExpectedCodes: "200",
		}
		if svc.PersistentCheck.CheckInterval != "" {
			t, _ := time.ParseDuration(svc.PersistentCheck.CheckInterval)
			hi.CheckInterval = t
		}
		if svc.PersistentCheck.CheckTimeout != "" {
			t, _ := time.ParseDuration(svc.PersistentCheck.CheckTimeout)
			hi.CheckTimeOut = t
		}
		if svc.PersistentCheck.CheckUnhealthyTimeout != "" {
			t, _ := time.ParseDuration(svc.PersistentCheck.CheckUnhealthyTimeout)
			hi.CheckUnhealthyTimeout = t
		}
	}

	hs := HealthCheckInfoStore{
		id:                      hi.ID,
		HCInfo:                  &hi,
		log:                     logger.GetLogger(),
		isHealthy:               true,
		healthCheckRetryEnabled: true,
		stop:                    true,
		stopCh:                  make(chan struct{}),
	}

	if svc.HealthStatus == unhealthyStatus {
		hs.isHealthy = false
	}

	return &hs
}

func services2healthCheckInfoStores(svcs []*types.Service) map[string]*HealthCheckInfoStore {
	ss := make(map[string]*HealthCheckInfoStore, 0)

	for _, svc := range svcs {
		ssTmp := service2HealthCheckInfoStore(svc)
		if ssTmp == nil {
			continue
		}
		ss[ssTmp.id] = ssTmp
	}

	return ss
}

func getUpdatedHealthCheckInfoStores(dataNew, dataOld map[string]*HealthCheckInfoStore) map[string]*HealthCheckInfoStore {
	updatedHCIS := make(map[string]*HealthCheckInfoStore, 0)
	for keyNew, sNew := range dataNew {
		sOld, ok := dataOld[keyNew]
		if !ok {
			updatedHCIS[keyNew] = sNew
			continue
		}
		if sNew.isHealthy != sOld.isHealthy || !reflect.DeepEqual(sNew.HCInfo, sOld.HCInfo) {
			updatedHCIS[keyNew] = sNew
			continue
		}
	}

	return updatedHCIS
}

func getDeletedHealthCheckInfoStoresKeys(dataNew, dataOld map[string]*HealthCheckInfoStore) []string {
	delKeys := make([]string, 0)

	for oldKey := range dataOld {
		if _, ok := dataNew[oldKey]; !ok {
			delKeys = append(delKeys, oldKey)
		}
	}
	return delKeys
}

func isExpectedStatusCode(statusCode int, expectedCodes string) bool {
	statusCodeStr := strconv.Itoa(statusCode)
	if statusCodeStr == expectedCodes {
		return true
	}
	if strings.Contains(expectedCodes, ",") {
		for _, ec := range strings.Split(expectedCodes, ",") {
			if statusCodeStr == ec {
				return true
			}
		}
	}

	if strings.Contains(expectedCodes, "-") {
		ecs := strings.SplitN(expectedCodes, "-", 2)
		if len(ecs) != 2 {
			return false
		}

		scMin, _ := strconv.Atoi(ecs[0])
		scMax, _ := strconv.Atoi(ecs[1])
		if statusCode >= scMin && statusCode <= scMax {
			return true
		}
	}

	return false
}

func randomStaggerTime(intv time.Duration) time.Duration {
	if intv == 0 {
		return 0
	}
	/* #nosec */
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}

func getUapregistryAddress() []string {
	serviceName := "uapregistry-health-check"
	svcs := sc.GetServicesCache().GetServicesByServiceName(serviceName)
	addresses := make([]string, 0)
	for _, svc := range svcs {
		if svc.HealthStatus == healthyStatus {
			// 除本节点外的其他节点的uapregistry做健康检查重试
			if svc.Host == utils.GetNodeIP() {
				continue
			}
			var address string
			if strings.Contains(svc.Host, ":") {
				address = "[" + svc.Host + "]:" + strconv.Itoa(svc.Port)
			} else {
				address = svc.Host + ":" + strconv.Itoa(svc.Port)
			}
			addresses = append(addresses, address)
		}
	}

	return addresses
}

func getInstanceHealthCheckResult(address, id string, checkTimeOut time.Duration, ch chan *types.HealthCheckResult) {
	url := "http://" + address + "/services/" + id + "/healthcheck"
	var hcResult *types.HealthCheckResult
	respBody, err := utils.HTTPGetWithTime(url, checkTimeOut)
	if err != nil {
		logger.GetLogger().Warnf("failed to getServiceHealthCheckResult:%v", err)
	}

	logger.GetLogger().Warnf("url:%s,respBody:%s", url, string(respBody))

	if len(respBody) != 0 {
		err = json.Unmarshal(respBody, &hcResult)
		if err != nil {
			logger.GetLogger().Warnf("failed to Unmarshal:%v", err)
		}
	}

	timeoutCh := time.After(1 * time.Second)
	select {
	case ch <- hcResult:
	case <-timeoutCh:
		logger.GetLogger().Debugf("timeout occurred while send health check result")
	}
}
