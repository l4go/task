package task

import (
	"sync"
)

type PoolFunc func(*Mission, ...interface{})
type PoolWeakFunc func(...interface{})

type poolParam struct {
	f   PoolFunc
	m   *Mission
	arg []interface{}
}

type poolWeakParam struct {
	f   PoolWeakFunc
	arg []interface{}
}

type Pool struct {
	param_ch chan interface{}
	m        *Mission

	cc_once sync.Once
}

func NewPool(m *Mission, cnt int) *Pool {
	pch := make(chan interface{}, cnt)
	p := &Pool{m: m, param_ch: pch}

	for cnt > 0 {
		go p.worker(p.m.New())
		cnt--
	}

	return p
}

func (p *Pool) worker(m *Mission) {
	defer m.Done()

	for i := range p.param_ch {
		switch param := i.(type) {
		case poolParam:
			select {
			case <-p.m.RecvCancel():
				param.m.Cancel()
			default:
			}
			param.f(param.m, param.arg...)
		case poolWeakParam:
			select {
			case <-p.m.RecvCancel():
			default:
				param.f(param.arg...)
			}
		}
	}
}

func (p *Pool) cancel() {
	p.m.Cancel()
	close(p.param_ch)
}

func (p *Pool) Cancel() {
	p.cc_once.Do(p.cancel)
}

func (p *Pool) Recv() <-chan struct{} {
	return p.m.RecvDone()
}

func (p *Pool) Close() {
	p.Cancel()
	p.m.Done()
}

func (p *Pool) Do(f PoolFunc, m *Mission, args ...interface{}) {
	select {
	case <-p.m.RecvCancel():
		m.Cancel()
		f(m, args...)
	case p.param_ch <- poolParam{f: f, m: m, arg: args}:
	}
}

func (p *Pool) WeakDo(f PoolWeakFunc, args ...interface{}) {
	select {
	case <-p.m.RecvCancel():
	case p.param_ch <- poolWeakParam{f: f, arg: args}:
	}
}
