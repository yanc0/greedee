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
}
