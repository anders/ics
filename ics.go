// Copyright (c) 2020 Anders Bergh.
// License: MIT.

// Package ics provides a simple iCalendar encoder.
package ics

import (
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/anders/utils"
)

// Event is a single calendar entry. The values must be either string or
// time.Time.
type Event map[string]interface{}

// Calendar is a list of events.
type Calendar []Event

// Encode writes a complete calendar to the specified writer.
func (cal Calendar) Encode(w io.Writer) error {
	if _, err := io.WriteString(w, "BEGIN:VCALENDAR\r\n"+
		"VERSION:2.0\r\n"+
		"PRODID:-//github.com/anders/ics\r\n"+
		"CALSCAL:GREGORIAN\r\n"); err != nil {
		return err
	}

	for _, ev := range cal {
		if err := ev.Encode(w); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return nil
}

// Encode writes a single event to the specified Writer.
// Supported values are strings and time.Time.
func (ev Event) Encode(w io.Writer) error {
	if _, err := io.WriteString(w, "BEGIN:VEVENT\r\n"); err != nil {
		return err
	}

	var keys []string
	for key := range ev {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		var s string
		switch v := ev[k].(type) {
		case string:
			s = v
		case time.Time:
			s = v.In(time.UTC).Format("20060102T150405Z")
		default:
			log.Printf("key %s: unsupported type %T in Event.Encode()", k, v)
			continue
		}

		s = strings.Replace(s, "\\", "\\\\", -1)
		s = strings.Replace(s, ";", "\\;", -1)
		s = strings.Replace(s, ",", "\\,", -1)
		s = strings.Replace(s, "\n", "\\n", -1)

		k = strings.ToUpper(k)

		lines := utils.SplitLength(k+":"+s, 72)
		for i, line := range lines {
			if i > 0 {
				line = " " + line
			}
			_, err := io.WriteString(w, line+"\r\n")
			if err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(w, "END:VEVENT\r\n"); err != nil {
		return err
	}

	return nil
}
