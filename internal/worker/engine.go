package worker

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	"github.com/rs/zerolog/log"
)

// Engine is the main entry point for the worker process. It polls the queue
// for messages and feeds them to a pool of workers that invoke the Handler.
type Engine interface {
	// Start runs the engine until ctx is cancelled or a shutdown signal is received.
	Start(context.Context)
}

type engine struct {
	run        bool
	queue      Queue
	numWorkers int
	handler    Handler
	sigChan    chan os.Signal
	pool       *Pool
}

// Start runs the main loop: pop messages from the queue, pass them to the pool,
// and exit when ctx is done or a signal is received. Pool is closed on exit.
func (e *engine) Start(ctx context.Context) {
	log.Info().Msg("Starting worker...")
	e.run = true
	e.pool = NewPool(ctx, e.handler, e.numWorkers)

	signal.Notify(e.sigChan, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	for e.run {
		select {
		case sig := <-e.sigChan:
			log.Info().Str("signal", sig.String()).Msg("Received shutdown signal, stopping worker...")
			e.run = false
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping worker...")
			e.run = false
		default:
			msg, err := e.queue.PopMessage(ctx)
			if err != nil {
				if errors.Is(err, queue.ErrNoMessage) {
					log.Info().Msg("queue is empty")
					time.Sleep(time.Second)
					continue
				}
				log.Error().Err(err).Msg("cant pop message")
				time.Sleep(time.Second)
				continue
			}
			e.pool.Consume(msg)
		}
	}
	e.pool.Close()
}

// NewEngine builds an Engine that uses the given Queue and Handler, with
// numWorkers goroutines processing messages. The engine does not start until
// Start is called. Shutdown is driven by Start's context (ctx.Done()) or by
// process signals; the pool is closed when the engine loop exits.
func NewEngine(queue Queue, handler Handler, numWorkers int) Engine {
	return &engine{
		queue:      queue,
		numWorkers: numWorkers,
		handler:    handler,
		sigChan:    make(chan os.Signal, 1),
	}
}
