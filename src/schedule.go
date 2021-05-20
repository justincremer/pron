package pron

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Table of time maps for a given job
type schedule struct {
	sec   map[int]struct{}
	min   map[int]struct{}
	hour  map[int]struct{}
	day   map[int]struct{}
	month map[int]struct{}
	dow   map[int]struct{}
}

// Regex patterns for schedule parsing
var (
	matchSpaces = regexp.MustCompile("\\s+")
	matchN      = regexp.MustCompile("(.*)/(\\d+)")
	matchRange  = regexp.MustCompile("^(\\d+)-(\\d+)$")
)

// Parse schedule string and create a schedule or error if syntax is wrong
func parseSchedule(s string) (t schedule, err error) {
	s = matchSpaces.ReplaceAllLiteralString(s, " ")
	parts := strings.Split(s, " ")
	if len(parts) != 6 {
		return schedule{}, errors.New("Schedule string must have six components like * * * * * *")
	}

	t.sec, err = parsePart(parts[0], 0, 59)
	if err != nil {
		return t, err
	}

	t.min, err = parsePart(parts[1], 0, 59)
	if err != nil {
		return t, err
	}

	t.hour, err = parsePart(parts[2], 0, 23)
	if err != nil {
		return t, err
	}

	t.day, err = parsePart(parts[3], 1, 31)
	if err != nil {
		return t, err
	}

	t.month, err = parsePart(parts[4], 1, 12)
	if err != nil {
		return t, err
	}

	t.dow, err = parsePart(parts[5], 0, 6)
	if err != nil {
		return t, err
	}

	//  day/dayOfWeek combination
	switch {
	case len(t.day) < 31 && len(t.dow) == 7: // day set, but not dayOfWeek, clear dayOfWeek
		t.dow = make(map[int]struct{})
	case len(t.dow) < 7 && len(t.day) == 31: // dayOfWeek set, but not day, clear day
		t.day = make(map[int]struct{})
	default:
		// both day and dayOfWeek are * or both are set, use combined
		// i.e. don't do anything here
	}

	return t, nil
}

// parsePart parse individual schedule part from schedule string
func parsePart(s string, min, max int) (map[int]struct{}, error) {
	res := make(map[int]struct{}, 0)

	// wildcard pattern
	if s == "*" {
		for i := min; i <= max; i++ {
			res[i] = struct{}{}
		}
		return res, nil
	}

	// */2 1-59/5 pattern
	if matches := matchN.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := matchRange.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", rng[1], s, rng[1], min, max)
				}
			} else {
				return nil, fmt.Errorf("Unable to parse %s part in %s", matches[1], s)
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			res[i] = struct{}{}
		}
		return res, nil
	}

	// 1,2,4  or 1,2,10-15,20,30-45 pattern
	parts := strings.Split(s, ",")
	for _, x := range parts {
		if rng := matchRange.FindStringSubmatch(x); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", x, s, x, min, max)
			}
			for i := localMin; i <= localMax; i++ {
				res[i] = struct{}{}
			}
		} else if i, err := strconv.Atoi(x); err == nil {
			if i < min || i > max {
				return nil, fmt.Errorf("Out of range for %d in %s. %d must be in range %d-%d", i, s, i, min, max)
			}
			res[i] = struct{}{}
		} else {
			return nil, fmt.Errorf("Unable to parse %s part in %s", x, s)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("Unable to parse %s", s)
	}

	return res, nil
}

// Time object
type tick struct {
	sec   int
	min   int
	hour  int
	day   int
	month int
	dow   int
}

// Gets the current time
func getTick(t time.Time) tick {
	return tick{
		sec:   t.Second(),
		min:   t.Minute(),
		hour:  t.Hour(),
		day:   t.Day(),
		month: int(t.Month()),
		dow:   int(t.Weekday()),
	}
}
