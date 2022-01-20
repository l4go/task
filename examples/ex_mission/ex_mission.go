package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/l4go/task"
	"github.com/l4go/timer"
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

const (
	WORKERS    = 3
	SUBWORKERS = 3
)

func main() {
	log.Println("START")
	defer log.Println("END")

	m := task.NewMission()
	defer m.Done()

	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	time.AfterFunc(10*time.Second, m.Cancel)

	w := WORKERS
	for w > 0 {
		go worker(m.New())
		w--
	}

	select {
	case <-m.Recv():
	case <-signal_ch:
		m.Cancel()
	}
}

func worker(m *task.Mission) {
	log.Println("worker START")
	defer log.Println("worker END")
	defer m.Done()

	var w int = SUBWORKERS
	for w > 0 {
		if m.IsCanceled() {
			return
		}
		go subworker(m.New(), w)
		w--
	}
}

func subworker(m *task.Mission, work_time int) {
	log.Println("subworker START")
	defer log.Println("subworker END")
	defer m.Done()

	if m.IsCanceled() {
		return
	}

	wk := timer.NewTimer()
	defer wk.Stop()
	wk.Start(time.Duration(work_time) * time.Second)

	select {
	case <-wk.Recv():
	case <-m.RecvCancel():
	}
}
