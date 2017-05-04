package plugins

import "github.com/yanc0/collectd-http-server/collectd"

type Plugin interface {
	Send(cMetric []collectd.CollectDMetric) error
}