package leadermanager

import (
	"errors"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
)

const (
	interval = 5
)

type Notifier interface {
	EventLeader(e bool)
}

type Election struct {
	Client           *consulapi.Client
	isLeader         bool
	leader           string
	sessionID        string
	Key              string
	Value            string
	TTL              string
	CheckInterval    time.Duration
	SessionLockDelay time.Duration
	Event            Notifier
}

type ElectionConfig struct {
	Client           *consulapi.Client
	Key              string
	Value            string
	TTL              string
	CheckInterval    time.Duration
	SessionLockDelay time.Duration
	Event            Notifier
}

func NewElection(c *ElectionConfig) *Election {
	e := &Election{
		Client:           c.Client,
		isLeader:         false,
		Key:              c.Key,
		Value:            c.Value,
		TTL:              c.TTL,
		CheckInterval:    c.CheckInterval,
		SessionLockDelay: c.SessionLockDelay,
		Event:            c.Event,
	}
	return e
}

func (e *Election) StartingElection() {
	go e.watchLeader()

	e.process()
	ticker := time.NewTicker(e.CheckInterval)
	for range ticker.C {
		e.process()
	}
}

func (e *Election) watchLeader() {
	queryOptions := &consulapi.QueryOptions{
		WaitIndex:         0,
		WaitTime:          10 * time.Minute,
		RequireConsistent: true,
	}

	for {
		p, qm, err := e.Client.KV().Get(e.Key, queryOptions)
		if err != nil {
			logger.GetLogger().Errorf("failed to get kv(key:%s):%v,retry after %d seconds", e.Key, err, interval)
			queryOptions.WaitIndex = 0
			time.Sleep(interval * time.Second)
			continue
		}

		var leaderNew string
		if p != nil && p.Session != "" && string(p.Value) != "" {
			leaderNew = string(p.Value)
		} else {
			leaderNew = ""
		}

		if leaderNew != e.leader {
			logger.GetLogger().Infof("health-check leader has changed,old:%s,new:%s", e.leader, leaderNew)
			e.leader = leaderNew
		}
		queryOptions.WaitIndex = qm.LastIndex
	}
}

func (e *Election) process() {
	e.waitSession()
	if !e.isLeader {
		if !e.isNeedAquire() {
			return
		}
		logger.GetLogger().Debugf("Try to acquire")
		res, err := e.acquire()
		if res && err == nil {
			e.enableLeader()
		}
	}
}

// 检查session是否存在，存在就renew,不存在就创建
func (e *Election) waitSession() {
	err := e.processSession()
	if err == nil {
		return
	}
	ticker := time.NewTicker(e.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		err = e.processSession()
		if err == nil {
			return
		}
	}
}

func (e *Election) processSession() error {
	isset, err := e.checkSession()
	if isset {
		_, _, err = e.Client.Session().Renew(e.sessionID, nil)
		if err != nil {
			logger.GetLogger().Debugf(" e.Client.Session().Renew failed")
		}
		return nil
	}
	e.disableLeader()
	if err != nil {
		logger.GetLogger().Debugf("Try to get session info again.")
		return err
	}

	err = e.createSession()
	if err == nil {
		logger.GetLogger().Debugf("Session " + e.sessionID + " created")
		return err
	}
	return nil
}

func (e *Election) createSession() (err error) {
	ses := &consulapi.SessionEntry{
		TTL:       e.TTL,
		LockDelay: e.SessionLockDelay,
	}
	e.sessionID, _, err = e.Client.Session().Create(ses, nil)
	if err != nil {
		logger.GetLogger().Errorf("failed to create Session:%v", err)
	}
	return
}

func (e *Election) checkSession() (bool, error) {
	if e.sessionID == "" {
		return false, nil
	}
	res, _, err := e.Client.Session().Info(e.sessionID, nil)
	if err != nil {
		logger.GetLogger().Errorf("failed to get Session(id:%s):%v", e.sessionID, err)
	}
	return res != nil, err
}

// 占用当前kv的sessionID为空才返回true
func (e *Election) isNeedAquire() bool {
	res, err := e.waitSessionData()
	if err != nil {
		return false
	}

	if e.sessionID != "" && e.sessionID == res {
		e.enableLeader()
	}
	if res == "" || res != e.sessionID {
		e.disableLeader()
	}

	return res == ""
}

// 获取当前哪个Session占用了kv
func (e *Election) waitSessionData() (string, error) {
	if res, err := e.getKvSession(); err == nil {
		return res, nil
	}

	e.disableLeader()
	ticker := time.NewTicker(e.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		if res, err := e.getKvSession(); err == nil {
			return res, nil
		}
	}

	// 理论上不会执行到这里
	return "", errors.New("unexpected exit from wait loop")
}

func (e *Election) getKvSession() (string, error) {
	p, _, err := e.Client.KV().Get(e.Key, nil)
	if err != nil {
		logger.GetLogger().Errorf("failed to get kv(key:%s):%v", e.Key, err)
		return "", err
	}
	if p == nil {
		return "", nil
	}
	return p.Session, nil
}

func (e *Election) acquire() (bool, error) {
	kv := &consulapi.KVPair{
		Key:     e.Key,
		Session: e.sessionID,
		Value:   []byte(e.Value),
	}
	res, _, err := e.Client.KV().Acquire(kv, nil)
	if err != nil {
		logger.GetLogger().Errorf("failed to acquire kv(key:%s,value:%s,sessionID:%s):%v", e.Key, e.Value, e.sessionID, err)
	}
	return res, err
}

func (e *Election) enableLeader() {
	e.isLeader = true
	logger.GetLogger().Info("I'm a leader!")
	if e.Event != nil {
		e.Event.EventLeader(true)
	}
}

func (e *Election) disableLeader() {
	if e.isLeader {
		e.isLeader = false
		logger.GetLogger().Info("I'm not a leader")
		if e.Event != nil {
			e.Event.EventLeader(false)
		}
	}
}
