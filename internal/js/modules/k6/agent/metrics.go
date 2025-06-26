package agent

import "go.k6.io/k6/metrics"

type Metrics struct {
	metrics map[string]map[string]metrics.Sink
}

func NewMetrics() *Metrics {
	return &Metrics{metrics: make(map[string]map[string]metrics.Sink)}
}

func (m *Metrics) addSmple(name string, s metrics.Sample) {
	_, nameSeen := m.metrics[name]
	if !nameSeen {
		m.metrics[name] = make(map[string]metrics.Sink)
	}

	metricName := s.Metric.Name
	_, metricSeen := m.metrics[name][metricName]
	if !metricSeen {
		m.metrics[name][metricName] = metrics.NewSink(s.Metric.Type)
	}

	m.metrics[name][metricName].Add(s)
}
