package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"uapregistry/logger"
	t "uapregistry/types"

	"github.com/xeipuuv/gojsonschema"
)

func CheckAndModifyService(svc *t.Service) error {
	if svc.Name == "" {
		return errors.New("service name is empty")
	}
	svc.Protocol = strings.ToLower(svc.Protocol)
	if !isValidProtocol(svc.Protocol) {
		return fmt.Errorf("protocol %s is not supported", svc.Protocol)
	}
	if !CheckAddress(svc.Host) {
		return fmt.Errorf("host %s is not valid", svc.Host)
	}
	if !CheckPort(svc.Port) {
		return fmt.Errorf("port %d is not valid", svc.Port)
	}
	if err := checkServiceRetries(int64(svc.Retries)); err != nil {
		return err
	}

	svc.AgentProtocol = strings.ToLower(svc.AgentProtocol)
	if err := checkAgentProtocol(svc.AgentProtocol); err != nil {
		return err
	}
	if svc.Ephemeral {

	}
	if err := CheckEphemeralCheck(svc.EphemeralCheck); err != nil {
		return err
	}
	if err := checkPersistentCheck(svc.PersistentCheck); err != nil {
		return err
	}

	// if svc.AgentInfo == nil && svc.AgentInfoUrl == "" {
	// 	return errors.New("agent_info and agent_info_url are both empty")
	// }

	if svc.AgentInfo == nil && svc.AgentInfoUrl != "" {
		agentInfo, err := getAgentInfo(svc.AgentProtocol, svc.AgentInfoUrl)
		if err != nil {
			return fmt.Errorf("failed to get agent info: %v", err)
		}
		svc.AgentInfo = agentInfo
	}

	if err := checkAgentInfo(svc); err != nil {
		return err
	}

	return nil
}

func checkAgentInfo(svc *t.Service) error {
	if svc.AgentInfo == nil {
		return nil
	}

	if svc.AgentProtocol == "a2a" {
		return checkA2AAgentCard(svc.AgentInfo.A2AAgentCard)
	}

	return nil
}

func getAgentInfo(agentProtocol, url string) (*t.AgentInfo, error) {
	resp, err := HTTPGet(url)
	if err != nil {
		return nil, err
	}

	logger.GetLogger().Debugf("get agent info from %s: %s", url, string(resp))

	var ac map[string]interface{}
	err = json.Unmarshal(resp, &ac)
	if err != nil {
		return nil, err
	}

	switch agentProtocol {
	case "a2a":
		return &t.AgentInfo{A2AAgentCard: ac}, nil
	case "acp":
		return &t.AgentInfo{AcpAgentManifest: ac}, nil
	case "mcp":
		return &t.AgentInfo{McpServer: ac}, nil
	}

	return nil, nil
}

