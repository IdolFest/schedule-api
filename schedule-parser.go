package main

import (
	"encoding/csv"
	"io"
	"log"
	"time"
)

func parseSchedule(reader *csv.Reader) (map[string][]Event, []time.Time, error) {
	eventLocation, _ := time.LoadLocation("America/New_York")

	rowOne, err := reader.Read()
	if err != nil {
		log.Println("Unable to reader header line 1 due to error: " + err.Error())
		return nil, nil, err
	}
	rowTwo, err := reader.Read()
	if err != nil {
		log.Println("Unable to reader header line 2 due to error: " + err.Error())
		return nil, nil, err
	}

	mapping, err := parseScheduleHeader(rowOne, rowTwo)
	if err != nil {
		return nil, nil, err
	}

	// We read two lines ahead of the loop
	lineNumber := 2

	var row []string = nil
	events := make(map[string][]Event)
	times := make([]time.Time, 0, 30)
	for {
		row, err = reader.Read()
		if row == nil && err == io.EOF {
			break
		}
		// increment the line number since we've now read that line
		lineNumber++

		if err != nil {
			log.Println("Unable to complete parsing of schedule due to error: " + err.Error())
			return nil, nil, err
		}
		err = parseScheduleRow(lineNumber, row, events, &times, mapping, eventLocation)
		if err != nil {
			log.Println("Unable to complete parsing of schedule due to error: " + err.Error())
			return nil, nil, err
		}

	}

	return events, times, nil
}
