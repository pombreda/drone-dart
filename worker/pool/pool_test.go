package pool

import (
	"testing"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
	"github.com/franela/goblin"
)

func TestPool(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Pool", func() {

		g.It("Should allocate workers", func() {
			w := mockWorker{}
			pool := New()
			pool.Allocate(&w)
			g.Assert(len(pool.workers)).Equal(1)
			g.Assert(len(pool.workerc)).Equal(1)
			g.Assert(pool.workers[&w]).Equal(true)
		})

		g.It("Should reserve a worker", func() {
			w := mockWorker{}
			pool := New()
			pool.Allocate(&w)
			g.Assert(<-pool.Reserve()).Equal(&w)
		})

		g.It("Should release a worker", func() {
			w := mockWorker{}
			pool := New()
			pool.Allocate(&w)
			g.Assert(len(pool.workerc)).Equal(1)
			g.Assert(<-pool.Reserve()).Equal(&w)
			g.Assert(len(pool.workerc)).Equal(0)
			pool.Release(&w)
			g.Assert(len(pool.workerc)).Equal(1)
			g.Assert(<-pool.Reserve()).Equal(&w)
			g.Assert(len(pool.workerc)).Equal(0)
		})

		g.It("Should not release an unallocated worker", func() {
			w := mockWorker{}
			pool := New()
			g.Assert(len(pool.workers)).Equal(0)
			g.Assert(len(pool.workerc)).Equal(0)
			pool.Release(&w)
			g.Assert(len(pool.workers)).Equal(0)
			g.Assert(len(pool.workerc)).Equal(0)
		})

		g.It("Should list all allocated workers", func() {
			w1 := mockWorker{}
			w2 := mockWorker{}
			pool := New()
			pool.Allocate(&w1)
			pool.Allocate(&w2)
			g.Assert(len(pool.workers)).Equal(2)
			g.Assert(len(pool.workerc)).Equal(2)
			g.Assert(len(pool.List())).Equal(2)
		})

		g.It("Should remove a worker")

	})
}

// fake worker for testing purpose only
type mockWorker struct {
	name string
}

func (*mockWorker) Do(c context.Context, w *worker.Work) {}
