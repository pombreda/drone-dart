package dispatch

import (
	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
	"github.com/drone/drone-dart/worker/director"
)

const BufferSize = 999

type request struct {
	context context.Context
	work    *worker.Work
}

// Dispatch is a thin wrapper around the Director that gives the
// ability to pause processing.
type Dispatch struct {
	*director.Director

	cancel  chan bool
	request chan *request
	stopped bool
}

func New() *Dispatch {
	return &Dispatch{
		cancel:  make(chan bool, 1),
		request: make(chan *request, BufferSize),
		stopped: true,
	}
}

// Do wraps the work request and pushes onto the
// channel for processing.
func (q *Dispatch) Do(c context.Context, work *worker.Work) {
	go func() {
		q.request <- &request{c, work}
	}()
}

// Start starts processing work requests in a background
// process.
func (q *Dispatch) Start() {
	q.Lock()
	q.stopped = false
	q.Unlock()

	go func() {
		for {
			select {
			// pickup a work request from the channel
			case r := <-q.request:
				q.Director.Do(r.context, r.work)
			// listen for a signal to exit
			case <-q.cancel:
				return
			}

			q.Lock()
			q.stopped = true
			q.Unlock()
		}
	}()
}

// Stop sends a signal to the Director to stop processing
// work requests.
func (q *Dispatch) Stop() {
	q.cancel <- true
}

// Stopped returns true if the Director is not processing
// work requests.
func (q *Dispatch) Stopped() bool {
	q.Lock()
	defer q.Unlock()
	return q.stopped
}
