package plugins

import (
	"github.com/yanc0/greedee/collectd"
	"log"
)

type MemStorePluginConfig struct {
	Active bool `yaml:"active"`
}

// In-memory store for collectd metrics
type MemStorePlugin struct {
	config  MemStorePluginConfig
	metrics map[string]collectd.CollectDMetric
}

// NewMemStore return initialized MemStore
func NewMemStorePlugin(config MemStorePluginConfig) *MemStorePlugin {
	return &MemStorePlugin{
		config:  config,
		metrics: make(map[string]collectd.CollectDMetric),
	}
}

// Save collectd metric in memory
func (ms *MemStorePlugin) Put(id string, metric collectd.CollectDMetric) error {
	id, err := metric.IdentifierSHA1Sum()
	if err != nil {
		log.Println("[ERR] Put metric in store failed:", err)
		return err
	}
	ms.metrics[id] = metric
	return nil
}

// Get last saved metric in memory
func (ms *MemStorePlugin) Get(id string) *collectd.CollectDMetric {
	m, ok := ms.metrics[id]
	if !ok {
		return nil
	}
	return &m
}
