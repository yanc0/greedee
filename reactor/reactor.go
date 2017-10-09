package reactor

import (
	"fmt"
	"github.com/yanc0/greedee/events"
	"github.com/yanc0/greedee/plugins"
	"log"
	"time"
)

type Reactor struct {
	EventPlugin plugins.EventPlugin
}

func (r *Reactor) Launch() {
	c := time.Tick(time.Second * 5)
	for range c {
		evts, err := r.EventPlugin.GetExpiredAndNotProcessed()
		if err != nil {
			log.Println("[WARN]", err.Error())
		} else {
			for _, e := range evts {
				// Make new Event for expired event
				newEv := events.Event{
					Source:         "greedee reactor",
					CreatedAt:      time.Now(),
					AuthUserSource: e.AuthUserSource,
					Name:           e.Name,
					Description:    fmt.Sprintf("%s (%s)", e.Description, "expired"),
				}
				newEv.Fail()
				newEv.Gen256Sum()
				r.EventPlugin.Send(newEv)
				r.EventPlugin.Process(e, true) // tag event as processed in plugin
				log.Println("[INFO] Event created for", e.Name, "expired at", e.ExpiresAt)
			}
		}
	}
}
