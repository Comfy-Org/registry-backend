package drip_metric

import (
	"context"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *CounterMetric) Increment(key any, i int64) int64 {
	v, _ := c.LoadOrStore(key, new(atomic.Int64))
	ai, ok := v.(*atomic.Int64)
	if !ok {
		ai = new(atomic.Int64)
	}
	ai.Add(i) // Initialize and increment atomically
	return ai.Load()
}

type customCounterKey struct {
	t string
	l string
}

type CustomCounterIncrement struct {
	Type   string
	Labels map[string]string
	Val    int64
}

func (c CustomCounterIncrement) key() customCounterKey {
	keys := make([]string, 0, len(c.Labels))
	for k, v := range c.Labels {
		keys = append(keys, k+":"+v)
	}
	sort.Strings(keys)
	return customCounterKey{
		t: c.Type,
		l: strings.Join(keys, ","),
	}
}

var customCounterCtxKey = struct{}{}

func AttachCustomCounterMetric(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, customCounterCtxKey, &([]CustomCounterIncrement{}))
	return ctx
}

func IncrementCustomCounterMetric(ctx context.Context, inc CustomCounterIncrement) {
	v := ctx.Value(customCounterCtxKey)
	cc, ok := v.(*[]CustomCounterIncrement)
	if !ok || cc == nil {
		return
	}
	*cc = append(*cc, inc)
}

var customCounterMetric = CounterMetric{Map: sync.Map{}}

func CreateCustomCounterMetrics(ctx context.Context) (ts []*monitoringpb.TimeSeries) {
	v := ctx.Value(customCounterCtxKey)
	cc, ok := v.(*[]CustomCounterIncrement)
	if !ok || cc == nil {
		return
	}

	for _, c := range *cc {
		val := customCounterMetric.Increment(c.key(), c.Val)
		ts = append(ts, &monitoringpb.TimeSeries{
			Metric: &metricpb.Metric{
				Type:   MetricTypePrefix + "/" + c.Type,
				Labels: c.Labels,
			},
			MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
			Points: []*monitoringpb.Point{
				{
					Interval: &monitoringpb.TimeInterval{
						StartTime: timestamppb.New(time.Now().Add(-time.Second)),
						EndTime:   timestamppb.New(time.Now()),
					},
					Value: &monitoringpb.TypedValue{
						Value: &monitoringpb.TypedValue_Int64Value{
							Int64Value: val,
						},
					},
				},
			},
		})
	}
	return
}
