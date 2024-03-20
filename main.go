package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

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

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, key, vel uint8
		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
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
