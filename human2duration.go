package human2duration

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var unitMap = map[string]time.Duration{
	"s":       time.Second,
	"sec":     time.Second,
	"second":  time.Second,
	"seconds": time.Second,
	"mi":      time.Minute,
	"min":     time.Minute,
	"minute":  time.Minute,
	"minutes": time.Minute,
	"h":       time.Hour,
	"hr":      time.Hour,
	"hour":    time.Hour,
	"hours":   time.Hour,
	"d":       24 * time.Hour,
	"day":     24 * time.Hour,
	"days":    24 * time.Hour,
	"w":       7 * 24 * time.Hour,
	"week":    7 * 24 * time.Hour,
	"weeks":   7 * 24 * time.Hour,
	"m":       30 * 24 * time.Hour,
	"mo":      30 * 24 * time.Hour,
	"month":   30 * 24 * time.Hour,
	"months":  30 * 24 * time.Hour,
	"y":       365 * 24 * time.Hour,
	"year":    365 * 24 * time.Hour,
	"years":   365 * 24 * time.Hour,
}

var fuzzyMap = map[string]time.Duration{
	"half an hour":       30 * time.Minute,
	"an hour and a half": 90 * time.Minute,
	"half a day":         12 * time.Hour,

	"a second": time.Second,
	"a minute": time.Minute,
	"an hour":  time.Hour,
	"a day":    24 * time.Hour,
	"a week":   7 * 24 * time.Hour,
	"a month":  30 * 24 * time.Hour,
	"a year":   365 * 24 * time.Hour,
}

var (
	unitRegex             = regexp.MustCompile(`(?i)([\d.]+)\s*([a-z]+)`)
	goStyleCompactPattern = regexp.MustCompile(`^(?:\d+\s*h\s*\d+\s*m(?:\s*\d+\s*s)?|\d+\s*m\s*\d+\s*s|\d+\s*h\s*\d+\s*m?$)$`)
)

func stripPrefixIgnoreCase(s, prefix string) string {
	if strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix)) {
		return s[len(prefix):]
	}
	return s
}

func stripSuffixIgnoreCase(s, suffix string) string {
	if strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix)) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func ParseSinceDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	input = stripSuffixIgnoreCase(input, " ago")

	duration, err := ParseDuration(input)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return -duration, nil
}

func ParseAfterDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	input = stripPrefixIgnoreCase(input, "in ")
	input = stripPrefixIgnoreCase(input, "after ")

	duration, err := ParseDuration(input)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

func ParseDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)

	// Timestamp input (RFC3339, etc.)
	if ts, err := tryParseTimestamp(input); err == nil {
		return time.Until(ts), nil
	}

	// normalize input
	input = strings.ToLower(input)

	// Fuzzy match
	if d, ok := fuzzyMap[input]; ok {
		return d, nil
	}

	// Go compact format like "1h30m", "2m30s", etc.
	if goStyleCompactPattern.MatchString(strings.ReplaceAll(input, " ", "")) {
		if d, err := time.ParseDuration(strings.ReplaceAll(input, " ", "")); err == nil {
			return d, nil
		}
	}

	// General unit pattern: "2d 3h"
	matches := unitRegex.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	total := time.Duration(0)
	previousUnit := ""
	for _, match := range matches {
		valStr, unit := match[1], match[2]
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", valStr)
		}

		// Special disambiguation: if unit is "m", contextually decide
		if unit == "m" {
			if previousUnit == "h" || previousUnit == "hr" || previousUnit == "hour" || previousUnit == "hours" {
				unit = "min"
			} else {
				unit = "month"
			}
		}

		durUnit, ok := unitMap[unit]
		if !ok {
			return 0, fmt.Errorf("unknown unit: %s", unit)
		}
		total += time.Duration(val * float64(durUnit))
		previousUnit = unit
	}

	return total, nil
}

func tryParseTimestamp(s string) (time.Time, error) {
	s = strings.ToUpper(s)
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("not a timestamp")
}
