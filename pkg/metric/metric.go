package metric

import (
	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Defines its Prometheus metrics variables
var (
	GetFanArt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: config.Get_Fanart,
			Help: "Get fanart data",
		},
		[]string{"vtuber", "group", "author", "isLewd", "state"},
	)

	GetSubs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: config.Get_Subscriber,
			Help: "Get subscriber/follower count",
		},
		[]string{"vtuber", "group", "state"},
	)

	GetViews = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: config.Get_Viewers,
			Help: "Get viewers count",
		},
		[]string{"vtuber", "group", "state"},
	)

	GetLive = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: config.Get_Live,
			Help: "Get vtuber live count",
		},
		[]string{"vtuber", "group", "state"},
	)
)

// Init will register a Prometheus metrics with the specified variables
func Init() {
	prometheus.MustRegister(
		GetFanArt,
		GetSubs,
		GetViews,
		GetLive,
	)
}
