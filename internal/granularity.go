package internal

import (
	"github.com/thediveo/enumflag/v2"
)

type Granularity enumflag.Flag

const (
	Day Granularity = iota
	Week
	Month
	Quarter
	Year
)

var GranularityIds = map[Granularity][]string{
	Day:     {"day"},
	Week:    {"week"},
	Month:   {"month"},
	Quarter: {"quarter"},
	Year:    {"year"},
}
