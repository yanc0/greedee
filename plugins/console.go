package plugins

import (
	"fmt"
	"github.com/yanc0/collectd-http-server/collectd"
	"log"
)

type ConsolePlugin struct {
	ConsolePluginConfig *ConsolePluginConfig
}

type ConsolePluginConfig struct {
	Active bool `yaml:"active"`
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
	log.Println("[INFO] Console Plugin Initialized")
	return nil
}

func (console *ConsolePlugin) Send(cMetrics []collectd.CollectDMetric) error {
	for _, cMetric := range cMetrics {
		identifier, err := cMetric.CollectDIdentifier()
		if err != nil {
			log.Println("[WARN] Console:", err.Error())
		} else {
			fmt.Println("Console Plugin:", identifier, cMetric.Values)
		}
	}
	return nil
}