func checkA2AAgentCard(ac map[string]interface{}) error {
	if ac == nil {
		return errors.New("a2a_agent_card is nil")
	}
	wd, err := os.Getwd()
	if err != nil {
		logger.GetLogger().Errorf("get current working directory failed: %v", err)
		return err
	}

	fileName := filepath.Join(wd, "config", "schema", "a2a.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + filepath.ToSlash(fileName) + "#/definitions/AgentCard")
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(ac)
	result, err := schema.Validate(gojsonschema.NewStringLoader(string(j)))
	if err != nil {
		logger.GetLogger().Errorf("schemaLoader.Validate return error: %v", err)
		return err
	}

	if !result.Valid() {
		logger.GetLogger().Error("a2a agent card is invalid: %v", result.Errors())
		return fmt.Errorf("a2a agent card is invalid: %v)", result.Errors())
	}

	return nil
}

func checkServiceRetries(retries int64) error {
	if retries >= 0 && retries <= 32767 {
		return nil
	}
	return fmt.Errorf("retries %d is not valid,should be in range [0,32767]", retries)
}

func checkAgentProtocol(ap string) error {
	if ap == "" || ap == "a2a" || ap == "acp" || ap == "mcp" {
		return nil
	}
	return fmt.Errorf("agent protocol %s is not valid,should be a2a/acp/mcp or empty", ap)
}

func CheckEphemeralCheck(check *t.EphemeralCheckInfo) error {
	if check == nil {
		return nil
	}

	check.CheckType = strings.ToUpper(check.CheckType)
	if check.CheckType != "" && check.CheckType != "TTL" {
		return fmt.Errorf("check type %s is not valid,should be TTL", check.CheckType)
	}
	if !checkTime(check.RenewalDeleteTimeout) {
		return fmt.Errorf("renewal delete timeout %s is not valid", check.RenewalDeleteTimeout)
	}
	if !checkTime(check.RenewalUnhealthyTimeout) {
		return fmt.Errorf("renewal unhealthy timeout %s is not valid", check.RenewalUnhealthyTimeout)
	}
	if !checkTime(check.RenewalInterval) {
		return fmt.Errorf("renewal interval %s is not valid", check.RenewalInterval)
	}
	return nil
}

func FillServiceDefaultValue(svc *t.Service) {
	svc.ID = BuildServiceID(svc.Name, svc.Host, svc.Port)

	svc.CreatedAt = t.Timestamp(time.Now())
	svc.UpdatedAt = t.Timestamp(time.Now())
	if svc.ConnectTimeout == 0 {
		svc.ConnectTimeout = 60000
	}
	if svc.WriteTimeout == 0 {
		svc.WriteTimeout = 60000
	}
	if svc.ReadTimeout == 0 {
		svc.ReadTimeout = 60000
	}
	if svc.Retries == 0 {
		svc.Retries = 5
	}
	if svc.Tags == nil {
		svc.Tags = []string{}
	}
	fillEphemeralCheckDefaultValue(svc)
	fillPersistentCheckDefaultValue(svc)
}

func fillEphemeralCheckDefaultValue(svc *t.Service) {
	if !svc.Ephemeral {
		svc.EphemeralCheck = nil
		return
	}

	if svc.EphemeralCheck == nil {
		return
	}
	if svc.EphemeralCheck.CheckType == "" {
		svc.EphemeralCheck.CheckType = "TTL"
	}
	if svc.EphemeralCheck.RenewalDeleteTimeout == "" {
		svc.EphemeralCheck.RenewalDeleteTimeout = "60s"
	}
	if svc.EphemeralCheck.RenewalUnhealthyTimeout == "" {
		svc.EphemeralCheck.RenewalUnhealthyTimeout = "15s"
	}
	if svc.EphemeralCheck.RenewalInterval == "" {
		svc.EphemeralCheck.RenewalInterval = "30s"
	}
}

func fillPersistentCheckDefaultValue(svc *t.Service) {
	if svc.Ephemeral {
		svc.PersistentCheck = nil
		return
	}

	if svc.PersistentCheck == nil {
		return
	}
	if svc.PersistentCheck.CheckUnhealthyTimeout == "" {
		svc.PersistentCheck.CheckUnhealthyTimeout = "30s"
	}
	if svc.PersistentCheck.CheckInterval == "" {
		svc.PersistentCheck.CheckInterval = "10s"
	}
	if svc.PersistentCheck.CheckTimeout == "" {
		svc.PersistentCheck.CheckTimeout = "5s"
	}
	if svc.PersistentCheck.CheckHTTPMethod == "" {
		svc.PersistentCheck.CheckHTTPMethod = "GET"
	}
}

func checkPersistentCheck(check *t.PersistentCheckInfo) error {
	if check == nil {
		return nil
	}

	if check.CheckHTTPURL == "" {
		return errors.New("check_http_url is not valid")
	}

	check.CheckType = strings.ToUpper(check.CheckType)
	if check.CheckType != "HTTP" && check.CheckType != "HTTPS" {
		return errors.New("check_type is not valid, support type:(HTTP,HTTPS)")
	}
	if !checkTime(check.CheckInterval) {
		return fmt.Errorf("check_interval is not valid")
	}
	if !checkTime(check.CheckTimeout) {
		return fmt.Errorf("check_timeout is not valid")
	}
	if !checkTime(check.CheckUnhealthyTimeout) {
		return fmt.Errorf("check_unhealthy_timeout is not valid")
	}

	check.CheckHTTPMethod = strings.ToUpper(check.CheckHTTPMethod)
	if check.CheckHTTPMethod != "" && check.CheckHTTPMethod != "GET" && check.CheckHTTPMethod != "HEAD" && check.CheckHTTPMethod != "OPTIONS" {
		return fmt.Errorf("check_http_method is not valid, support method:(GET,HEAD,OPTIONS)")
	}
	return nil
}

// checkAddress
var CheckAddress = func(address string) bool {
	if address == "0.0.0.0" || address == "::" {
		return false
	}

	// ip
	if IsValidIP(address) {
		return true
	}

	// domain
	return IsValidDomain(address)
}

func IsValidDomain(domain string) bool {
	if len(domain) > 255 {
		return false
	}
	if _, err := strconv.Atoi(domain[strings.LastIndex(domain, ".")+1:]); err == nil {
		return false
	}
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, domain)

	return isMatch
}

func IsValidIP(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	if ip.String() == "0.0.0.0" || ip.String() == "::" {
		return false
	}

	return true
}

func IsValidIPV6(address string) bool {
	if strings.Contains(address, ":") && IsValidIP(address) {
		return true
	}
	return false
}

func IsValidIPV4(address string) bool {
	if strings.Contains(address, ".") && IsValidIP(address) {
		return true
	}
	return false
}

// check Port
func CheckPort(port int) bool {
	if port < 1 || port > 65535 {
		return false
	}
	return true
}

func isValidProtocol(protocol string) bool {
	validProtocols := []string{"http", "https", "tcp", "udp", "tls", "grpc", "grpcs"}
	for _, p := range validProtocols {
		if protocol == p {
			return true
		}
	}

	return false
}

// chcck time
func checkTime(t string) bool {
	if t == "" {
		return true
	}
	match, _ := regexp.MatchString(`^([0-9]\d{0,2}[m|s])$|^([1-9]\d{0,2}m(([0-9])|([1-5][0-9]))s)$`, t)
	return match
}
