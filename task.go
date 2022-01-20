package task

import "context"

type Canceller interface {
	Cancel()
	RecvCancel() <-chan struct{}
	AsContext() context.Context
}

func IsCanceled(cc Canceller) bool {
	select {
	case <-cc.RecvCancel():
		return true
	default:
	}
	return false
}

type Finisher interface {
	Done()
	RecvDone() <-chan struct{}
	AsContext() context.Context
}
