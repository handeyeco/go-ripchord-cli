package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Ripchord struct {
	Map map[uint8]RipchordMapping
}

func ripchordFromXML(rcxml RipchordXML) (*Ripchord, error) {
	var rv Ripchord
	newMap := make(map[uint8]RipchordMapping)

	for _, input := range rcxml.Preset.Inputs {
		triggerNoteInt, err := strconv.Atoi(input.Note)
		if err != nil {
			return nil, err
		}
		triggerNote := uint8(triggerNoteInt)

		outputNotesStr := strings.Split(input.Chord.Notes, ";")
		var outputNotes []uint8
		for _, note := range outputNotesStr {
			parsedNoteInt, err := strconv.Atoi(note)
			if err != nil {
				return nil, err
			}
			parsedNote := uint8(parsedNoteInt)
			outputNotes = append(outputNotes, parsedNote)
		}

		var newMapping RipchordMapping
		newMapping.TriggerNote = triggerNote
		newMapping.OutputNotes = outputNotes
		newMapping.Name = input.Chord.Name
		newMap[triggerNote] = newMapping
	}

	rv.Map = newMap

	return &rv, nil
}

func parseRipchordXML(ripchordXML *RipchordXML, ripchordFileName string) error {
	xmlFile, err := os.Open(ripchordFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("file not found: %v\n", ripchordFileName))
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return errors.New("unable to read XML file")
	}

	err = xml.Unmarshal(byteValue, &ripchordXML)
	if err != nil {
		return errors.New("unable to marshall XML file")
	}

	return nil
}

type RipchordMapping struct {
	TriggerNote uint8
	OutputNotes []uint8
	Name        string
}

type RipchordXML struct {
	XMLName xml.Name  `xml:"ripchord"`
	Preset  PresetXML `xml:"preset"`
}

type PresetXML struct {
	XMLName xml.Name   `xml:"preset"`
	Inputs  []InputXML `xml:"input"`
}

type InputXML struct {
	XMLName xml.Name `xml:"input"`
	Note    string   `xml:"note,attr"`
	Chord   ChordXML `xml:"chord"`
}

type ChordXML struct {
	XMLName xml.Name `xml:"chord"`
	Name    string   `xml:"name,attr"`
	Notes   string   `xml:"notes,attr"`
}
