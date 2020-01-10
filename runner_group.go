package runnergroup

import (
	"context"
	"sync"
)

type worker struct {
	execute   func() error
	interrupt func()
}

func (w *worker) run() error {
	return w.execute()
}

func (w *worker) stop() {
	w.interrupt()
}

type Runnergroup struct {
	runWg   sync.WaitGroup
	workers []worker
	once    sync.Once
	errors  chan error
}

func New() *Runnergroup {
	return &Runnergroup{
		workers: make([]worker, 0),
		errors:  make(chan error),
	}
}

func (runner *Runnergroup) Register(execute func() error, interrupt func()) {
	runner.workers = append(runner.workers, worker{
		execute:   execute,
		interrupt: interrupt,
	})
}

func (runner *Runnergroup) Run(ctx context.Context) (err error) {
	defer runner.closeErrorsChanel()

	if len(runner.workers) == 0 {
		return nil
	}

	// Run each worker.
	for _, w := range runner.workers {
		runner.run(w)
	}

loop:
	for {
		select {
		case <-ctx.Done():
			err = nil
			break loop
		case err = <-runner.errors:
			break loop
		}
	}

	for _, w := range runner.workers {
		w.stop()
	}

	runner.runWg.Wait()

	// Return the original error.
	return err
}

func (runner *Runnergroup) run(w worker) {
	runner.runWg.Add(1)
	go func(w worker, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := w.run(); err != nil {
			select {
			case runner.errors <- err:
			default:
				return
			}
		}
	}(w, &runner.runWg)
}

func (runner *Runnergroup) closeErrorsChanel() {
	runner.once.Do(func() {
		close(runner.errors)
	})
}
