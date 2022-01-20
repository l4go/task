package main

import (
	"log"
	"time"

	"github.com/l4go/task"
	"github.com/l4go/timer"
)

const (
	WORKERS = 3
)

func main() {
	log.Println("START boss")
	defer log.Println("END boss")

	f := task.NewFinish()

	w := WORKERS
	for w > 0 {
		go worker(f)
		w --
	}

	select {
	case <-f.RecvDone():
		log.Println("Got 'done' from worker!!")
	}
	tm := timer.NewTimer()
	defer tm.Stop()
	tm.Start(time.Second * 3)
	<-tm.Recv()
}

func worker(f task.Finisher) {
	log.Println("START worker")
	defer log.Println("END worker")

	tm := timer.NewTimer()
	defer tm.Stop()
	tm.Start(time.Second * 3)
	select {
	case <-f.RecvDone():
		log.Println("Got 'done' from other one!!")
		return
	case <- tm.Recv():
		log.Println("Worked!!")
	}
	f.Done()
}
