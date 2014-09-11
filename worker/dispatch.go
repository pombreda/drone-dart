package worker

import (
	"log"

	"github.com/drone/drone-dart/dart"
)

// http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

// Queue is a generic implementation of a build queue
// that can accept work requests.
type Queue interface {
	Send(r *Request)
	SendPackage(pkg *dart.Package, sdk *dart.SDK)
}

// Dispatch implements a simple FIFO queue, dispatching work
// requests to the first available worker node.
type Dispatch struct {
	requests chan *Request
	workers  chan chan *Request
	quit     chan bool
}

func NewDispatch(requests chan *Request, workers chan chan *Request) *Dispatch {
	return &Dispatch{
		requests: requests,
		workers:  workers,
		quit:     make(chan bool),
	}
}

// Start tells the dispatcher to start listening
// for work requests and dispatching to workers.
func (d *Dispatch) Start() {
	go func() {
		for {
			select {
			// pickup a request from the queue
			case request := <-d.requests:
				go func() {
					// find an available worker and
					// send the request to that worker
					worker := <-d.workers
					worker <- request
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

// Send sends a Request to the queue to be dispatched
// to a worker node.
func (d *Dispatch) Send(r *Request) {
	go func() { d.requests <- r }()
}

// SendPackage sends a Package to the queue to be
// dispatched to a worker node.
func (d *Dispatch) SendPackage(pkg *dart.Package, sdk *dart.SDK) {
	log.Printf("Queue build %s for sdk %s\n", pkg.Name, sdk.Version)
	d.Send(&Request{pkg, sdk})
}
