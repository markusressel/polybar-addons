package util

import (
	"fmt"
	"math"
	"strings"
)

const (
	BytesInKilobyte = 1024
	BytesInMegabyte = BytesInKilobyte * 1024
	BytesInGigabyte = BytesInMegabyte * 1024
)

func FormatDataRate(bytesPerMilliSecond float64) string {
	decimalPlaces := 1
	unit := " B/s"
	unitFactor := 1.0

	bytesPerSecond := bytesPerMilliSecond / 1000

	if bytesPerSecond >= BytesInKilobyte/10 {
		decimalPlaces = 1
		unitFactor = 1.0 / BytesInKilobyte
		unit = "KB/s"
	}

	if bytesPerSecond >= BytesInMegabyte {
		decimalPlaces = 1
		unitFactor = 1.0 / BytesInMegabyte
		unit = "MB/s"
	}

	if bytesPerSecond >= BytesInGigabyte {
		decimalPlaces = 1
		unitFactor = 1.0 / BytesInGigabyte
		unit = "GB/s"
	}

	// round values appropriately
	var roundingFactor float64
	if decimalPlaces > 0 {
		roundingFactor = float64(10 * decimalPlaces)
	} else {
		roundingFactor = 1
	}

	value := float64(bytesPerSecond) * unitFactor

	value = value * roundingFactor
	value = math.Floor(value)
	value = value / roundingFactor

	prefix := ""
	if value >= 100 {
		prefix = ""
	} else if 10 <= value && value < 100 {
		prefix = " "
	} else {
		prefix = "  "
	}

	formattedValue := betterFormat(value)

	result := fmt.Sprintf("%s %s%s", prefix, formattedValue, unit)

	return result
}

func betterFormat(num float64) string {
	s := fmt.Sprintf("%.4f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}
