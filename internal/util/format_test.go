package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatDataRateBytesSingleDigit(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 5

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	assert.Equal(t, "   5.0 B/s", result)
}

func TestFormatDataRateBytesDoubleDigit(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 50

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	assert.Equal(t, "  50.0 B/s", result)
}

func TestFormatDataRateKBytes(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 512

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	assert.Equal(t, "   0.5KB/s", result)
}

func TestFormatDataRateMBytes(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 512 * 1000

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	assert.Equal(t, " 500.0KB/s", result)
}
