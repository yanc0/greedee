package collectd

import (
	"log"
)

// In-memory store for collectd metrics
type MemStore struct {
	Metrics map[string]CollectDMetric
}


// Save collectd metric in memory
func (ms * MemStore) Put(id string, metric CollectDMetric) error {
	id, err := metric.Identifier256Sum()
	if err != nil {
		log.Println("[ERR] Put metric in store failed:", err)
		return err
	}
	ms.Metrics[id] = metric
	return nil
}

// Get last saved metric in memory
func (ms *MemStore) Get(id string) *CollectDMetric {
	m, ok := ms.Metrics[id]
	if ! ok {
		return nil
	}
	return &m
}