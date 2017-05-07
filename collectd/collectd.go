package collectd

import (
	"errors"
)

type CollectDMetric struct {
	Host           string            `json:"host"`
	Plugin         string            `json:"plugin"`
	PluginInstance string            `json:"plugin_instance"`
	Type           string            `json:"type"`
	TypeInstance   string            `json:"type_instance"`
	Time           float64           `json:"time"`
	Interval       int               `json:"interval"`
	DSTypes        []string          `json:"dstypes"`
	DSNames        []string          `json:"dsnames"`
	Values         []float64         `json:"values"`
	Meta           map[string]string `json:"meta"`
}

func (cMetric *CollectDMetric) CollectDIdentifier() (string, error) {
	if cMetric.Host == "" || cMetric.Plugin == "" || cMetric.Type == "" {
		return "", errors.New("Invalid Collectd Metric")
	}
	ident := cMetric.Host
	ident = ident + "/" + cMetric.Plugin
	if cMetric.PluginInstance != "" {
		ident = ident + "-" + cMetric.PluginInstance
	}
	ident = ident + "/" + cMetric.Type
	if cMetric.TypeInstance != "" {
		ident = ident + "-" + cMetric.TypeInstance
	}
	return ident, nil
}
