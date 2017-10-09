package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/yanc0/greedee/plugins"
	pluginEvent "github.com/yanc0/greedee/plugins/event"
	pluginMetric "github.com/yanc0/greedee/plugins/metric"
	pluginStore "github.com/yanc0/greedee/plugins/store"
	"github.com/yanc0/greedee/reactor"
	"github.com/yanc0/greedee/transformer"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var metricPluginList []plugins.MetricPlugin
var eventPluginList []plugins.EventPlugin
var storePlugin plugins.StorePlugin
var transform *transformer.Transformer

var config Config

type BasicAuth struct {
	Active   bool     `toml:"active"`
	Accounts []string `toml:"accounts"`
}

type Config struct {
	Listen         string                             `toml:"listen"`
	Port           int                                `toml:"port"`
	BasicAuth      *BasicAuth                         `toml:"basic_auth"`
	GraphitePlugin *pluginMetric.GraphitePluginConfig `toml:"graphite_plugin"`
	ConsolePlugin  *pluginMetric.ConsolePluginConfig  `toml:"console_plugin"`
	MySQLPlugin    *pluginEvent.MySQLPluginConfig     `toml:"mysql_plugin"`
	MemStorePlugin *pluginStore.MemStorePluginConfig  `toml:"memstore_plugin"`
}

func loadConfig(configPath string) {
	t0 := time.Now()
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

	fmt.Println("[INFO] Configuration loaded in", time.Since(t0))
}

func loadPlugins(config *Config) {
	t0 := time.Now()
	// Metrics Plugins
	if config.GraphitePlugin != nil && config.GraphitePlugin.Active {
		metricPluginList = append(metricPluginList, pluginMetric.NewGraphitePlugin(config.GraphitePlugin))
	}
	if config.ConsolePlugin != nil && config.ConsolePlugin.Active {
		metricPluginList = append(metricPluginList, pluginMetric.NewConsolePlugin(config.ConsolePlugin))
	}

	// Events Plugins
	if config.MySQLPlugin != nil && config.MySQLPlugin.Active {
		eventPluginList = append(eventPluginList, pluginEvent.NewMySQLPlugin(config.MySQLPlugin))
	}

	nbPlugins := len(metricPluginList) + len(eventPluginList)
	if nbPlugins < 1 {
		log.Println("[WARN] No plugins loaded")
	} else {
		log.Println("[INFO]", nbPlugins, "Plugins loaded in", time.Since(t0))
	}

	//Store Plugin
	if storePlugin == nil {
		storePlugin = pluginStore.NewMemStorePlugin(*config.MemStorePlugin)
		log.Println("[INFO] Memstore plugin loaded")
	}
}

func initPlugins() {
	t0 := time.Now()
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
			log.Println("[INFO]", p.Name(), "Plugin Initialized in", time.Since(t0))
		}

	}

}

func initReactors() {
	t0 := time.Now()
	for _, ep := range eventPluginList {
		r := reactor.Reactor{
			EventPlugin: ep,
		}
		go r.Launch()
	}
	log.Println("[INFO] Reactors launched in", time.Since(t0))
}

func initTransformer() {
	t0 := time.Now()
	transform = transformer.NewTransformer(storePlugin)
	log.Println("[INFO] Metrics transformer launched in", time.Since(t0))

}

func main() {
	configPath := flag.String("config",
		"/etc/greedee/config.toml",
		"Config path")
	flag.Parse()

	loadConfig(*configPath)
	loadPlugins(&config)
	initPlugins()
	initReactors()
	initTransformer()

	listen := fmt.Sprintf("%s:%d", config.Listen, config.Port)

	http.HandleFunc("/metrics", auth(handlerMetricPost))
	http.HandleFunc("/events", auth(handlerEventPost))
	log.Fatal(http.ListenAndServe(listen, nil))
}
