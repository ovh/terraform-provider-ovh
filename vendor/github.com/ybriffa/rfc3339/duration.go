package rfc3339

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Minute = 60
	Hour   = Minute * 60
	Day    = Hour * 24
	Week   = Day * 7
)

var (
	rfc3339regexp = regexp.MustCompile(`^-?P(\d+Y)?(\d+M)?(\d+W)?(\d+D)?(T(\d+H)?(\d+M)?(\d+S)?)?$`)
)

type threshold struct {
	unit  string
	limit int64
}

func computeDuration(v int64, thresholds []threshold) (string, int64) {
	var result string
	for _, t := range thresholds {

		var remainder int64
		if t.limit == -1 {
			remainder = v
		} else {
			remainder = v % t.limit
			v = v / t.limit
		}
		if remainder > 0 {
			result = fmt.Sprintf("%d%s%s", remainder, t.unit, result)
		}
	}
	return result, v
}

// FormatDuration returns a textual representation of the time.Duration according to the RFC3339.
func FormatDuration(d time.Duration) string {
	return FormatSeconds(int64(d / time.Second))
}

// FormatSeconds returns a textual representation of the number of seconds according to the RFC3339.
func FormatSeconds(seconds int64) string {
	prefix := "P"
	if seconds < 0 {
		prefix = "-P"
		seconds = -seconds
	}

	timeResult, remaining := computeDuration(seconds, []threshold{{"S", 60}, {"M", 60}, {"H", 24}})
	if timeResult != "" {
		timeResult = "T" + timeResult
	}

	// XXX: only take the days since it's complicated for the month, and I currently don't think we'll need it
	result, _ := computeDuration(remaining, []threshold{{"D", 7}, {"W", -1}})
	if timeResult != "" {
		result += timeResult
	}

	// XXX: H4X
	if result == "" {
		result = "T0S"
	}

	return fmt.Sprintf("%s%s", prefix, result)
}

// ParseDuration parses a formatted string to a RFC3339 duration and returns the time duration it represents.
func ParseDuration(s string) (time.Duration, error) {
	seconds, err := ParseSeconds(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

// ParseSeconds parses a formatted string to a RFC3339 duration and returns the time duration it represents.
func ParseSeconds(s string) (int64, error) {
	if !rfc3339regexp.MatchString(s) {
		return 0, fmt.Errorf("%s does not match RFC3339 duration", s)
	}

	multiplier := 1
	if strings.HasPrefix(s, "-") {
		multiplier = -1
		s = strings.TrimPrefix(s, "-")
	}

	matches := rfc3339regexp.FindStringSubmatch(s)
	if len(matches) != 9 {
		return 0, fmt.Errorf("matching number elements must be 9, got %d", len(matches))
	}
	// A bit dirty, but according to the regexp, all the values must be set
	durationResult, err := parsePeriod(matches[1:5])
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %s", err)
	}
	timeResult, err := parseTime(matches[6:9])
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %s", err)
	}

	return (durationResult + timeResult) * int64(multiplier), nil
}

func parsePeriod(array []string) (int64, error) {
	var result int64

	for _, v := range array {
		switch {
		case strings.HasSuffix(v, "W"):
			i, err := strconv.ParseInt(strings.TrimSuffix(v, "W"), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of weeks: %s", err)
			}
			result += (i * Week)

		case strings.HasSuffix(v, "D"):
			i, err := strconv.ParseInt(strings.TrimSuffix(v, "D"), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of days: %s", err)
			}
			result += (i * Day)
		}
	}

	return result, nil
}

func parseTime(array []string) (int64, error) {
	var result int64

	for _, v := range array {
		switch {
		case strings.HasSuffix(v, "H"):
			i, err := strconv.ParseInt(strings.TrimSuffix(v, "H"), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of hours: %s", err)
			}
			result += (i * Hour)

		case strings.HasSuffix(v, "M"):
			i, err := strconv.ParseInt(strings.TrimSuffix(v, "M"), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of minutes: %s", err)
			}
			result += (i * Minute)

		case strings.HasSuffix(v, "S"):
			i, err := strconv.ParseInt(strings.TrimSuffix(v, "S"), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of seconds: %s", err)
			}
			result += i
		}
	}
	return result, nil
}

func ToStringDuration(s string) (string, error) {
	seconds, err := ParseSeconds(s)
	if err != nil {
		return "", err
	}

	switch {
	case seconds%(Minute) > 0:
		return fmt.Sprintf("%ds", seconds), nil
	case seconds%(Hour) > 0:
		return fmt.Sprintf("%dm", seconds/Minute), nil
	case seconds%(Day) > 0:
		return fmt.Sprintf("%dh", seconds/Hour), nil
	}

	return fmt.Sprintf("%dd", seconds/Day), nil
}
