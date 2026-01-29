# :watch: duration

[![Build status](https://github.com/peterhellberg/duration/actions/workflows/test.yml/badge.svg)](https://github.com/peterhellberg/duration/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/peterhellberg/duration)](https://goreportcard.com/report/peterhellberg/duration)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/peterhellberg/duration)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/duration/blob/master/LICENSE)

Parse a [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) duration string into `time.Duration`

There are probably a few unsupported edge cases still to be fixed, please help me find them :)

The following constants are used to do the calculations for longer durations:

```
HoursPerDay = 24.0
HoursPerWeek = 168.0
HoursPerMonth = 730.4841667
HoursPerYear = 8765.81
```

Look in the test for examples of both valid and invalid duration strings.

## Installation

    go get -u github.com/peterhellberg/duration

Feel free to copy this package into your own codebase.

## Usage

```go
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
```

## RFC3339 grammar for durations

```
   dur-second        = 1*DIGIT "S"
   dur-minute        = 1*DIGIT "M" [dur-second]
   dur-hour          = 1*DIGIT "H" [dur-minute]
   dur-time          = "T" (dur-hour / dur-minute / dur-second)
   dur-day           = 1*DIGIT "D"
   dur-week          = 1*DIGIT "W"
   dur-month         = 1*DIGIT "M" [dur-day]
   dur-year          = 1*DIGIT "Y" [dur-month]
   dur-date          = (dur-day / dur-month / dur-year) [dur-time]

   duration          = "P" (dur-date / dur-time / dur-week)
```
