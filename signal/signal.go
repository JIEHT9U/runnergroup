package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ShutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

type Shutdown struct {
	once    sync.Once
	ctx     context.Context
	cancel  context.CancelFunc
	signals []os.Signal
}

func New(signals []os.Signal) *Shutdown {
	var shutdown = &Shutdown{
		signals: signals,
	}
	shutdown.ctx, shutdown.cancel = context.WithCancel(context.Background())
	return shutdown
}

func (s *Shutdown) Listen() context.Context {
	stopSignal := make(chan os.Signal, 2)
	signal.Notify(stopSignal, s.signals...)
	go func() {
		<-stopSignal
		s.Cancel()
		<-stopSignal
		os.Exit(1)
	}()
	return s.ctx
}

func (s *Shutdown) Cancel() {
	s.once.Do(func() {
		s.cancel()
	})
}
