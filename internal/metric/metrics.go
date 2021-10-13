package metric

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/mtulio/go-lab-api/internal/event"
)

type MetricsHandler struct {
	Time                time.Time `json:"time"`
	AppTermination      bool      `json:"app_termination"`
	AppHealthy          bool      `json:"app_healthy"`
	TargetHealthy       bool      `json:"tg_healthy"`
	TargetHealthCount   uint8     `json:"tg_health_count"`
	TargetUnhealthCount uint8     `json:"tg_unhealth_count"`
	ReqCountService     uint64    `json:"reqc_service"`
	ReqCountHC          uint64    `json:"reqc_hc"`
	mxReqService        sync.Mutex
	mxReqHC             sync.Mutex
	event               *event.EventHandler
}

func NewMetricHandler(e *event.EventHandler) *MetricsHandler {
	return &MetricsHandler{
		event: e,
	}
}

func (m *MetricsHandler) Set(metric string, value uint8) {
	return
}

func (m *MetricsHandler) Inc(metric string) {
	switch metric {
	case "requests_service":
		m.mxReqService.Lock()
		m.ReqCountService += 1
		m.mxReqService.Unlock()
	case "requests_hc":
		m.mxReqService.Lock()
		m.ReqCountHC += 1
		m.mxReqService.Unlock()
	}
	return
}

// StartPush is a routing to dump/push metrics to
// anywhere (ToDo). Only stdout is supported atm.
func (m *MetricsHandler) StartPusher() {
	for {
		m.Time = time.Now()
		data, err := json.Marshal(m)
		if err != nil {
			log.Println("Error building metrics...")
			time.Sleep(5 * time.Second)
			continue
		}
		m.event.Send("metrics", "metrics-push", string(data))
		time.Sleep(1 * time.Second)
	}
}