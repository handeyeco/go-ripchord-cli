package main

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Ripchord struct {
	Map map[int]RipchordMapping
}

func ripchordFromXML(rcxml RipchordXML) (*Ripchord, error) {
	var rv Ripchord
	newMap := make(map[int]RipchordMapping)

	for _, input := range rcxml.Preset.Inputs {
		triggerNote, err := strconv.Atoi(input.Note)
		if err != nil {
			return nil, err
		}

		outputNotesStr := strings.Split(input.Chord.Notes, ";")
		var outputNotes []int
		for _, note := range outputNotesStr {
			parsedNote, err := strconv.Atoi(note)
			if err != nil {
				return nil, err
			}
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

type RipchordMapping struct {
	TriggerNote int
	OutputNotes []int
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
