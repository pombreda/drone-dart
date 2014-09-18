package director

import (
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
	"github.com/drone/drone-dart/worker/pool"
)

// Director is a simple FIFO queue that delegates tasks
// accross multiple runners.
type Director struct {
	sync.Mutex

	pending map[*worker.Work]bool
	started map[*worker.Work]worker.Worker
}

func New() *Director {
	return &Director{
		pending: make(map[*worker.Work]bool),
		started: make(map[*worker.Work]worker.Worker),
	}
}

// Do processes the work request async.
func (d *Director) Do(c context.Context, work *worker.Work) {
	go func() {
		d.do(c, work)
	}()
}

// do is a blocking function that waits for an
// available worker to process work.
func (d *Director) do(c context.Context, work *worker.Work) {
	defer func() {
		recover()
	}()

	d.markPending(work)
	var pool = pool.FromContext(c)
	var worker = <-pool.Reserve()

	d.markStarted(work, worker)
	worker.Do(c, work)
	pool.Release(worker)
	d.markComplete(work)
}

// GetStarted returns a list of all jobs that
// are assigned and being worked on.
func (d *Director) GetStarted() []*worker.Work {
	d.Lock()
	defer d.Unlock()
	var started []*worker.Work
	for work, _ := range d.started {
		started = append(started, work)
	}
	return started
}

// GetPending returns a list of all work that
// is pending assignment to a worker.
func (d *Director) GetPending() []*worker.Work {
	d.Lock()
	defer d.Unlock()
	var pending []*worker.Work
	for work, _ := range d.pending {
		pending = append(pending, work)
	}
	return pending
}

// GetAssignments returns a list of assignments. The
// assignment type is a structure that stores the
// work being performed and the assigned worker.
func (d *Director) GetAssignemnts() []*Assignment {
	d.Lock()
	defer d.Unlock()
	var assignments []*Assignment
	for work, worker := range d.started {
		assignment := &Assignment{work, worker}
		assignments = append(assignments, assignment)
	}
	return assignments
}

func (d *Director) markPending(work *worker.Work) {
	d.Lock()
	defer d.Unlock()
	delete(d.started, work)
	d.pending[work] = true
}

func (d *Director) markStarted(work *worker.Work, worker worker.Worker) {
	d.Lock()
	defer d.Unlock()
	delete(d.pending, work)
	d.started[work] = worker
}

func (d *Director) markComplete(work *worker.Work) {
	d.Lock()
	defer d.Unlock()
	delete(d.pending, work)
	delete(d.started, work)
}
