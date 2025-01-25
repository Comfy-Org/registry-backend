package algolia

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// executor is a worker structure that manages jobs via a buffered channel.
type executor struct {
	jobs chan func()
}

// newExecutor creates a new executor with a specified maximum number of queued jobs.
func newExecutor(maxJob int) *executor {
	return &executor{
		jobs: make(chan func(), maxJob),
	}
}

// start begins processing jobs from the executor's job queue until the context is cancelled.
func (e *executor) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // Exit the loop if the context is cancelled.
			return
		case f := <-e.jobs: // Fetch and execute a job from the channel.
			f()
		}
	}
}

// schedule adds a new job to the executor's job queue for execution.
// If the job queue is full, it logs a warning and adds relevant attributes to the New Relic transaction.
func (e *executor) schedule(ctx context.Context, op string, f func() error) {
	// Retrieve the New Relic transaction from the context.
	txn := newrelic.FromContext(ctx)

	// Wrap the job function with additional logic for New Relic transaction handling.
	fwrap := func() {
		if txn == nil {
			// Execute the job directly if no New Relic transaction is available.
			f()
			return
		}

		// Start a new segment for this operation within the New Relic transaction.
		txnSegment := txn.NewGoroutine().StartSegment(op)
		defer txnSegment.End()

		// Execute the job and handle errors if they occur.
		if err := f(); err != nil {
			log.Ctx(ctx).Error().
				Err(err).
				Str("operation", op).
				Msg("Failed to execute Algolia operation")
			txn.AddAttribute("algolia.op", op)
			txn.AddAttribute("algolia.error", err.Error())
		}
	}

	// Attempt to enqueue the wrapped job function into the job queue.
	select {
	case e.jobs <- fwrap: // Job successfully enqueued.
	default:
		// Handle the case where the job queue is full.
		if txn != nil {
			txn.AddAttribute("algolia.error", "Algolia queue full")
		}
		log.Ctx(ctx).Warn().Msg("Algolia queue full")
	}
}
