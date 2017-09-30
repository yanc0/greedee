package transformer

import (
	"github.com/yanc0/greedee/collectd"
	"github.com/yanc0/greedee/plugins"
	"github.com/yanc0/greedee/utils"
	"log"
	"math"
)

// Transformer precalculate data of certain types like derive or counter.
// Before sending data to plugin, we have to apply certain formula and
// store previous metrics.
type Transformer struct {
	store plugins.StorePlugin
}

// NewTransformer return a transformer
func NewTransformer(st plugins.StorePlugin) *Transformer {
	return &Transformer{
		store: st,
	}
}

// TransformMetrics apply tranformation on all metrics in a slice
// which have to be transformed
func (tr *Transformer) TransformMetrics(metrics []*collectd.CollectDMetric) {
	for _, m := range metrics {
		if isTransformable(*m) {
			tr.Transform(m)
		}
	}
}

// Tranform get old metric, store current data for later use and
// apply formula on metric values
func (tr *Transformer) Transform(m *collectd.CollectDMetric) {
	metricID, err := m.IdentifierSHA1Sum()
	if err != nil {
		log.Println("[WARN] Metric error:", err)
	}
	oldMetric := tr.store.Get(metricID)

	// Now oldMetric is saved, we can store the untouched
	// new metric, overwriting the old one
	err = tr.store.Put(metricID, m.Clone())
	if err != nil {
		log.Println("[WARN] Metric store failed ", err.Error())
	}

	// Transform value only if the old one is stored
	if oldMetric != nil {
		for i, dsType := range m.DSTypes {
			// Select the right formula for the dsType
			switch dsType {
			case "derive":
				m.Values[i] = derive(oldMetric.Values[i], m.Values[i], oldMetric.Time, m.Time)
			case "counter":
				m.Values[i] = counter(oldMetric.Values[i], m.Values[i], oldMetric.Time, m.Time)
			}
		}
		return
	}

	// We can't calculate derive or counter because we don't
	// have the previous values. Reset it to 0
	for i, _ := range m.Values {
		m.Values[i] = 0
	}

}

// isTransformable returns true if metric have to be transformed
func isTransformable(m collectd.CollectDMetric) bool {
	return utils.ContainsString(m.DSTypes, "derive") || utils.ContainsString(m.DSTypes, "counter")
}

// https://collectd.org/wiki/index.php/Data_source
func derive(valueOld float64, valueNew float64, timeOld float64, timeNew float64) float64 {
	return (valueNew - valueOld) / (timeNew - timeOld)
}

// https://collectd.org/wiki/index.php/Data_source
func counter(valueOld float64, valueNew float64, timeOld float64, timeNew float64) float64 {
	var width float64
	if valueOld < math.Pow(2, 32) {
		width = 32
	} else {
		width = 64
	}
	return (math.Pow(2, width) - valueOld + valueNew) / (timeNew - timeOld)
}
