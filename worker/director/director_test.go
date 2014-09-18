package director

import (
	"testing"

	"code.google.com/p/go.net/context"
	"github.com/drone/drone-dart/worker/pool"
	"github.com/franela/goblin"
)

func TestDirector(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Supervisor", func() {
		g.It("Should return pending work")
		g.It("Should return started work")
		g.It("Should return allocations")
		g.It("Should mark work as pending")
		g.It("Should mark work as started")
		g.It("Should mark work as complete")
	})
}

// fake worker for testing purpose only
type mockWorker struct {
	name string
}

func (*mockWorker) Do(c context.Context, w *worker.Work) {}
