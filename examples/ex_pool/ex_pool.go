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
	WORKERS        = 5
	TASKS_PER_STEP = 10
	TASK_MSEC      = 250
)

func main() {
	top_m := task.NewMission()
	defer top_m.Done()

	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	log.Println("START")
	defer log.Println("END")
	p := task.NewPool(top_m.New(), WORKERS)
	defer p.Close()

	step_m := top_m.New()
	defer step_m.Done()

	for i := 1; i <= TASKS_PER_STEP; i++ {
		p.Do(must_work, step_m.New(), "must", i)
		p.WeakDo(weak_work, "weak", i)
	}

	log.Println("CANCEL")
	top_m.Cancel()

	log.Println("WAIT")
	select {
	case <-step_m.Recv():
	case <-signal_ch:
		top_m.Cancel()
	}
}

func must_work(m *task.Mission, args ...interface{}) {
	defer m.Done()

	s := args[0].(string)
	i := args[1].(int)
	if m.IsCanceled() {
		log.Println("work", s, i, "canceled")
		return
	}

	log.Println("work", s, i, "start")

	tm := timer.NewTimer()
	defer tm.Stop()

	tm.Start(TASK_MSEC * time.Millisecond)
	select {
	case <-m.RecvCancel():
		log.Println("work", s, i, "cancel")
	case <-tm.Recv():
		log.Println("work", s, i, "end")
	}
}

func weak_work(args ...interface{}) {
	s := args[0].(string)
	i := args[1].(int)

	log.Println("work", s, i, "start")
	defer log.Println("work", s, i, "end")

	tm := timer.NewTimer()
	defer tm.Stop()

	tm.Start(TASK_MSEC * time.Millisecond)
	<-tm.Recv()
}
