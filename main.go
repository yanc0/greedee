package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/yanc0/greedee/plugins"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var metricPluginList 	    []plugins.MetricPlugin
var eventPluginList []plugins.EventPlugin
var config Config

type BasicAuth struct {
	Active   bool     `toml:"active"`
	Accounts []string `toml:"accounts"`
}

type Config struct {
	Listen         string                        `toml:"listen"`
	Port           int                           `toml:"port"`
	BasicAuth      *BasicAuth                    `toml:"basic_auth"`
	GraphitePlugin *plugins.GraphitePluginConfig `toml:"graphite_plugin"`
	ConsolePlugin  *plugins.ConsolePluginConfig  `toml:"console_plugin"`
	MySQLPlugin    *plugins.MySQLPluginConfig    `toml:"mysql_plugin"`
}

func loadConfig(configPath string) {
	configPath = os.ExpandEnv(configPath)

	configStr, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = toml.Unmarshal(configStr, &config) //global config
	if err != nil {
		log.Fatalln(err.Error())
	}

	if config.Listen == "" {
		config.Listen = "127.0.0.1"
	}

	if config.Port == 0 {
		config.Port = 9223
	}
}

func loadPlugins(config *Config) {
	// Metrics Plugins
	if config.GraphitePlugin != nil && config.GraphitePlugin.Active {
		metricPluginList = append(metricPluginList, plugins.NewGraphitePlugin(config.GraphitePlugin))
	}
	if config.ConsolePlugin != nil && config.ConsolePlugin.Active {
		metricPluginList = append(metricPluginList, plugins.NewConsolePlugin(config.ConsolePlugin))
	}

	//Events Plugins
	if config.MySQLPlugin != nil && config.MySQLPlugin.Active {
		eventPluginList = append(eventPluginList, plugins.NewMySQLPlugin(config.MySQLPlugin))
	}
	if len(metricPluginList) < 1 {
		log.Println("[WARN] No plugins loaded")
	} else {
		log.Println("[INFO]", len(metricPluginList) + len(eventPluginList), "Plugins loaded")
	}
}

func initPlugins() {
	for _, p := range metricPluginList {
		err := p.Init()
		if err != nil {
			log.Println("[WARN]", p.Name(), err.Error())
		} else {
			log.Println("[INFO]", p.Name(), "Plugin Initialized")
		}

	}

	for _, p := range eventPluginList {
		err := p.Init()
		if err != nil {
			log.Println("[WARN]", p.Name(), err.Error())
		} else {
			log.Println("[INFO]", p.Name(), "Plugin Initialized")
		}

	}

}

func main() {
	configPath := flag.String("config",
		"/etc/greedee/config.toml",
		"Config path")
	flag.Parse()

	loadConfig(*configPath)
	loadPlugins(&config)
	initPlugins()

	listen := fmt.Sprintf("%s:%d", config.Listen, config.Port)

	http.HandleFunc("/metrics", auth(handlerMetricPost))
	http.HandleFunc("/events", auth(handlerEventPost))
	log.Fatal(http.ListenAndServe(listen, nil))
}
