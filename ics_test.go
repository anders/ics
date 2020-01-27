package ics

import (
	"bytes"
	"log"
	"testing"
)

func TestEncode(t *testing.T) {
	cal := NewCalendar()
	if err := cal.Add(Event{"test": "abcd"}); err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}
	if err := cal.Encode(buf); err != nil {
		log.Fatal(err)
	}
}
