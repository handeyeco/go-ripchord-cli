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
	fmt.Println(ripchord.Map)

	in, err := midi.FindInPort("microKEY-37 KEYBOARD")
	if err != nil {
		fmt.Println("Unable to find microKEY-37 KEYBOARD")
		return
	}

	fmt.Println(in)

	var notesPressed []uint8
	// var notesActive []uint8
	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, key, vel uint8
		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			if !includes(notesPressed, key) {
				notesPressed = append(notesPressed, key)
			}
			fmt.Println(notesPressed)
		case msg.GetNoteEnd(&ch, &key):
			notesPressed = filter(notesPressed, key)
			fmt.Println(notesPressed)
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
