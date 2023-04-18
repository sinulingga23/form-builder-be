package monitoring

import "github.com/prometheus/client_golang/prometheus"

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
		}, []string{"service_name", "http_method"}),
		DurationRequestEndpoint: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "form_builder_be",
			Name:      "duration_request_endpoint",
			Help:      "It's show duration request for each endpoint",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"service_name", "http_method"}),
	}

	registry.MustRegister(m.TotalRequestEndpoint, m.DurationRequestEndpoint)
	return m
}
