package main

import (
	"log"
	"strings"
	"time"
)

func parseScheduleRow(lineNumber int, row []string, events EventSchedule, times *[]time.Time, mapping []roomMap, rowLength time.Duration, eventLocation *time.Location) error {
	switch {
	case row == nil:
		return ErrArgumentNil{"row", "parseScheduleRow"}
	case events == nil:
		return ErrArgumentNil{"events", "parseScheduleRow"}
	case mapping == nil:
		return ErrArgumentNil{"mapping", "parseScheduleRow"}
	case eventLocation == nil:
		return ErrArgumentNil{"eventLocation", "parseScheduleRow"}
	}

	// parse the time
	startTime, err := time.ParseInLocation("2006-01-02 3:04 PM", row[0], eventLocation)
	if err != nil {
		log.Printf("Skipping line %d due to date parse error: %s\n", lineNumber, err.Error())
		return nil
	}
	endTime := startTime.Add(rowLength)

	*times = append(*times, startTime)
	rowLen := len(row)
	for _, room := range mapping {
		if !room.isSafelyReadable(rowLen) {
			log.Printf("Unable to populate room at row %d due it being outside of set of columns.", lineNumber)
			continue
		}
		// pull the room's events, if possible
		roomEvents, ok := events[room.roomName]
		if !ok {
			roomEvents = make([]Event, 0, 1)
			events[room.roomName] = roomEvents
		}

		// Build event
		event := Event{
			StartTime:   startTime,
			EndTime:     endTime,
			ID:          row[room.idColumn],
			Title:       row[room.titleColumn],
			Panelists:   row[room.panelistsColumn],
			Description: row[room.descColumn],
			IsGuest:     strings.EqualFold(row[room.isGuestColumn], "yes"),
            Recording:   row[room.recordingColumn],
            CallMix:     row[room.callMixColumn],
		}

		// merge events that are longer than half-hour
		if rl := len(roomEvents); rl > 0 {
			prevEvent := &roomEvents[rl-1]
			if prevEvent.SharesID(&event) {
				prevEvent.EndTime = event.EndTime
				continue
			}
		}
		if event.IsValid() {
			if events.ContainsID(event.ID) {
				log.Printf("Skipped event on line %d in room %s due to duplicate UID %s\n", lineNumber, room.roomName, event.ID)
			} else {
				roomEvents = append(roomEvents, event)
				events[room.roomName] = roomEvents
			}
		}
	}

	return nil
}

func (this EventSchedule) ContainsID(id string) bool {
	if this == nil {
		return false
	}
	for _, events := range this {
		for _, event := range events {
			if event.HasID(id) {
				return true
			}
		}
	}
	return false
}
