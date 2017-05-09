package plugins

import "github.com/yanc0/greedee/collectd"

type Plugin interface {
	Send(cMetric []collectd.CollectDMetric) error
	Init() error
	Name() string
}
