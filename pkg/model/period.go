package model

import (
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/spout/pkg/utils"
)

type (
	RangeType string
)

type period struct {
	baseTime  string
	duration  string
	rangeType RangeType
	begin     time.Time
	end       time.Time
}

type Period interface {
	Begin() time.Time
	End() time.Time
}

func (x *period) Begin() time.Time { return x.begin }
func (x *period) End() time.Time   { return x.end }

func NewPeriod(baseTime, duration string, rangeType RangeType) (Period, error) {
	p := &period{
		baseTime:  baseTime,
		duration:  duration,
		rangeType: rangeType,
	}

	var ts time.Time
	if baseTime != "" {
		out, err := parseTime(baseTime)
		if err != nil {
			return nil, err
		}
		ts = *out
	} else {
		ts = time.Now().UTC()
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, goerr.Wrap(err, "invalid duration")
	}

	switch rangeType {
	case "before":
		p.begin = ts.Add(-d)
		p.end = ts
	case "after":
		p.begin = ts
		p.end = ts.Add(d)
	case "around":
		p.begin = ts.Add(-d / 2)
		p.end = ts.Add(d / 2)
	}

	utils.Logger.With("begin", p.begin).With("end", p.end).Debug("duration configured")
	return p, nil
}

func parseTime(s string) (*time.Time, error) {
	dateTimeLayouts := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
	}
	timeOnlyFormats := []string{
		"15:04:05",
		"15:04",
	}
	for _, layout := range dateTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return &t, nil
		}
	}

	for _, layout := range timeOnlyFormats {
		if t, err := time.Parse(layout, s); err == nil {
			now := time.Now()
			today := time.Date(
				now.Year(), now.Month(), now.Day(),
				t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

			if today.After(now) {
				return ptr.To(today.AddDate(0, 0, -1)), nil
			}

			return &today, nil
		}
	}
	return nil, goerr.New("No suitable time format").With("input", s)
}
