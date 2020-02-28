package main

import (
	"fmt"
	"strconv"
	"time"
)

type timeframe struct {
	a, b timecode
}

type timecode struct {
	Hour, Minute, Second, Millisecond int
}

// SetFromMatches updates the pointer with the passed matches []string
// these matches will correspond to both timecodes; hour, minute, second, millisecond (index 0 to 7)
func (t *timeframe) SetFromMatches(matches []string) *timeframe {
	var m []int
	for i, name := range matches {
		if i != 0 && name != "" {
			v, _ := strconv.Atoi(matches[i])
			m = append(m, v)
		}
	}
	t.a.Hour, t.a.Minute, t.a.Second, t.a.Millisecond = m[0], m[1], m[2], m[3]
	t.b.Hour, t.b.Minute, t.b.Second, t.b.Millisecond = m[4], m[5], m[6], m[7]

	return t
}

// String formats according to the SubRip-format for a timeframe and returns the resulting string
func (t *timeframe) String() string {
	return fmt.Sprintf(
		"%02d:%02d:%02d,%03d --> %02d:%02d:%02d,%03d",
		t.a.Hour, t.a.Minute, t.a.Second, t.a.Millisecond,
		t.b.Hour, t.b.Minute, t.b.Second, t.b.Millisecond,
	)
}

// Offset manipulates the timeframe based on the given offset in milliseconds
func (t *timeframe) Offset(offset string) *timeframe {
	o, _ := strconv.Atoi(offset)
	t.a.offset(o)
	t.b.offset(o)

	return t
}

func (t *timecode) offset(offset int) *timecode {
	millisecond := int(time.Duration(t.Millisecond) * time.Millisecond)
	newTime := time.Date(
		2010, 1, 1, // placeholder date
		t.Hour, t.Minute, t.Second, millisecond, time.UTC,
	).Add(time.Millisecond * time.Duration(offset))

	t.Hour = newTime.Hour()
	t.Minute = newTime.Minute()
	t.Second = newTime.Second()
	t.Millisecond = newTime.Nanosecond() / int(time.Millisecond)

	return t
}
