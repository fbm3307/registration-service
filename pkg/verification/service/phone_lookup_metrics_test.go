package service

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	promtestutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterPhoneLookupMetrics(t *testing.T) {
	// given
	reg := prometheus.NewRegistry()

	// when
	require.NotPanics(t, func() {
		RegisterPhoneLookupMetrics(reg)
	})

	// then — observe once so Gather includes the vectors
	PhoneLookupTotal.WithLabelValues("allowed", "low").Inc()
	PhoneLookupErrorsTotal.WithLabelValues("api_error").Inc()

	gathered, err := reg.Gather()
	require.NoError(t, err)
	names := make(map[string]bool, len(gathered))
	for _, m := range gathered {
		names[m.GetName()] = true
	}
	assert.True(t, names["sandbox_signup_phone_lookup_total"])
	assert.True(t, names["sandbox_signup_phone_lookup_errors_total"])
}

func TestPhoneLookupTotalIncrements(t *testing.T) {
	// given — use a fresh CounterVec so tests don't share global state
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sandbox_signup_phone_lookup_total_test",
			Help: "test",
		},
		[]string{"result", "risk_category"},
	)

	// when
	total.WithLabelValues("allowed", "low").Inc()
	total.WithLabelValues("blocked", "high").Inc()
	total.WithLabelValues("blocked", "high").Inc()

	// then
	assert.Equal(t, float64(1), promtestutil.ToFloat64(total.WithLabelValues("allowed", "low")))
	assert.Equal(t, float64(2), promtestutil.ToFloat64(total.WithLabelValues("blocked", "high")))
	assert.Equal(t, 2, promtestutil.CollectAndCount(total))
}

func TestPhoneLookupErrorsTotalIncrements(t *testing.T) {
	// given
	errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sandbox_signup_phone_lookup_errors_total_test",
			Help: "test",
		},
		[]string{"error_type"},
	)

	// when
	errors.WithLabelValues("api_error").Inc()
	errors.WithLabelValues("api_error").Inc()
	errors.WithLabelValues("timeout").Inc()

	// then
	assert.Equal(t, float64(2), promtestutil.ToFloat64(errors.WithLabelValues("api_error")))
	assert.Equal(t, float64(1), promtestutil.ToFloat64(errors.WithLabelValues("timeout")))
	assert.Equal(t, 2, promtestutil.CollectAndCount(errors))
}
