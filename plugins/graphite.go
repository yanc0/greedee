package plugins

import (
	"fmt"
	"math"
	"github.com/yanc0/collectd-http-server/collectd"
	"github.com/marpaia/graphite-golang"
	"errors"
	"time"
)


type PluginGraphite struct {

}

func (graphite *PluginGraphite) Send(cMetrics []collectd.CollectDMetric) error {
	for _, cMetric := range cMetrics {
		gMetrics, err := fromCollectDMetric(cMetric)
		if err != nil {
			return err
		}
		for _, gMetric := range gMetrics {
			gMetric.print()
		}
	}
	return nil
}

type graphiteMetric struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
	Timestamp time.Time `json:"time"`
}

func (gMetric graphiteMetric) print() {
	fmt.Println(gMetric.Name, gMetric.Value, gMetric.Timestamp.Unix())
}

func fromCollectDMetric(cMetric collectd.CollectDMetric) ([]graphiteMetric, error) {
	metrics := make([]graphiteMetric, len(cMetric.Values))

	if cMetric.Host == "" || cMetric.Plugin == "" || cMetric.Type == "" {
		return nil, errors.New("Graphite Plugin: Invalid Collectd Metric")
	}
	ident := cMetric.Host
	ident = ident + "." + cMetric.Plugin
	if cMetric.PluginInstance != "" {
		ident = ident + "." + cMetric.PluginInstance
	}
	ident = ident + "." + cMetric.Type
	if cMetric.TypeInstance != "" {
		ident = ident + "." + cMetric.TypeInstance
	}
	for i, dsName := range cMetric.DSNames {
		valueName := dsName
		metrics[i].Name = ident
		if dsName != "value" && dsName != "" {
			metrics[i].Name = ident + "." + valueName
		}
		metrics[i].Value = cMetric.Values[i]
		metrics[i].Timestamp = time.Unix(round(cMetric.Time), 0)
	}
	return metrics, nil
}


func round(f float64) int64 {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int64(f + math.Copysign(0.5, f))
}

