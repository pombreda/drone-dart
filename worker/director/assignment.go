package director

import (
	"github.com/drone/drone-dart/worker"
)

type Assignment struct {
	Work   *worker.Work
	Worker worker.Worker
}
