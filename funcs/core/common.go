package core

//NewMetricValue decorate metric object,return new metric with tags
func NewMetricValue(endpoint, metric string, val interface{}, dataType string, tags map[string]string) *MetricValue {
	mv := MetricValue{
		Endpoint:     endpoint,
		Metric:       metric,
		ValueUntyped: val,
		CounterType:  dataType,
		TagsMap:      map[string]string{},
	}

	for k, v := range tags {
		mv.TagsMap[k] = v
	}

	return &mv
}

//GaugeValue Gauge type
func GaugeValue(endpoint, metric string, val interface{}, tags map[string]string) *MetricValue {
	return NewMetricValue(endpoint, metric, val, GAUGE, tags)
}

//CounterValue Gauge type
func CounterValue(endpoint, metric string, val interface{}, tags map[string]string) *MetricValue {
	return NewMetricValue(endpoint, metric, val, COUNTER, tags)
}
