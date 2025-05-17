package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	URLShortenRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "url_shorten_requests_total",
		Help: "Total number of URL shorten requests",
	})

	URLRedirectRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shorten_redirects_total",
		Help: "Total number of redirects",
	}, []string{"url"},
	)
)

func Register() {
	prometheus.MustRegister(URLShortenRequests, URLRedirectRequests)
}
