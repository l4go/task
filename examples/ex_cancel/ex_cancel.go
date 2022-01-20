package main

import (
	"log"
	"time"

	"github.com/l4go/task"
)

const (
	WORKERS = 3
)

func main() {
	log.Println("START")
	defer log.Println("END")

	c := task.NewCancel()

	w := WORKERS
	for w > 0 {
		go worker(c)
		w --
	}

	time.Sleep(3 * time.Second)
	c.Cancel()
	time.Sleep(3 * time.Second)
}

func worker(c task.Canceller) {
	log.Println("START worker")
	defer log.Println("END worker")

	go sub_worker(c)
	go sub_worker(c)

	select {
	case <-c.RecvCancel():
		log.Println("canceled!!")
		return
	}
}

func sub_worker(c task.Canceller) {
	log.Println("START sub_worker")
	defer log.Println("END sub_worker")

	select {
	case <-c.RecvCancel():
		log.Println("canceled!!")
		return
	}
}
