package main

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
)

type roomMap struct {
	name            string
	timeColumn      int
	nameColumn      int
	panelistsColumn int
	descColumn      int
}

func (room *roomMap) isComplete() bool {
	switch {
	case room == nil:
	case room.name == "":
	case room.nameColumn < 1:
	case room.panelistsColumn < 1:
	case room.descColumn < 1:
		return false
	default:
		return true
	}
	return true
}

func parseScheduleHeader(firstRecord []string, secondRecord []string) ([]*roomMap, error) {
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

	roomMappings := make([]*roomMap, 0, 20)

	var thisRoom *roomMap = new(roomMap)
	thisRoom.timeColumn = 0

	var roomStartAt int = 1

	for i := 1; i < len(firstRecord); i++ {
		switch {
		case firstRecord[i] != "":
			thisRoom.name = firstRecord[i]
		case strings.EqualFold(secondRecord[i], "Name"):
			thisRoom.nameColumn = i
		case strings.EqualFold(secondRecord[i], "Panelists"):
			thisRoom.panelistsColumn = i
		case strings.EqualFold(secondRecord[i], "Public Description"):
			thisRoom.descColumn = i
		default:
			continue
		}
		if thisRoom.isComplete() {
			roomMappings = append(roomMappings, thisRoom)
			thisRoom = nil
			roomStartAt = i + 1
		} else if (i - roomStartAt) >= 3 {
			return nil, ErrHeaderParseFailed("Unable to map room starting at column " + strconv.Itoa(roomStartAt))
		}
	}

	return roomMappings, nil
}

func parseScheduleRow(row []string, events *map[string]*Event, mapping []*roomMap) error {
	return nil
}

func parseSchedule(reader *csv.Reader) (*map[string]*Event, error) {
	rowOne, err := reader.Read()
	if err != nil {
		log.Println("Unable to reader header line 1 due to error: " + err.Error())
		return nil, err
	}
	rowTwo, err := reader.Read()
	if err != nil {
		log.Println("Unable to reader header line 2 due to error: " + err.Error())
		return nil, err
	}

	mapping, err := parseScheduleHeader(rowOne, rowTwo)
	if err != nil {
		return nil, err
	}

	var row []string = nil
	events := make(map[string]*Event)
	for {
		row, err = reader.Read()
		if row == nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Unable to complete parsing of schedule due to error: " + err.Error())
			return nil, err
		}
		err = parseScheduleRow(row, &events, mapping)
		if err != nil {
			log.Println("Unable to complete parsing of schedule due to error: " + err.Error())
			return nil, err
		}
	}

	return &events, nil
}
