package tracing

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func TraceDefaultSegment(ctx context.Context, segmentName string) func() {
	_, f := TraceSegment(ctx, segmentName)
	return f
}

func TraceSegment(ctx context.Context, segmentName string, opts ...func(*newrelic.Transaction)) (*newrelic.Transaction, func()) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return nil, func() {}
	}
	for _, opt := range opts {
		opt(txn)
	}
	segment := txn.StartSegment(segmentName)
	return txn, func() { segment.End() }
}
