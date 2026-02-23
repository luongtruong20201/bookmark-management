package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Pool manages a fixed number of worker goroutines that read messages from
// an internal channel and call the Handler. Panicking workers are reported
// on errChan and restarted after a delay.
type Pool struct {
	handler    Handler
	numWorkers int
	wg         *sync.WaitGroup
	messages   chan []byte
	errChan    chan *worker
}

// NewPool creates a pool with numWorkers goroutines running Handler. Context
// is used for graceful shutdown (workers check ctx on panic). Call Close
// when done to stop all workers.
func NewPool(ctx context.Context, handler Handler, numWorkers int) *Pool {
	messageChan := make(chan []byte, numWorkers)
	pool := &Pool{
		handler:    handler,
		numWorkers: numWorkers,
		wg:         &sync.WaitGroup{},
		messages:   messageChan,
		errChan:    make(chan *worker, numWorkers),
	}
	pool.init(ctx)

	return pool
}

func (p *Pool) init(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		w := &worker{
			id:       i + 1,
			handler:  p.handler,
			wg:       p.wg,
			messages: p.messages,
			errChan:  p.errChan,
		}
		p.wg.Add(1)
		go w.Work(ctx)
	}

	go func() {
		for w := range p.errChan {
			log.Error().Msg(fmt.Sprintf("worker %d exit with err: %s", w.id, w.err.Error()))
			w.err = nil
			time.Sleep(2 * time.Second)
			log.Info().Msg(fmt.Sprintf("Restart worrker id %d", w.id))
			go w.Work(ctx)
		}
	}()
}

// Consume sends a message to the pool; one of the workers will process it.
func (p *Pool) Consume(msg []byte) {
	p.messages <- msg
}

// Close stops the pool: closes the message channel, waits for all workers
// to finish, then closes the error channel. Call after no more Consume calls.
func (p *Pool) Close() {
	close(p.messages)
	p.wg.Wait()       
	close(p.errChan)  
	log.Debug().Msg("All workers have been closed")
}
