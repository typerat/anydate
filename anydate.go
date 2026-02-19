package anydate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	dateFormats = []string{
		"2.1.2006",
		"2006-1-2",
		"2006.1.2",
		"Jan 2, 2006",
		"2006:01:02",
	}
	monthNames = map[string]int{
		"jan": 1,
		"jän": 1,
		"feb": 2,
		"mär": 3,
		"mar": 3,
		"apr": 4,
		"mai": 5,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"okt": 10,
		"oct": 10,
		"nov": 11,
		"dez": 12,
		"dec": 12,
	}
)

// Parse scans a string and returns a list of all dates contained in the string.
func Parse(s string) []time.Time {
	dates := []time.Time{}

	elements := strings.Fields(s)

	for i, e := range elements {
		yearStart := strings.Index(e, "20")
		if yearStart < 0 {
			// not a year
			continue
		}
		if len(e)-yearStart < 4 {
			// no space for year
			continue
		}

		for j := 0; j < i && j < 3; j++ {
			e = strings.Join(elements[i-j:i+1], " ")
			d, err := parseDate(e)
			if err != nil {
				continue
			}
			dates = append(dates, d)
			break
		}

	}

	return dates
}

func parseDate(s string) (time.Time, error) {
	// try standard formats
	for _, f := range dateFormats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}

	// parse other formats
	elements := strings.Fields(s)
	if len(elements) != 3 {
		return time.Now(), errors.New("could not read date")
	}

	month := -1
	day := -1
	year := -1

	for _, e := range elements {
		// parse number if it exists
		i, err := strconv.Atoi(strings.TrimSuffix(e, "."))
		if err != nil {
			// if it's not a number, it might be a month name
			for k, v := range monthNames {
				if strings.Contains(strings.ToLower(e), k) {
					month = v
				}
			}
			continue
		}

		if i >= 2000 {
			// if it's 2000 or greater, it must be the year
			year = i
			continue
		}

		if strings.HasSuffix(e, ".") || i <= 31 {
			// if it has a dot, it must be a day
			day = i
			continue
		}
	}

	if day < 0 || month < 0 || year < 0 {
		return time.Now(), errors.New("could not read date")
	}
	return time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.Local), nil
}
