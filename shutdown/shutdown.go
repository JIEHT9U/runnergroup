package shutdown

import (
	"context"
	"sync"
)

type Shutdown struct {
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context) *Shutdown {
	ctx, cancel := context.WithCancel(ctx)
	return &Shutdown{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Shutdown) Shutdown() {
	s.once.Do(func() {
		s.cancel()
	})
}

func (s *Shutdown) Context() context.Context {
	return s.ctx
}
