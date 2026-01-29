/*

Package duration parses RFC3339 duration strings into time.Duration

Installation

Just go get the package:

    go get -u github.com/peterhellberg/duration

Usage

A small usage example

		package main

		import (
			"fmt"

			"github.com/peterhellberg/duration"
		)

		func main() {
			if d, err := duration.Parse("P1DT30H4S"); err == nil {
				fmt.Println(d) // Output: 54h0m4s
			}
		}

*/
package duration

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// HoursPerDay is the number of hours per day according to Google
	HoursPerDay = 24.0

	// HoursPerWeek is the number of hours per week according to Google
	HoursPerWeek = 168.0

	// HoursPerMonth is the number of hours per month according to Google
	HoursPerMonth = 730.4841667

	// HoursPerYear is the number of hours per year according to Google
	HoursPerYear = 8765.81
)

var (
	// ErrInvalidString is returned when passed an invalid string
	ErrInvalidString = fmt.Errorf("invalid duration string")

	// ErrUnsupportedFormat is returned when parsing fails
	ErrUnsupportedFormat = fmt.Errorf("unsupported duration string format")

	pattern = regexp.MustCompile(`\A(-)?P((?P<years>[\d\.]+)Y)?((?P<months>[\d\.]+)M)?((?P<weeks>[\d\.]+)W)?((?P<days>[\d\.]+)D)?(T((?P<hours>[\d\.]+)H)?((?P<minutes>[\d\.]+)M)?((?P<seconds>[\d\.]+?)S)?)?\z`)

	invalidStrings = []string{"", "P", "PT"}
)

// Parse a RFC3339 duration string into time.Duration
func Parse(s string) (time.Duration, error) {
	if contains(invalidStrings, s) || strings.HasSuffix(s, "T") {
		return 0, ErrInvalidString
	}

	var (
		match  []string
		prefix string
	)

	if pattern.MatchString(s) {
		match = pattern.FindStringSubmatch(s)
	} else {
		return 0, ErrUnsupportedFormat
	}

	if strings.HasPrefix(s, "-") {
		prefix = "-"
	}

	return durationFromMatchAndPrefix(match, prefix)
}

func durationFunc(prefix string) func(string, float64) time.Duration {
	return func(format string, f float64) time.Duration {
		if d, err := time.ParseDuration(fmt.Sprintf(prefix+format, f)); err == nil {
			return d
		}

		return time.Duration(0)
	}
}

func durationFromMatchAndPrefix(match []string, prefix string) (time.Duration, error) {
	d := time.Duration(0)

	duration := durationFunc(prefix)

	for i, name := range pattern.SubexpNames() {
		value := match[i]
		if i == 0 || name == "" || value == "" {
			continue
		}

		if f, err := strconv.ParseFloat(value, 64); err == nil {
			switch name {
			case "years":
				d += duration("%fh", f*HoursPerYear)
			case "months":
				d += duration("%fh", f*HoursPerMonth)
			case "weeks":
				d += duration("%fh", f*HoursPerWeek)
			case "days":
				d += duration("%fh", f*HoursPerDay)
			case "hours":
				d += duration("%fh", f)
			case "minutes":
				d += duration("%fm", f)
			case "seconds":
				d += duration("%fs", f)
			}
		}
	}

	return d, nil
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
