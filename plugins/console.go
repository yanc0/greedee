package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/yanc0/greedee/collectd"
	"log"
)

type ConsolePlugin struct {
	ConsolePluginConfig *ConsolePluginConfig
}

type ConsolePluginConfig struct {
	Active bool `yaml:"active"`
	Json   bool `yaml:"json"`
}

func NewConsolePlugin(config *ConsolePluginConfig) *ConsolePlugin {
	return &ConsolePlugin{
		ConsolePluginConfig: config,
	}
}

func (console *ConsolePlugin) Name() string {
	return "Console"
}

func (console *ConsolePlugin) Init() error {
	return nil
}

func (console *ConsolePlugin) Send(cMetrics []*collectd.CollectDMetric) error {
	for _, cMetric := range cMetrics {

		// If json mode is disabled, print metric identifier instead
		if !console.ConsolePluginConfig.Json {
			identifier, err := cMetric.Identifier()
			if err != nil {
				log.Println("[WARN] Console:", err.Error())
			} else {
				fmt.Println("Console Plugin:", identifier, cMetric.Values)
			}
		} else { // JSON mode
			jsn, err := json.Marshal(cMetric)
			if err != nil {
				log.Println("[WARN] Console:", err.Error())
			} else {
				fmt.Println("Console Plugin:", string(jsn))
			}
		}
	}
	return nil
}
