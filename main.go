package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/yanc0/collectd-http-server/collectd"
	"github.com/yanc0/collectd-http-server/plugins"
)

var pluginList []plugins.Plugin

func loadPlugins() {
	pluginList = append(pluginList, &plugins.PluginGraphite{})
	pluginList = append(pluginList, &plugins.PluginConsole{})
	log.Println("Pluging loaded")
}

func main() {
	loadPlugins()

	post, err := ioutil.ReadFile("./example.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var metrics []collectd.CollectDMetric
	err = json.Unmarshal(post, &metrics)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range pluginList {
		go p.Send(metrics)
	}
}
