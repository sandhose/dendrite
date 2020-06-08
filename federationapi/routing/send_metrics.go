package routing

import "github.com/prometheus/client_golang/prometheus"

var (
	metricSendTransactionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_duration",
			Buckets:   prometheus.LinearBuckets(0.0, 1.0, 30),
		},
	)
	metricSendTransactionRxPDUs = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_pdus",
			Buckets:   []float64{},
		},
		[]string{
			"result",
		},
	)
	metricSendTransactionRxEDUs = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_edus",
			Buckets:   []float64{},
		},
		[]string{
			"result",
		},
	)
)

func init() {
	// Register prometheus metrics. They must be registered to be exposed.
	prometheus.MustRegister(metricSendTransactionDuration)
	prometheus.MustRegister(metricSendTransactionRxPDUs)
	prometheus.MustRegister(metricSendTransactionRxEDUs)
}
