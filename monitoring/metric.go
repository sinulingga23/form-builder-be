package monitoring

import (
	"errors"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	MetricType = int
)

const (
	MetricTotalRequestEndpoint    MetricType = 1
	MetricDurationRequestEndpoint MetricType = 2
)

type Metric struct {
	TotalRequestEndpoint    *prometheus.CounterVec
	DurationRequestEndpoint *prometheus.HistogramVec
}

func NewMetric(registry prometheus.Registerer) *Metric {

	m := &Metric{
		TotalRequestEndpoint: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "form_builder_be",
			Name:      "total_request_endpoint",
			Help:      "It's show total request for each endpoint",
		}, []string{"service_name", "http_method", "http_status"}),
		DurationRequestEndpoint: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "form_builder_be",
			Name:      "duration_request_endpoint",
			Help:      "It's show duration request for each endpoint",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"service_name", "http_method", "http_status"}),
	}

	registry.MustRegister(m.TotalRequestEndpoint, m.DurationRequestEndpoint)
	return m
}

/*
Capture all metrics of Metric object
*/
func (m *Metric) CaptureMetrics(
	serviceName string,
	httpMethod string,
	httpStatus string,
	now *time.Time,
	onlyMetrics ...MetricType) error {

	// only selected metrics will be captured
	if len(onlyMetrics) > 0 {
		for _, metricType := range onlyMetrics {
			switch metricType {
			case MetricTotalRequestEndpoint:
				m.TotalRequestEndpoint.WithLabelValues(serviceName, httpMethod, httpStatus).Inc()
				break
			case MetricDurationRequestEndpoint:
				if now == nil || *now == (time.Time{}) {
					return errors.New("Error: 'now' can't empty for MetricDurationRequestEndpoint")
				}
				m.DurationRequestEndpoint.WithLabelValues(serviceName, httpMethod, httpStatus).Observe(float64(time.Since(*now).Nanoseconds()))
				break
			default:
				log.Printf("Error: Unknown MetricType.")
				break
			}
		}
		return nil
	}

	if now == nil || *now == (time.Time{}) {
		return errors.New("Error: 'now' can't empty for MetricDurationRequestEndpoint")
	}

	m.TotalRequestEndpoint.WithLabelValues(serviceName, httpMethod, httpStatus).Inc()
	m.DurationRequestEndpoint.WithLabelValues(serviceName, httpMethod, httpStatus).Observe(float64(time.Since(*now).Nanoseconds()))
	return nil
}
