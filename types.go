package main

import (
	"time"
)

type ErrHeaderParseFailed string

func (e ErrHeaderParseFailed) Error() string {
	return "Header parsing failed: " + string(e)
}

type Event struct {
	time      time.Time
	room      string
	name      string
	desc      string
	panelists []string
}
