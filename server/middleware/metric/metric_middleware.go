package drip_metric

import (
	"context"
	"os"
	"registry-backend/config"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	metricpb "google.golang.org/genproto/googleapis/api/metric"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MetricTypePrefix = "custom.googleapis.com/comfy_api_frontend"
	batchInterval    = 5 * time.Minute // Batch interval for sending metrics
)

var (
	environment = os.Getenv("DRIP_ENV")
	metricsCh   = make(chan *monitoringpb.TimeSeries, 1000)
)

func init() {
	go processMetricsBatch()
}

// MetricsMiddleware creates a middleware to capture and send metrics for HTTP requests.
func MetricsMiddleware(client *monitoring.MetricClient, config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			err := next(c)
			endTime := time.Now()

			// Generate metrics for the request duration, count, and errors.
			if config.DripEnv != "localdev" {
				enqueueMetrics(
					createDurationMetric(c, startTime, endTime),
					createRequestMetric(c),
					createErrorMetric(c, err),
				)
			}
			return err
		}
	}
}

// CounterMetric safely increments counters using concurrent maps and atomic operations.
type CounterMetric struct{ sync.Map }

func (c *CounterMetric) increment(key any, i int64) int64 {
	v, loaded := c.LoadOrStore(key, new(atomic.Int64))
	ai := v.(*atomic.Int64)
	if !loaded {
		ai.Add(i) // Initialize and increment atomically
	}
	return ai.Load()
}

// EndpointMetricKey provides a unique identifier for metrics based on request properties.
type EndpointMetricKey struct {
	endpoint   string
	method     string
	statusCode string
}

func endpointMetricKeyFromEcho(c echo.Context) EndpointMetricKey {
	return EndpointMetricKey{
		endpoint:   c.Path(),
		method:     c.Request().Method,
		statusCode: strconv.Itoa(c.Response().Status),
	}
}

func (e EndpointMetricKey) toLabels() map[string]string {
	return map[string]string{
		"endpoint":   e.endpoint,
		"method":     e.method,
		"statusCode": e.statusCode,
		"env":        environment,
	}
}

func enqueueMetrics(series ...*monitoringpb.TimeSeries) {
	for _, s := range series {
		if s != nil {
			metricsCh <- s
		}
	}
}

func processMetricsBatch() {
	ticker := time.NewTicker(batchInterval)
	for range ticker.C {
		sendBatchedMetrics()
	}
}

func sendBatchedMetrics() {
	var series []*monitoringpb.TimeSeries
	for {
		select {
		case s := <-metricsCh:
			series = append(series, s)
			if len(series) >= 1000 {
				sendMetrics(series)
				series = nil
			}
		default:
			if len(series) > 0 {
				sendMetrics(series)
			}
			return
		}
	}
}

func sendMetrics(series []*monitoringpb.TimeSeries) {
	if len(series) == 0 {
		return
	}

	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create metric client")
		return
	}
	defer client.Close()

	req := &monitoringpb.CreateTimeSeriesRequest{
		Name:       "projects/" + os.Getenv("PROJECT_ID"),
		TimeSeries: series,
	}

	if err := client.CreateTimeSeries(ctx, req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create time series")
	}
}

// createDurationMetric constructs a metric for the request processing duration.
func createDurationMetric(c echo.Context, startTime, endTime time.Time) *monitoringpb.TimeSeries {
	key := endpointMetricKeyFromEcho(c)
	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type:   MetricTypePrefix + "/request_duration",
			Labels: key.toLabels(),
		},
		Points: []*monitoringpb.Point{{
			Interval: &monitoringpb.TimeInterval{
				EndTime: timestamppb.New(endTime),
			},
			Value: &monitoringpb.TypedValue{
				Value: &monitoringpb.TypedValue_DoubleValue{
					DoubleValue: endTime.Sub(startTime).Seconds(),
				},
			},
		}},
	}
}

var reqCountMetric = CounterMetric{Map: sync.Map{}}

// createRequestMetric constructs a cumulative metric for counting requests.
func createRequestMetric(c echo.Context) *monitoringpb.TimeSeries {
	key := endpointMetricKeyFromEcho(c)
	val := reqCountMetric.increment(key, 1)
	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type:   MetricTypePrefix + "/request_count",
			Labels: key.toLabels(),
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
	}
}

var reqErrCountMetric = CounterMetric{Map: sync.Map{}}

// createErrorMetric constructs a cumulative metric for counting request errors.
func createErrorMetric(c echo.Context, err error) *monitoringpb.TimeSeries {
	if c.Response().Status < 400 && err == nil {
		return nil // No error occurred, no metric needed
	}

	key := endpointMetricKeyFromEcho(c)
	val := reqErrCountMetric.increment(key, 1)
	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type:   MetricTypePrefix + "/request_errors",
			Labels: key.toLabels(),
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
	}
}
