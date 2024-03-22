package main

import (
	"errors"
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func getMidiIn(in *drivers.In, midiInPort string, drv *rtmididrv.Driver) error {
	inputs, err := drv.Ins()
	if err != nil || len(inputs) == 0 {
		return errors.New("unable to get MIDI inputs")
	}

	fmt.Println("Available input ports are:")
	for _, i := range inputs {
		fmt.Println(i)
	}

	if midiInPort == "" {
		*in = inputs[0]
		fmt.Printf("Defaulted to first MIDI in: %v\n", *in)
	} else {
		inPort, err := midi.FindInPort(midiInPort)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to find MIDI in: %v\n", midiInPort))
		}
		*in = inPort
		fmt.Printf("Found provided MIDI input: %v\n", *in)
	}

	return nil
}

func getMidiOut(out *drivers.Out, midiOutPort string, drv *rtmididrv.Driver) error {
	outputs, err := drv.Outs()
	if err != nil || len(outputs) == 0 {
		return errors.New("unable to get MIDI outputs")
	}

	fmt.Println("Available output ports are:")
	for _, i := range outputs {
		fmt.Println(i)
	}

	if midiOutPort == "" {
		*out = outputs[0]
		fmt.Printf("Defaulted to first MIDI out: %v\n", *out)
	} else {
		outPort, err := midi.FindOutPort(midiOutPort)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to find MIDI out: %v\n", midiOutPort))
		}
		*out = outPort
		fmt.Printf("Found provided MIDI output: %v\n", *out)
	}

	return nil
}
