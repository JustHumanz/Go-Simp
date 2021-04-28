package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Defines its Prometheus metrics variables
var (
	GetFanArt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_fanart",
			Help: "Get fanart data",
		},
		[]string{"vtuber", "group", "author", "isLewd", "state"},
	)

	GetSubs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_subscriber",
			Help: "Get subscriber/follower count",
		},
		[]string{"vtuber", "group", "state"},
	)

	GetViews = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_viewers",
			Help: "Get viewers count",
		},
		[]string{"vtuber", "group", "state"},
	)

	GetLive = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_live",
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
