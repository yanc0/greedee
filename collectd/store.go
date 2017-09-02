package collectd

type Store interface {
	Put(id string, metric CollectDMetric) error
	Get(id string) *CollectDMetric
}

