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
	PhoneLookupTotal.Reset()
	PhoneLookupErrorsTotal.Reset()

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
	// given
	PhoneLookupTotal.Reset()

	// when
	PhoneLookupTotal.WithLabelValues("allowed", "low").Inc()
	PhoneLookupTotal.WithLabelValues("blocked", "high").Inc()
	PhoneLookupTotal.WithLabelValues("blocked", "high").Inc()

	// then
	assert.InDelta(t, float64(1), promtestutil.ToFloat64(PhoneLookupTotal.WithLabelValues("allowed", "low")), 0.01)
	assert.InDelta(t, float64(2), promtestutil.ToFloat64(PhoneLookupTotal.WithLabelValues("blocked", "high")), 0.01)
	assert.Equal(t, 2, promtestutil.CollectAndCount(PhoneLookupTotal))
}

func TestPhoneLookupErrorsTotalIncrements(t *testing.T) {
	// given
	PhoneLookupErrorsTotal.Reset()

	// when
	PhoneLookupErrorsTotal.WithLabelValues("api_error").Inc()
	PhoneLookupErrorsTotal.WithLabelValues("api_error").Inc()
	PhoneLookupErrorsTotal.WithLabelValues("timeout").Inc()

	// then
	assert.InDelta(t, float64(2), promtestutil.ToFloat64(PhoneLookupErrorsTotal.WithLabelValues("api_error")), 0.01)
	assert.InDelta(t, float64(1), promtestutil.ToFloat64(PhoneLookupErrorsTotal.WithLabelValues("timeout")), 0.01)
	assert.Equal(t, 2, promtestutil.CollectAndCount(PhoneLookupErrorsTotal))
}
