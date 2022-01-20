package task

import (
	"context"
)

type Finish Cancel

func NewFinish() *Finish {
	return (*Finish)(NewCancel())
}

func (f *Finish) Done() {
	(*Cancel)(f).Cancel()
}

func (f *Finish) RecvDone() <-chan struct{} {
	return (*Cancel)(f).RecvCancel()
}

func (f *Finish) AsCanceller() Canceller {
	return (*Cancel)(f)
}

func (f *Finish) AsContext() context.Context {
	return (*Cancel)(f).AsContext()
}

