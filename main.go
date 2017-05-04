package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	"math"
	"strings"
)


type CollectDMetric struct {
	Host string `json:"host"`
	Plugin string `json:"plugin"`
	PluginInstance string`json:"plugin_instance"`
	Type string `json:"type"`
	TypeInstance string `json:"type_instance"`
	Time float64 `json:"time"`
	Interval int `json:"interval"`
	DSTypes []string `json:"dstypes"`
	DSNames []string `json:"dsnames"`
	Values []float64 `json:"values"`
	Meta map[string]string `json:"meta"`
}

func main() {
	post, err := ioutil.ReadFile("./example.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var tab []CollectDMetric
	err = json.Unmarshal(post, &tab)
	if err != nil {
		log.Fatal(err)
	}
	PluginGraphite(tab)
}

func PluginGraphite(metrics []CollectDMetric) error {
	for _, metric := range metrics {
		for i, dsName := range metric.DSNames {
			valueName := dsName
			if dsName == "value" {
				valueName = ""
			}
			nameTab := []string{metric.Host, metric.Plugin, metric.TypeInstance, valueName}
			name := graphiteConcat(nameTab)
			value := metric.Values[i]
			time := round(metric.Time)
			fmt.Println(name, value, time)
		}
	}
	return nil
}

func graphiteConcat(names []string) string {
	var join []string
	for i, name := range names {
		if name == "" {
			join = remove(names, i)
		}
	}
	concat := strings.Join(join, ".")
	return concat
}

func round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}