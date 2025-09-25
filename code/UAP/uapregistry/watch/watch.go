package watch

import (
	"sync"
	"time"

	"uapregistry/logger"
)

// Watch is the external interface that's common to all the different flavors.
type Watch interface {
	// Wait registers the given channel and calls it back when the watch
	// fires.
	Wait(notifyCh chan struct{})

	// Clear deregisters the given channel.
	Clear(notifyCh chan struct{})

	//******************************
	Notify()

	NotifyWithTimeOut(timeOut time.Duration)

	NoBlockingNotify()

	// Wait registers the given channel and calls it back when the watch
	// fires.
	WaitErrCh(notifyCh chan error)

	// Clear deregisters the given channel.
	ClearErrCh(notifyCh chan error)

	//******************************
	NotifyErrCh(error)

	WaitAnyCh(notifyCh chan AnyObject)

	ClearAnyCh(notifyCh chan AnyObject)

	NotifyAnyWithTimeOut(obj AnyObject, timeOut time.Duration)
}

type NotifyGroup struct {
	l         sync.Mutex
	notify    map[chan struct{}]struct{}
	notifyErr map[chan error]error
	notifyAny map[chan AnyObject]AnyObject
}

type AnyObject interface{}

// NewNotifyGroup returns a new watch.
func NewNotifyGroup() *NotifyGroup {
	return &NotifyGroup{}
}

// Notify will do a blocking send to all waiting channels, and
// clear the notify list
func (n *NotifyGroup) Notify() {
	n.l.Lock()
	defer n.l.Unlock()
	for ch := range n.notify {
		select {
		case ch <- struct{}{}:
			//default:
		}
	}
	n.notify = nil
}

// Notify will do a blocking send to all waiting channels, and do not clear the notify list
func (n *NotifyGroup) NotifyWithTimeOut(timeOut time.Duration) {
	n.l.Lock()
	defer n.l.Unlock()
	for ch := range n.notify {
		timeoutCh := time.After(timeOut)
		select {
		case ch <- struct{}{}:
		case <-timeoutCh:
			logger.GetLogger().Warn("timeout occurred while notify")
		}
	}
}

func (n *NotifyGroup) NoBlockingNotify() {
	n.l.Lock()
	defer n.l.Unlock()
	for ch := range n.notify {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	//n.notify = nil
}

// Wait adds a channel to the notify group
func (n *NotifyGroup) Wait(ch chan struct{}) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notify == nil {
		n.notify = make(map[chan struct{}]struct{})
	}
	n.notify[ch] = struct{}{}
}

// Clear removes a channel from the notify group
func (n *NotifyGroup) Clear(ch chan struct{}) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notify == nil {
		return
	}
	delete(n.notify, ch)
}

// Notify will do a blocking send to all waiting channels
func (n *NotifyGroup) NotifyErrCh(err error) {
	n.l.Lock()
	defer n.l.Unlock()
	for ch := range n.notifyErr {
		select {
		case ch <- err:
		default:
		}
	}
}

// Wait adds a channel to the notify group
func (n *NotifyGroup) WaitErrCh(ch chan error) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notifyErr == nil {
		n.notifyErr = make(map[chan error]error)
	}
	n.notifyErr[ch] = nil
}

// Clear removes a channel from the notify group
func (n *NotifyGroup) ClearErrCh(ch chan error) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notifyErr == nil {
		return
	}
	delete(n.notifyErr, ch)
}

// Notify will do a blocking send to all waiting channels
func (n *NotifyGroup) NotifyAnyWithTimeOut(obj AnyObject, timeOut time.Duration) {
	n.l.Lock()
	defer n.l.Unlock()
	for ch := range n.notifyAny {
		timeoutCh := time.After(timeOut)
		select {
		case ch <- obj:
		case <-timeoutCh:
			logger.GetLogger().Warn("timeout occurred while notify")
		}
	}
}

// Wait adds a channel to the notify group
func (n *NotifyGroup) WaitAnyCh(ch chan AnyObject) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notifyAny == nil {
		n.notifyAny = make(map[chan AnyObject]AnyObject)
	}
	n.notifyAny[ch] = nil
}

// Clear removes a channel from the notify group
func (n *NotifyGroup) ClearAnyCh(ch chan AnyObject) {
	n.l.Lock()
	defer n.l.Unlock()
	if n.notifyAny == nil {
		return
	}
	delete(n.notifyAny, ch)
}
