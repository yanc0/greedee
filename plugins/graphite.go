package plugins

import (
	"errors"
	gr "github.com/marpaia/graphite-golang"
	"github.com/yanc0/collectd-http-server/collectd"
	"log"
	"strconv"
)

type GraphitePlugin struct {
	Server *gr.Graphite
	Config *GraphitePluginConfig
}

type GraphitePluginConfig struct {
	Active bool   `toml:"active"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	Proto  string `toml:"proto"`
	Prefix string `toml:"prefix"`
}

func NewGraphitePlugin(config *GraphitePluginConfig) *GraphitePlugin {
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}

	if config.Port == 0 {
		config.Port = 2003
	}

	if config.Proto == "" {
		config.Proto = "tcp"
	}

	return &GraphitePlugin{
		Config: config,
	}
}

func (graphite *GraphitePlugin) Name() string {
	return "Graphite"
}

func (graphite *GraphitePlugin) Init() error {
	err := graphite.Connect()
	if err != nil {
		log.Println("[WARN] Graphite", err.Error())
	}
	return nil
}

func (graphite *GraphitePlugin) Connect() error {
	if graphite.Config == nil {
		return errors.New("Graphite config not set")
	}
	server, err := gr.GraphiteFactory(graphite.Config.Proto,
		graphite.Config.Host,
		graphite.Config.Port,
		graphite.Config.Prefix)
	if err != nil {
		return err
	} else {
		log.Println("[INFO] Graphite connection succeed")
	}
	graphite.Server = server
	return nil
}

func (graphite *GraphitePlugin) Send(cMetrics []collectd.CollectDMetric) error {
	var toSend []gr.Metric
	if graphite.Server == nil {
		log.Println("[WARN] Graphite is not connected, retrying...")
		err := graphite.Connect()
		if err != nil {
			return err
		}
	}

	for _, cMetric := range cMetrics {
		gMetrics, err := fromCollectDMetric(cMetric)
		if err != nil {
			log.Println("[WARN]", err.Error())
		} else {
			for _, m := range gMetrics {
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

	// Construct graphite identifier according to
	// https://collectd.org/wiki/index.php/Naming_schema
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
		metrics[i].Value = strconv.FormatFloat(cMetric.Values[i], 'f', -1, 64)
		metrics[i].Timestamp = int64(cMetric.Time)
	}
	return metrics, nil
}