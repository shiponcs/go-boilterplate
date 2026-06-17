package utils

import (
	"math"
	"strings"
	"time"
)

func GetCurrentTimeUtc() time.Time {
	return time.Now().UTC()
}

func ToSnakeCase(str string) string {
	parts := strings.Fields(str)
	return strings.ToLower(strings.Join(parts, "_"))
}

// RoundTo2 rounds a float to 2 decimal places.
func RoundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}
