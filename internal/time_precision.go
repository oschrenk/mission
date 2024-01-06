package internal

import (
	"fmt"
	"time"
)

type TimePrecision int64

const (
	Day TimePrecision = iota
	Week
	Month
	Year
)

func (tp TimePrecision) hours() int64 {
	switch tp {
	case Day:
		duration, _ := time.ParseDuration("24h")
		return int64(duration.Hours())
	case Week:
		duration, _ := time.ParseDuration("168h")
		return int64(duration.Hours())
	case Month:
		duration, _ := time.ParseDuration("744h")
		return int64(duration.Hours())
	case Year:
		duration, _ := time.ParseDuration("8760h")
		return int64(duration.Hours())
	}
	return -1
}

func (tp TimePrecision) String() string {
	switch tp {
	case Day:
		return "Day"
	case Week:
		return "Week"
	case Month:
		return "Month"
	case Year:
		return "Year"
	default:
		return fmt.Sprintf("Unknown (%d)", int(tp))
	}
}

func isInt(val float64) bool {
	return val == float64(int(val))
}

func BuildTimePrecision(d time.Duration) (TimePrecision, error) {
	fullDurationInHours := (d + time.Second).Hours()
	if !isInt(fullDurationInHours) {
		return -1, fmt.Errorf("unsupported duration %s. Must be full hour", d)
	}
	hours := int64(fullDurationInHours)

	if int64(hours/Year.hours()) == 1 {
		return Year, nil
	}

	if int64(hours/Month.hours()) == 1 {
		return Month, nil
	}

	if int64(hours/Week.hours()) == 1 {
		return Week, nil
	}

	if int64(hours/Day.hours()) == 1 {
		return Day, nil
	}

	return -1, fmt.Errorf("unsupported duration %s. Unknown precision", d)

}
