package main

import (
	"log"
	"time"
)


func parseScheduleRow(lineNumber int, row []string, events EventSchedule, mapping []roomMap, eventLocation *time.Location) error {
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
	startTime, err := time.ParseInLocation("2006-01-02 15:04", row[0], eventLocation)
	if err != nil {
		log.Printf("Skipping line %d due to date parse error: %s\n", lineNumber, err.Error())
		return nil
	}
	endTime := startTime.Add(time.Duration(30) * time.Minute)

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
		}

		// merge events that are longer than half-hour
		if rl := len(roomEvents); rl > 0 {
			prevEvent := roomEvents[rl-1]
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
