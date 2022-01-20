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

	for w := 0; w < WORKERS; w++ {
		go worker(m.New(), w)
	}

	select {
	case <-m.Recv():
	case <-signal_ch:
		m.Cancel()
	}
}

func worker(m *task.Mission, w int) {
	log.Println("worker", w, "START")
	defer log.Println("worker", w, "END")
	defer m.Done()

	for sw := 0; sw < SUBWORKERS; sw++ {
		if task.IsCanceled(m) {
			return
		}
		go subworker(m.New(), w, sw)
	}
}

func subworker(m *task.Mission, w int, sw int) {
	log.Println("subworker", w, sw, "START")
	defer log.Println("subworker", w, sw, "END")
	defer m.Done()

	work_time := time.Duration(w) * time.Second

	if task.IsCanceled(m) {
		return
	}

	wt := timer.NewTimer()
	defer wt.Stop()
	wt.Start(work_time)

	job(m.NewCancel(), w, sw)

	select {
	case <-wt.Recv():
	case <-m.RecvCancel():
	}
}

func job(cc task.Canceller, w int, sw int) {
	log.Println("job", w, sw, "START")
	defer log.Println("job", w, sw, "END")
	defer cc.Cancel()

	job_time := time.Duration(w) * time.Second / 2

	tt := timer.NewTimer()
	defer tt.Stop()
	tt.Start(job_time)

	select {
	case <-tt.Recv():
	case <-cc.RecvCancel():
	}
}
