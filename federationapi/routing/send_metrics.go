package routing

import "github.com/prometheus/client_golang/prometheus"

var (
	metricSendTransactionDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_duration",
		},
	)
	metricSendTransactionRxPDUs = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_pdus",
			Buckets:   prometheus.LinearBuckets(0.0, 1.0, 128),
		},
		[]string{
			"result",
		},
	)
	metricSendTransactionRxEDUs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_edus",
			Buckets:   prometheus.LinearBuckets(0.0, 1.0, 128),
		},
	)
)

func init() {
	// Register prometheus metrics. They must be registered to be exposed.
	prometheus.MustRegister(metricSendTransactionDuration)
	prometheus.MustRegister(metricSendTransactionRxPDUs)
	prometheus.MustRegister(metricSendTransactionRxEDUs)
}
