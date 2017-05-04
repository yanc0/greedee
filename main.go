package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/yanc0/collectd-http-server/collectd"
	"github.com/yanc0/collectd-http-server/plugins"
)

func main() {
	post, err := ioutil.ReadFile("./example.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var metrics []collectd.CollectDMetric
	err = json.Unmarshal(post, &metrics)
	if err != nil {
		log.Fatal(err)
	}

	graphite := plugins.PluginGraphite{}
	graphite.Send(metrics)
}
