package worker

import (
	"log"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/dart"
)

// http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

// Queue is a generic implementation of a worker queue
// that can accept work requests.
type Queue interface {
	Send(context.Context, *Worker)
}

// Dispatch implements a simple FIFO queue, dispatching work
// requests to the first available worker node.
type Dispatch struct {
	work    chan *Work
	workers chan chan *Work
	quit    chan bool
}

func NewDispatch(work chan *Work, workers chan chan *Work) *Dispatch {
	return &Dispatch{
		work:    work,
		workers: workers,
		quit:    make(chan bool),
	}
}

// Start tells the dispatcher to start listening
// for work requests and dispatching to workers.
func (d *Dispatch) Start() {
	go func() {
		for {
			select {
			// pickup a work request from the queue
			case work := <-d.work:
				go func() {
					// find an available worker and
					// send the request to that worker
					worker := <-d.workers
					worker <- work
				}()
			// listen for a signal to exit
			case <-d.quit:
				return
			}
		}
	}()

}

// Stop tells the dispatcher to stop listening for new
// work requests.
func (d *Dispatch) Stop() {
	go func() { d.quit <- true }()
}

// Send sends a work request to the queue to be dispatched
// to a worker node.
func (d *Dispatch) Send(w *Work) {
	go func() { d.work <- w }()
}
