package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	var ripchordFileName, midiInPort, midiOutPort string
	flag.StringVar(&ripchordFileName, "f", "", "Ripchord preset file (.rpc)")
	flag.StringVar(&midiInPort, "i", "", "MIDI input port")
	flag.StringVar(&midiOutPort, "o", "", "MIDI output port")
	flag.Parse()

	// ==================================
	// Get Ripchord struct from file name
	// ==================================
	if ripchordFileName == "" {
		fmt.Println("No Ripchord file specified")
		fmt.Println("Run with '-f' to specify a flag: -f \"PATH/TO_FILE\"")
		return
	}

	var ripchordXML RipchordXML
	err := parseRipchordXML(&ripchordXML, ripchordFileName)
	if err != nil {
		fmt.Println("Unable to parse Ripchord XML file")
		fmt.Println(err)
		return
	}

	ripchord, err := ripchordFromXML(ripchordXML)
	if err != nil {
		fmt.Println("Unable to convert marshalled XML to Ripchord struct")
		fmt.Println(err)
		return
	}

	// ============
	// Get MIDI I/O
	// ============
	drv, err := rtmididrv.New()

	var in drivers.In
	err = getMidiIn(&in, midiInPort, drv)
	if err != nil {
		fmt.Println(err)
		return
	}

	var out drivers.Out
	err = getMidiOut(&out, midiOutPort, drv)
	if err != nil {
		fmt.Println(err)
		return
	}

	send, _ := midi.SendTo(out)

	// ====================
	// Handle MIDI messages
	// ====================
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
