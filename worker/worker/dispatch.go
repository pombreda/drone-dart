package worker

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
)

// Dispatch implements a simple FIFO queue, dispatching work
// requests to the first available worker node.
type Dispatch struct {
	work    chan *worker.Work
	workers chan chan *worker.Work
	quit    chan bool
}

func NewDispatch(work chan *worker.Work, workers chan chan *worker.Work) *Dispatch {
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
func (d *Dispatch) Send(c context.Context, w *worker.Work) {
	go func() { d.work <- w }()
}
