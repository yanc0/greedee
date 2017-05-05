package plugins

import (
	"fmt"
	"github.com/yanc0/collectd-http-server/collectd"
	gr "github.com/marpaia/graphite-golang"
	"errors"
	"log"
)


type GraphitePlugin struct {
	Server *gr.Graphite
	Config *GraphitePluginConfig
}

type GraphitePluginConfig struct {
	Active bool `yaml:"active"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Prefix string `yaml:"prefix"`
}

func NewGraphitePlugin(config *GraphitePluginConfig) *GraphitePlugin {
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}

	if config.Port == 0 {
		config.Port = 2003
	}

	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	return &GraphitePlugin{
		Config: config,
	}
}

func (graphite *GraphitePlugin) Name() string {
	return "Graphite"
}

func (graphite *GraphitePlugin) Init() error {
	server, err := gr.GraphiteFactory(graphite.Config.Protocol,
		graphite.Config.Host,
		graphite.Config.Port,
		graphite.Config.Prefix)
	if err != nil {
		log.Println("[WARN] Graphite", err.Error())
	} else {
		log.Println("[INFO] Graphite connection success")
	}
	graphite.Server = server
	log.Println("[INFO] Graphite Plugin Initialized")
	return nil
}

func (graphite *GraphitePlugin) Send(cMetrics []collectd.CollectDMetric) error {
	var toSend []gr.Metric

	if graphite.Server == nil{
		log.Println("[WARN] Graphite is not connected, retrying...")

		server, err := gr.NewGraphite("127.0.0.1", 2003)
		if err != nil {
			return err
		}
		graphite.Server = server
		log.Println("[INFO] Graphite connection succeed")
	}


	for _, cMetric := range cMetrics {
		gMetrics, err := fromCollectDMetric(cMetric)
		if err != nil {
			log.Println("[WARN]", err.Error())
		} else {
			for _, m := range gMetrics{
				toSend = append(toSend, m)
			}
		}
	}

	err := graphite.Server.SendMetrics(toSend)
	if err != nil {
		log.Println("[WARN]", err.Error())
		err = graphite.Server.Connect()
		if err == nil {
			log.Println("[INFO] Graphite connection success")
			graphite.Server.SendMetrics(toSend)
		}
	}

	return nil
}

func fromCollectDMetric(cMetric collectd.CollectDMetric) ([]gr.Metric, error) {
	metrics := make([]gr.Metric, len(cMetric.Values))

	if cMetric.Host == "" || cMetric.Plugin == "" || cMetric.Type == "" {
		return nil, errors.New("Graphite Plugin: Invalid Collectd Metric")
	}
	ident := cMetric.Host
	ident = ident + "." + cMetric.Plugin
	if cMetric.PluginInstance != "" {
		ident = ident + "." + cMetric.PluginInstance
	}
	ident = ident + "." + cMetric.Type
	if cMetric.TypeInstance != "" {
		ident = ident + "." + cMetric.TypeInstance
	}
	for i, dsName := range cMetric.DSNames {
		valueName := dsName
		metrics[i].Name = ident
		if dsName != "value" && dsName != "" {
			metrics[i].Name = ident + "." + valueName
		}
		metrics[i].Value = fmt.Sprintf("%.4f", cMetric.Values[i])
		metrics[i].Timestamp = int64(cMetric.Time)
	}
	return metrics, nil
}