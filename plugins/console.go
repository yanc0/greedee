package plugins

import (
	"github.com/yanc0/collectd-http-server/collectd"
	"fmt"
	"log"
)

type PluginConsole struct {

}

func (graphite *PluginConsole) Send(cMetrics []collectd.CollectDMetric) error {
	for _, cMetric := range cMetrics {
		identifier, err := cMetric.CollectDIdentifier()
		if err != nil {
			log.Println("[WARN] Console:", err.Error())
		} else {
			fmt.Println("Console Plugin:", identifier)
		}
	}
	return nil
}
