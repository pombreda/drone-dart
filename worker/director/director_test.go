package director

import (
	"testing"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker"
	"github.com/franela/goblin"
)

func TestDirector(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Supervisor", func() {

		g.It("Should mark work as pending", func() {
			d := New()
			d.markPending(&worker.Work{})
			d.markPending(&worker.Work{})
			g.Assert(len(d.GetPending())).Equal(2)
		})
		g.It("Should mark work as started", func() {
			d := New()
			w1 := worker.Work{}
			w2 := worker.Work{}
			d.markPending(&w1)
			d.markPending(&w2)
			g.Assert(len(d.GetPending())).Equal(2)
			d.markStarted(&w1, &mockWorker{})
			g.Assert(len(d.GetStarted())).Equal(1)
			g.Assert(len(d.GetPending())).Equal(1)
			d.markStarted(&w2, &mockWorker{})
			g.Assert(len(d.GetStarted())).Equal(2)
			g.Assert(len(d.GetPending())).Equal(0)
		})
		g.It("Should mark work as complete", func() {
			d := New()
			w1 := worker.Work{}
			w2 := worker.Work{}
			d.markStarted(&w1, &mockWorker{})
			d.markStarted(&w2, &mockWorker{})
			g.Assert(len(d.GetStarted())).Equal(2)
			d.markComplete(&w1)
			g.Assert(len(d.GetStarted())).Equal(1)
			d.markComplete(&w2)
			g.Assert(len(d.GetStarted())).Equal(0)
		})
	})
}

// fake worker for testing purpose only
type mockWorker struct {
	name string
}

func (*mockWorker) Do(c context.Context, w *worker.Work) {}
