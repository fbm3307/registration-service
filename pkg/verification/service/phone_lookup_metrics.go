package service

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// PhoneLookupTotal counts successful Twilio Lookup API responses by result and risk category.
	PhoneLookupTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sandbox_signup_phone_lookup_total",
			Help: "Total successful phone lookup operations by result and risk category.",
		},
		[]string{"result", "risk_category"},
	)

	// PhoneLookupErrorsTotal counts Twilio Lookup API failures (fail-open cases).
	PhoneLookupErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sandbox_signup_phone_lookup_errors_total",
			Help: "Phone lookup API errors (fail-open cases).",
		},
		[]string{"error_type"},
	)
)

// RegisterPhoneLookupMetrics registers phone lookup metrics with the given registry.
func RegisterPhoneLookupMetrics(reg *prometheus.Registry) {
	reg.MustRegister(PhoneLookupTotal, PhoneLookupErrorsTotal)
}
