package main

import (
	"strconv"
	"strings"
)

type ErrHeaderParseFailed string

func (e ErrHeaderParseFailed) Error() string {
	return "Header parsing failed: " + string(e)
}

type roomMap struct {
	roomName        string
	idColumn        int
	titleColumn     int
	panelistsColumn int
	descColumn      int
}

const roomMaxColumns = 4

func (room *roomMap) isComplete() bool {
	switch {
	case room == nil:
	case room.roomName == "":
	case room.idColumn < 1:
	case room.titleColumn < 1:
	case room.panelistsColumn < 1:
	case room.descColumn < 1:
		return false
	default:
		return true
	}
	return true
}

func (room *roomMap) isSafelyReadable(columnCount int) bool {
	return room != nil &&
		room.idColumn < columnCount &&
		room.titleColumn < columnCount &&
		room.panelistsColumn < columnCount &&
		room.descColumn < columnCount
}

func parseScheduleHeader(firstRecord []string, secondRecord []string) ([]roomMap, error) {
	// Validation
	switch {
	case firstRecord == nil:
		return nil, ErrHeaderParseFailed("firstRecord cannot be nil")
	case secondRecord == nil:
		return nil, ErrHeaderParseFailed("secondRecord cannot be nil")
	case len(firstRecord) != len(secondRecord):
		return nil, ErrHeaderParseFailed("Unable to parse schedule header due to mismatched first/second record length.")
	case len(firstRecord) < 4:
		return nil, ErrHeaderParseFailed("Header is not sufficently long to contain at least the time and one room set of columns (4).")
	case !strings.Contains(firstRecord[0], "Time") && !strings.Contains(secondRecord[0], "Time"):
		return nil, ErrHeaderParseFailed("Unable to find 'Time' in first column.")
	}

	roomMappings := make([]roomMap, 0, 20)

	thisRoom := roomMap{}

	var roomStartAt int = 1

	for i := 1; i < len(firstRecord); i++ {
		switch {
		case firstRecord[i] != "":
			thisRoom.roomName = firstRecord[i]
		case strings.EqualFold(secondRecord[i], "ID"):
			thisRoom.idColumn = i
		case strings.EqualFold(secondRecord[i], "Title"):
			thisRoom.titleColumn = i
		case strings.EqualFold(secondRecord[i], "Panelists"):
			thisRoom.panelistsColumn = i
		case strings.EqualFold(secondRecord[i], "Public Description"):
			thisRoom.descColumn = i
		default:
			continue
		}
		if thisRoom.isComplete() {
			roomMappings = append(roomMappings, thisRoom)
			thisRoom = roomMap{}
			roomStartAt = i + 1
		} else if (i - roomStartAt) >= roomMaxColumns {
			return nil, ErrHeaderParseFailed("Unable to map room starting at column " + strconv.Itoa(roomStartAt))
		}
	}

	return roomMappings, nil
}
