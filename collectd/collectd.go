package collectd

import (
	"errors"
	"crypto/sha256"
	"fmt"
)

type CollectDMetric struct {
	Host           string            `json:"host"`
	Plugin         string            `json:"plugin"`
	PluginInstance string            `json:"plugin_instance"`
	Type           string            `json:"type"`
	TypeInstance   string            `json:"type_instance"`
	Time           float64           `json:"time"`
	Interval       float64           `json:"interval"`
	DSTypes        []string          `json:"dstypes"`
	DSNames        []string          `json:"dsnames"`
	Values         []float64         `json:"values"`
	Meta           map[string]string `json:"meta"`
}

// Generate Metric identifier in SHA256 format
func (cMetric *CollectDMetric) Identifier() (string, error) {
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

// Generate Metric identifier in SHA256 format
func (cMetric *CollectDMetric) Identifier256Sum() (string, error) {
	ident, err := cMetric.Identifier()
	if err != nil {
		return "", err
	}
	eventBytes := []byte(ident)
	h := sha256.New()
	h.Write(eventBytes)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Clone return a deep copy of the CollectDMetric
func (cMetric *CollectDMetric) Clone() CollectDMetric {
	m := CollectDMetric{
		Host: cMetric.Host,
		Plugin: cMetric.Plugin,
		PluginInstance: cMetric.PluginInstance,
		Type: cMetric.Type,
		TypeInstance: cMetric.TypeInstance,
		Time: cMetric.Time,
		Interval: cMetric.Interval,
		DSTypes: make([]string, len(cMetric.DSTypes)),
		DSNames: make([]string, len(cMetric.DSNames)),
		Values: make([]float64, len(cMetric.Values)),
		Meta: cMetric.Meta,
	}
	copy(m.DSTypes, cMetric.DSTypes)
	copy(m.DSNames, cMetric.DSNames)
	copy(m.Values, cMetric.Values)

	return m
}