package leadermanager

import (
	"errors"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"uapregistry/logger"
)

var (
	lm *LeaderManager
)

type LeaderManager struct {
	sync.Mutex
	election *Election
	notify   map[chan bool]bool
}

type notify struct {
	selfIP string
}

func (n *notify) EventLeader(isleader bool) {
	if isleader {
		logger.GetLogger().Infof("%s is leader now", n.selfIP)
	} else {
		logger.GetLogger().Infof("%s is not leader now", n.selfIP)
	}

	if lm != nil {
		lm.NotifyLeader(isleader)
	}
}

func StartLeaderManager(client *consulapi.Client, selfIP string) {
	if client == nil {
		logger.GetLogger().Errorf("consul client is nil in StartLeaderManager")
		return
	}

	if selfIP == "" {
		logger.GetLogger().Errorf("selfIP is empty in StartLeaderManager")
		return
	}

	logger.GetLogger().Infof("starting leader manager,init isleader as false")

	n := &notify{
		selfIP: selfIP,
	}

	elconf := &ElectionConfig{
		Client:           client,
		Key:              "uapregistry/health-check/leader",
		Value:            selfIP,
		TTL:              "10s",
		CheckInterval:    2 * time.Second,
		SessionLockDelay: 2 * time.Second,
		Event:            n,
	}
	e := NewElection(elconf)
	go e.StartingElection()

	lm = &LeaderManager{
		election: e,
		notify:   make(map[chan bool]bool),
	}
}

func GetLeaderManager() (*LeaderManager, error) {
	if lm == nil {
		return nil, errors.New("leader manager is nil")
	}

	return lm, nil
}

func (m *LeaderManager) IsLeader() bool {
	return m.election.isLeader
}

func (m *LeaderManager) GetLeader() string {
	return m.election.leader
}

func (m *LeaderManager) NotifyLeader(isLeader bool) {
	m.Lock()
	defer m.Unlock()
	for ch := range m.notify {
		select {
		case ch <- isLeader:
		}
	}
}

func (m *LeaderManager) WaitLeaderCh(ch chan bool) {
	m.Lock()
	defer m.Unlock()
	if m.notify == nil {
		m.notify = make(map[chan bool]bool)
	}
	m.notify[ch] = false
}

func (m *LeaderManager) ClearLeaderCh(ch chan bool) {
	m.Lock()
	defer m.Unlock()
	if m.notify == nil {
		return
	}
	delete(m.notify, ch)
}
