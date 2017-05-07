package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/yanc0/collectd-http-server/plugins"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var pluginList []plugins.Plugin
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
}

func loadPlugins(config *Config) {
	if config.GraphitePlugin != nil && config.GraphitePlugin.Active {
		pluginList = append(pluginList, plugins.NewGraphitePlugin(config.GraphitePlugin))
	}
	if config.ConsolePlugin != nil && config.ConsolePlugin.Active {
		pluginList = append(pluginList, plugins.NewConsolePlugin(config.ConsolePlugin))
	}
	log.Println("[INFO] Plugins loaded")
}

func initPlugins() {
	for _, p := range pluginList {
		err := p.Init()
		if err != nil {
			log.Fatalln("[FATAL]", p.Name(), err.Error())
		}
	}
}

func main() {
	configPath := flag.String("configPath",
		"/etc/collectd-http-server/config.yaml",
		"Config path")
	flag.Parse()

	loadConfig(*configPath)
	loadPlugins(&config)
	initPlugins()

	listen := fmt.Sprintf("%s:%d", config.Listen, config.Port)

	http.HandleFunc("/", auth(handlerMetricPost))
	log.Fatal(http.ListenAndServe(listen, nil))
}
