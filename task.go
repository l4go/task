package task

type Canceller interface {
	Cancel()
	RecvCancel() <-chan struct{}
}
