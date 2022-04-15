package worker

import (
	"errors"
	"time"

	"github.com/dostow/worker/pkg/queues/machinery"
)

// Worker a worker that communicates with rave
type Worker struct {
	Timeout   time.Duration `help:"gocent timeout"`
	ID        string        `help:"worker id"`
	Build     string        `help:"build"`
	Dry       bool          `help:"dry run"`
	DostowAPI string        `help:"dostow api url" env:"DOSTOW_API"`
}

// Run run the worker
func (w *Worker) Run() error {
	return machinery.Worker(w.ID, map[string]interface{}{
		"rave": func(args ...string) error {
			return doRave(w.DostowAPI, args[0], args[1], args[2], args[3], w.Dry)
		},
	})
}

// Send a job to another worker
func (w *Worker) Send() error {
	return errors.New("not implemented")
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	return &Worker{Timeout: 5 * time.Second, Build: build, Dry: false}
}
