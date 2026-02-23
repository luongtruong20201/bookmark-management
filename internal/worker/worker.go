package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// worker is a single goroutine that reads messages from the pool channel
// and invokes the Handler. On panic it reports to errChan for restart
// unless ctx is already cancelled (e.g. during shutdown).
type worker struct {
	id       int
	err      error
	handler  Handler
	wg       *sync.WaitGroup
	messages <-chan []byte
	errChan  chan<- *worker
}

// Work runs the worker loop until the messages channel is closed. Each
// message is passed to the Handler. Panics are recovered and sent to
// errChan so the pool can restart the worker.
func (w *worker) Work(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				w.err = err
			} else {
				w.err = fmt.Errorf("panic happened with %v", r)
			}
			if ctx.Err() != nil {
				log.Debug().Int("worker_id", w.id).Err(w.err).Msg("shutting down, skip reporting panic to pool")
				w.wg.Done()
				return
			}
			w.errChan <- w
		} else {
			w.wg.Done()
		}
	}()

	for {
		msg, ok := <-w.messages
		if !ok {
			log.Debug().Int("worker id", w.id).Msg("message channel closed, worker exiting")
			return
		}
		log.Debug().Msg(fmt.Sprintf("worker %d received message: %s", w.id, string(msg)))
		if err := w.handler.Handle(ctx, msg); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("worker %d received message error", w.id))
		} else {
			log.Debug().Msg(fmt.Sprintf("worker %d finished processing message: %s", w.id, string(msg)))
		}
	}
}
