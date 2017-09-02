package collectd

import (
	"log"
	"math"
)

type Transformer struct {
	TypesDBPath string
	Store Store
}

func NewTransformer(typesDBPath string, st Store) *Transformer{
	return &Transformer{
		TypesDBPath: typesDBPath,
		Store: st,
	}
}

func (tr *Transformer) TransformMetrics(metrics []*CollectDMetric){
	for _, m := range metrics {
		tr.Transform(m)
	}
}

func (tr *Transformer) Transform(m *CollectDMetric) {
	metricID, err := m.Identifier256Sum()
	if err != nil {
		log.Println("[WARN] Metric error:", err)
	}
	oldMetric := tr.Store.Get(metricID)
	// Now oldMetric is saved, we can store the untouched
	// new metric and replace the old one
	err = tr.Store.Put(metricID, *m)
	if err != nil {
		log.Println("[WARN] Metric store failed ", err.Error())
	}
	// Transform value only if the old one is stored
	if oldMetric != nil {
		for i, dsType := range m.DSTypes {
			switch dsType {
			case "derive":
				m.Values[i] = derive(oldMetric.Values[i], m.Values[i], oldMetric.Time, m.Time)
			case "counter":
				m.Values[i] = counter(oldMetric.Values[i], m.Values[i], oldMetric.Time, m.Time)
			}
		}
	} else {
		// We can't calculate derive because we don't have the previous value
		// Reset it to 0
		for i, dsType := range m.DSTypes {
			switch dsType {
			case "derive":
				m.Values[i] = 0
			case "counter":
				m.Values[i] = 0
			}
		}
	}
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