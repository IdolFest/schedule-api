package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	// "sync"
	// "os"
)

var webClient = new(http.Client)

// var eventLock = new(sync.RWMutex)
// var eventCacheTime = new(time.Time)
var eventsCache EventSchedule

func readEventCache() (EventSchedule, error) {
	if len(eventsCache) == 0 {
		err := updateEventCache()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func updateEventCache() error {
	rawEvents, err := fetchRawEvents()
	if err != nil {
		log.Println("Unable to update event cache due to an error during fetchRawEvents: " + err.Error())
		return err
	}

	reader := csv.NewReader(bytes.NewReader(rawEvents))

	newEvents, err := parseSchedule(reader)
	if err != nil {
		log.Println("Unable to update event cache due to an error during parseSchedule: " + err.Error())
		return err
	} else {
		log.Printf("I got %d events\n", len(newEvents))
		
	}

	return nil
}

func fetchRawEvents() ([]byte, error) {
	// sheetUrl := os.Getenv("SHEET_URL")
	sheetUrl := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRqFQ_57wfOI6gNj8ORyAGjrJQ5JwDguv6_0lWtdraKoPYr2_VgHtdHqF010vKt6K2JgOn92Q7zcjpX/pub?gid=0&single=true&output=tsv"
	log.Printf("Going to call out to the following sheets URL: %s\n", sheetUrl)

	resp, err := webClient.Get(sheetUrl)
	if err != nil {
		log.Printf("Failed to fetch Events Schedule sheet with error: %s\n", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Events Schedule sheet with error: %s\n", err)
		return nil, err
	}

	return body, nil
}
