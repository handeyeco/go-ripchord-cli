package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func includes(s []uint8, e uint8) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func filter(s []uint8, e uint8) (ret []uint8) {
	for _, v := range s {
		if v != e {
			ret = append(ret, v)
		}
	}

	return
}

func calculateNotesActive(pressedNotes []uint8, ripchord Ripchord) (ret []uint8) {
	for _, pressedNote := range pressedNotes {
		mapping, ok := ripchord.Map[pressedNote]

		if ok {
			for _, mappedNote := range mapping.OutputNotes {
				if !includes(ret, mappedNote) {
					ret = append(ret, mappedNote)
				}
			}
		} else {
			if !includes(ret, pressedNote) {
				ret = append(ret, pressedNote)
			}
		}
	}

	return
}

func main() {
	defer midi.CloseDriver()

	// Open our xmlFile
	xmlFile, err := os.Open("ripchord/MP Neo Soul X-16.rpc")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened XML")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var ripchordXML RipchordXML
	xml.Unmarshal(byteValue, &ripchordXML)
	// fmt.Println(ripchordXML)

	ripchord, err := ripchordFromXML(ripchordXML)
	if err != nil {
		fmt.Println("Unable to convert marshalled XML to Ripchord struct")
		return
	}

	in, err := midi.FindInPort("microKEY-37 KEYBOARD")
	if err != nil {
		fmt.Println("Unable to find microKEY-37 KEYBOARD")
		return
	}

	fmt.Println(in)

	var notesPressed []uint8
	var notesActive []uint8
	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, key, vel uint8
		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			if !includes(notesPressed, key) {
				notesPressed = append(notesPressed, key)
				notesActive = calculateNotesActive(notesPressed, *ripchord)
			}
			// fmt.Println(notesPressed)
			fmt.Println(notesActive)
		case msg.GetNoteEnd(&ch, &key):
			notesPressed = filter(notesPressed, key)
			notesActive = calculateNotesActive(notesPressed, *ripchord)
			// fmt.Println(notesPressed)
			fmt.Println(notesActive)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	for {
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			stop()
			return
		}
	}
}
