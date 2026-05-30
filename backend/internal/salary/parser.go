package salary

import (
	"regexp"
	"strconv"
	"strings"
)

// ParsedSalary holds the extracted min and max salary values in USD
type ParsedSalary struct {
	Min *int64
	Max *int64
}

// Parse extracts salary values from a salary string like "120k - 150k" or "$80,000"
// Returns min and max as int64 pointers (nil if not found)
func Parse(salaryStr string) ParsedSalary {
	if salaryStr == "" {
		return ParsedSalary{}
	}

	// Remove common currency symbols and whitespace
	cleaned := strings.ToLower(salaryStr)
	cleaned = strings.TrimSpace(cleaned)

	// Regex to find numbers followed by optional multipliers (k, m)
	// Matches patterns like: "120k", "80,000", "150.5k", "1.2m"
	re := regexp.MustCompile(`(\d+(?:[.,]\d+)*)\s*([km])?`)
	matches := re.FindAllStringSubmatch(cleaned, -1)

	if len(matches) == 0 {
		return ParsedSalary{}
	}

	var values []int64

	for _, match := range matches {
		numberStr := match[1]
		multiplier := match[2]

		// Replace comma with dot for consistent parsing
		numberStr = strings.Replace(numberStr, ",", ".", -1)

		// Parse as float first to handle decimals
		val, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			continue
		}

		// Apply multiplier
		switch multiplier {
		case "k":
			val *= 1000
		case "m":
			val *= 1000000
		}

		values = append(values, int64(val))
	}

	if len(values) == 0 {
		return ParsedSalary{}
	}

	// If we have 2 values, it's likely a range (min and max)
	// If we have 1 value, use it as the minimum
	if len(values) >= 2 {
		min := values[0]
		max := values[len(values)-1]

		// Ensure min <= max
		if min > max {
			min, max = max, min
		}

		return ParsedSalary{
			Min: &min,
			Max: &max,
		}
	}

	// Single value - use as minimum
	min := values[0]
	return ParsedSalary{
		Min: &min,
	}
}