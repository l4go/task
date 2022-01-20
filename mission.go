package task

import (
	"context"
	"sync/atomic"
	"time"
)

type Mission struct {
	p *Mission

	cancel  chan struct{}
	cc_once uint32

	done    chan struct{}
	dn_cnt  int32
	dn_mode uint32
}

const (
	mode_off uint32 = iota
	mode_on
)

func NewMission() *Mission {
	return &Mission{
		p:       nil,
		cancel:  make(chan struct{}),
		cc_once: 1,
		done:    make(chan struct{}),
		dn_cnt:  1,
		dn_mode: mode_off,
	}
}

func (m *Mission) on_done() {
	if atomic.CompareAndSwapUint32(&m.dn_mode, mode_off, mode_on) {
		m.decr_done()
	}
}

func (m *Mission) incr_done() {
	atomic.AddInt32(&m.dn_cnt, 1)
}

func (m *Mission) decr_done() {
	cnt := atomic.AddInt32(&m.dn_cnt, -1)
	if cnt == 0 {
		close(m.done)
	}
}

func (m *Mission) Parson() *Mission {
	return m.p
}

func (p *Mission) New() *Mission {
	if p == nil {
		panic("nil mission")
	}
	p.incr_done()
	select {
	case <-p.done:
		panic("use of done mission")
	default:
	}

	c := NewMission()
	c.p = p

	go c.chain_cancel()
	return c
}

func (m *Mission) chain_cancel() {
	select {
	case <-m.cancel:
	case <-m.p.done:
	case <-m.p.cancel:
		m.Cancel()
	}
}

func (m *Mission) Activate() {
	m.on_done()
}

func (m *Mission) Done() {
	m.on_done()
	<-m.done
	if m.p != nil {
		m.p.decr_done()
	}
}

func (m *Mission) NowaitDone() {
	m.on_done()
	if m.p != nil {
		m.p.decr_done()
	}
}

func (m *Mission) Cancel() {
	win := atomic.SwapUint32(&m.cc_once, 0)
	if win != 0 {
		close(m.cancel)
	}
}

func (m *Mission) Abort() {
	top := m
	for top.p != nil {
		top = top.p
	}
	top.Cancel()
}

func (m *Mission) IsCanceled() bool {
	select {
	case <-m.cancel:
		return true
	default:
	}

	return false
}

func (m *Mission) RecvCancel() <-chan struct{} {
	return m.cancel
}

func (m *Mission) RecvDone() <-chan struct{} {
	return m.done
}

func (m *Mission) Recv() <-chan struct{} {
	m.Activate()
	return m.RecvDone()
}

func (m *Mission) AsContext() context.Context {
	return (*MissionContext)(m)
}

type MissionContext Mission

func (c *MissionContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c *MissionContext) Done() <-chan struct{} {
	return (*Mission)(c).RecvCancel()
}

func (c *MissionContext) Err() error {
	select {
	case <-(*Mission)(c).RecvCancel():
		return context.Canceled
	default:
	}

	return nil
}

func (*MissionContext) Value(key interface{}) interface{} {
	return nil
}
