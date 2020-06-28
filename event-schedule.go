package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	// "os"
)

var webClient = http.Client{}

var eventCacheTime time.Time
var eventsCache EventSchedule
var timesCache []time.Time
var roomOrderCache []string

func readEventCache(sheetUrl string, cacheTimeLength int, timezone string) (EventSchedule, []time.Time, []string, error) {
	if time.Now().After(eventCacheTime) {
		err := updateEventCache(sheetUrl, cacheTimeLength, timezone)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return eventsCache, timesCache, roomOrderCache, nil
}

func updateEventCache(sheetUrl string, cacheTimeLength int, timezone string) error {

	rawEvents, err := fetchRawEvents(sheetUrl)
	if err != nil {
		log.Println("Unable to update event cache due to an error during fetchRawEvents: " + err.Error())
		return err
	}

	reader := csv.NewReader(bytes.NewReader(rawEvents))

	newEvents, newTimes, newRoomOrder, err := parseSchedule(reader, timezone)
	if err != nil {
		log.Println("Unable to update event cache due to an error during parseSchedule: " + err.Error())
		return err
	} else {
		log.Printf("I got %d room\n", len(newEvents))
		eventsCache = newEvents
		timesCache = newTimes
		roomOrderCache = newRoomOrder
		eventCacheTime = time.Now().Add(time.Duration(cacheTimeLength) * time.Second)
	}

	return nil
}

func fetchRawEvents(sheetUrl string) ([]byte, error) {
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
