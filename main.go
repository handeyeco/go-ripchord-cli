package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/v2/drivers"
	"io"
	"os"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
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

	// List available MIDI I/O
	drv, err := rtmididrv.New()

	var ripchordFile, midiInPort, midiOutPort string
	flag.StringVar(&ripchordFile, "f", "", "Ripchord preset file (.rpc)")
	flag.StringVar(&midiInPort, "i", "", "MIDI input port")
	flag.StringVar(&midiOutPort, "o", "", "MIDI output port")
	flag.Parse()

	if ripchordFile == "" {
		fmt.Println("No Ripchord file specified")
		fmt.Println("Run with '-f' to specify a flag: -f \"PATH/TO_FILE\"")
		return
	}

	var in drivers.In
	inputs, err := drv.Ins()
	if midiInPort == "" {
		if err != nil {
			fmt.Println("Unable to get MIDI inputs")
			fmt.Println("Available ports are:")
			for _, i := range inputs {
				fmt.Println(i)
			}
			return
		}
		in = inputs[0]
		fmt.Printf("Defaulted to first MIDI in: %v\n", in)
	} else {
		in, err = midi.FindInPort(midiInPort)
		if err != nil {
			fmt.Printf("Unable to find MIDI in: %v\n", midiInPort)
			return
		}
		fmt.Printf("Found provided MIDI input: %v\n", in)
	}

	var out drivers.Out
	outputs, err := drv.Outs()
	if midiOutPort == "" {
		if err != nil {
			fmt.Println("Unable to get MIDI outputs")
			return
		}
		out = outputs[0]
		fmt.Printf("Defaulted to first MIDI out: %v\n", out)
	} else {
		out, err = midi.FindOutPort(midiOutPort)
		if err != nil {
			fmt.Printf("Unable to find MIDI out: %v\n", midiOutPort)
			fmt.Println("Available ports are:")
			for _, i := range outputs {
				fmt.Println(i)
			}
			return
		}
		fmt.Printf("Found provided MIDI output: %v\n", out)
	}

	// Open our xmlFile
	xmlFile, err := os.Open(ripchordFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Printf("File not found: %v\n", ripchordFile)
		return
	}

	fmt.Println("Successfully opened XML")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)
	var ripchordXML RipchordXML
	xml.Unmarshal(byteValue, &ripchordXML)

	ripchord, err := ripchordFromXML(ripchordXML)
	if err != nil {
		fmt.Println("Unable to convert marshalled XML to Ripchord struct")
		return
	}

	send, _ := midi.SendTo(out)

	var notesPressed []uint8
	var notesActive []uint8
	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, key, vel uint8
		switch {
		case msg.GetNoteStart(&ch, &key, &vel):
			if !includes(notesPressed, key) {
				notesPressed = append(notesPressed, key)
				newNotesActive := calculateNotesActive(notesPressed, *ripchord)

				for _, newNote := range newNotesActive {
					if !includes(notesActive, newNote) {
						fmt.Printf("MIDI On: %v\n", newNote)
						send(midi.NoteOn(ch, newNote, 100))
					}
				}

				notesActive = newNotesActive
			}
			fmt.Println(notesActive)
		case msg.GetNoteEnd(&ch, &key):
			if includes(notesPressed, key) {
				notesPressed = filter(notesPressed, key)
				newNotesActive := calculateNotesActive(notesPressed, *ripchord)

				for _, oldNote := range notesActive {
					if !includes(newNotesActive, oldNote) {
						fmt.Printf("MIDI Off: %v\n", oldNote)
						send(midi.NoteOff(ch, oldNote))
					}
				}

				notesActive = newNotesActive
			}
			fmt.Println(notesActive)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	fmt.Println("Ready...")

	for {
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			stop()
			return
		}
	}
}
