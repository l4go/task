package task

import (
	"context"
	"sync/atomic"
	"time"
)

type Cancel struct {
	done    chan struct{}
	c_once  uint32
}

func NewCancel() *Cancel {
	return &Cancel{
		c_once:   1,
		done:    make(chan struct{}),
	}
}

func (c *Cancel) Cancel() {
	c.do_cancel()
}

func (c *Cancel) RecvCancel() <-chan struct{} {
	return c.done
}

func (c *Cancel) do_cancel() {
	win := atomic.SwapUint32(&c.c_once, 0)
	if win == 0 {
		return
	}
	close(c.done)
}

func (c *Cancel) AsContext() context.Context {
	return (*CancelContext)(c)
}

type CancelContext Cancel

func (cc *CancelContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (cc *CancelContext) Done() <-chan struct{} {
	return (*Cancel)(cc).RecvCancel()
}

func (cc *CancelContext) Err() error {
	select {
	case <-(*Cancel)(cc).RecvCancel():
		return context.Canceled
	default:
	}

	return nil
}

func (*CancelContext) Value(key interface{}) interface{} {
	return nil
}
