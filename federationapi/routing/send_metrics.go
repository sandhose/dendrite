package routing

import "github.com/prometheus/client_golang/prometheus"

var (
	metricSendTransactionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_duration",
		},
	)
	metricSendTransactionRxPDUs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_pdus",
		},
	)
	metricSendTransactionSuccessfulPDUs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_pdus_successful",
		},
	)
	metricSendTransactionFailedPDUs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_pdus_failed",
		},
	)
	metricSendTransactionRxEDUs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "dendrite",
			Subsystem: "federationapi",
			Name:      "transaction_edus",
		},
	)
)

func init() {
	// Register prometheus metrics. They must be registered to be exposed.
	prometheus.MustRegister(metricSendTransactionDuration)
	prometheus.MustRegister(metricSendTransactionRxPDUs)
	prometheus.MustRegister(metricSendTransactionSuccessfulPDUs)
	prometheus.MustRegister(metricSendTransactionFailedPDUs)
	prometheus.MustRegister(metricSendTransactionRxEDUs)
}
