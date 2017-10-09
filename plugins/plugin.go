package plugins

import (
	"github.com/yanc0/greedee/collectd"
	"github.com/yanc0/greedee/events"
)

type MetricPlugin interface {
	Send(cMetric []*collectd.CollectDMetric) error
	Init() error
	Name() string
}

type EventPlugin interface {
	Send(event events.Event) error
	Init() error
	Name() string
	GetExpiredAndNotProcessed() ([]events.Event, error)
	Process(e events.Event, expired bool) error
	ProcessAll(e events.Event) error
}

type StorePlugin interface {
	Put(id string, metric collectd.CollectDMetric) error
	Get(id string) *collectd.CollectDMetric
}
