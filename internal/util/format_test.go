package util

import (
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
	if result != "   5.0 B/s" {
		t.Fail()
	}
}

func TestFormatDataRateBytesDoubleDigit(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 50

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	if result != "  50.0 B/s" {
		t.Fail()
	}
}

func TestFormatDataRateKBytes(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 512

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	if result != "   0.5KB/s" {
		t.Fail()
	}
}

func TestFormatDataRateMBytes(t *testing.T) {
	// GIVEN
	duration := 1 * time.Second
	var bytes int64 = 512 * 1000

	// WHEN
	result := FormatDataRate(bytes, duration)

	// THEN
	if result != " 500.0KB/s" {
		t.Fail()
	}
}
