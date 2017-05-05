package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/yanc0/collectd-http-server/collectd"
	"github.com/yanc0/collectd-http-server/plugins"
	"net/http"
	"flag"
	"gopkg.in/yaml.v2"
	"os"
	"fmt"
)

var pluginList []plugins.Plugin

type Config struct {
	Listen string `yaml:"listen"`
	Port int `yaml:"port"`
	GraphitePlugin *plugins.GraphitePluginConfig `yaml:"graphite_plugin"`
	ConsolePlugin  *plugins.ConsolePluginConfig `yaml:"console_plugin"`
}

func loadConfig(configPath string) *Config {
	var config Config
	configPath = os.ExpandEnv(configPath)

	configStr, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = yaml.Unmarshal(configStr, &config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &config
}

func loadPlugins(config *Config) {
	if config.GraphitePlugin != nil && config.GraphitePlugin.Active{
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

	config := loadConfig(*configPath)
	listen := fmt.Sprintf("%s:%d", config.Listen, config.Port)
	gr := config.GraphitePlugin

	fmt.Println(gr)

	loadPlugins(config)
	initPlugins()



	http.HandleFunc("/", handlerMetricPost)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func handlerMetricPost(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST" {
		http.Error(w, "405 Method Not Allowed - POST Only", http.StatusMethodNotAllowed)
		return
	}

	post, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer req.Body.Close()

	var metrics []collectd.CollectDMetric
	err = json.Unmarshal(post, &metrics)
	if err != nil {
		log.Println("[WARN]", err.Error())
	}

	for _, p := range pluginList {
		err := p.Send(metrics)
		if err != nil {
			log.Println("[WARN]", err.Error())
		}
	}
	return
}
