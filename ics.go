// Copyright (c) 2020 Anders Bergh.
// License: MIT.

// Package ics provides a simple iCalendar encoder.
package ics

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/anders/utils"
)

// Event is a single calendar entry. Supported types are strings.Stringer,
// string and time.Time.
type Event map[string]interface{}

// Calendar holds a list of Events.
type Calendar struct {
	Properties map[string]interface{} // VCALENDAR fields
	Events     []Event                // a list of Events
}

// NewCalendar returns a new instance of Calendar with some properties preset.
func NewCalendar() *Calendar {
	return &Calendar{
		Properties: map[string]interface{}{
			"VERSION": "2.0",
			"PRODID":  "-//github.com/anders/ics",
			"CALSCAL": "GREGORIAN",
		},
	}
}

// Get returns a property.
func (c *Calendar) Get(key string) interface{} {
	return c.Properties[key]
}

// Set sets a VCALENDAR property.
func (c *Calendar) Set(key string, value interface{}) {
	c.Properties[key] = value
}

// Add adds an object to the calendar.
func (c *Calendar) Add(obj interface{}) error {
	switch v := obj.(type) {
	case Event:
		c.Events = append(c.Events, v)
		return nil
	default:
		return fmt.Errorf("Add(): unsupported type %T", obj)
	}
}

// Encode writes a complete calendar to the specified writer.
func (cal Calendar) Encode(w io.Writer) error {
	if _, err := io.WriteString(w, "BEGIN:VCALENDAR\r\n"); err != nil {
		return err
	}

	keys := []string{}
	for key := range cal.Properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := cal.Properties[key]
		if err := encodeKV(w, key, value); err != nil {
			return err
		}
	}

	// TODO: handle other types of iCalendar entries.

	for _, ev := range cal.Events {
		if err := ev.Encode(w); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return nil
}

func encodeKV(w io.Writer, key string, value interface{}) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case time.Time:
		str = v.In(time.UTC).Format("20060102T150405Z")
	case fmt.Stringer:
		str = v.String()
	default:
		log.Printf("encodeKV(): key %s: unsupported type %T", key, v)
	}

	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, ";", "\\;", -1)
	str = strings.Replace(str, ",", "\\,", -1)
	str = strings.Replace(str, "\n", "\\n", -1)

	key = strings.ToUpper(key)

	lines := utils.SplitLength(key+":"+str, 72)
	for i, line := range lines {
		if i > 0 {
			line = " " + line
		}
		_, err := io.WriteString(w, line+"\r\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Encode writes a single event to the specified Writer.
// Supported values are strings and time.Time, as well as any time that conforms
// to the Stringer interface.
func (ev Event) Encode(w io.Writer) error {
	if _, err := io.WriteString(w, "BEGIN:VEVENT\r\n"); err != nil {
		return err
	}

	var keys []string
	for key := range ev {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if err := encodeKV(w, key, ev[key]); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "END:VEVENT\r\n"); err != nil {
		return err
	}

	return nil
}
