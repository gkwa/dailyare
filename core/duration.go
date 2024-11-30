package core

import (
	"fmt"
	"strconv"
	"time"
)

func parseDuration(since string) (time.Duration, error) {
	if len(since) < 2 {
		return 0, fmt.Errorf("invalid duration format: %s", since)
	}

	// Check if there are multiple 'd' characters
	dCount := 0
	for _, c := range since {
		if c == 'd' {
			dCount++
		}
	}
	if dCount != 1 {
		return 0, fmt.Errorf("invalid duration unit: must have exactly one 'd'")
	}

	unit := since[len(since)-1:]
	if unit != "d" {
		return 0, fmt.Errorf("invalid duration unit: %s (only 'd' is supported)", unit)
	}

	numStr := since[:len(since)-1]
	days, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse days: %w", err)
	}

	if days < 0 {
		return 0, fmt.Errorf("duration cannot be negative: %d", days)
	}

	duration := fmt.Sprintf("%dh", 24*days)
	return time.ParseDuration(duration)
}
