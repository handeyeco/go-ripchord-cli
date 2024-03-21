> [!CAUTION]
> This is not meant to be an active or consumable module. Use at your own risk!
>
> Sharing because I'm caring! Not because I'm going to maintain it.

# Go Ripchord thing

**tl:dr** - this is a toy CLI project that listens to incoming MIDI notes and sends multiple MIDI notes out based on [Ripchord](https://trackbout.com/ripchord) mappings

## About

During a personal development day, I decided to try to get more comfortable with Golang - the backend language my company uses. It's a pretty straight-forward language, but some of the nuances of Golang tooling were lost on me because they were abstracted away by our Infra team.

This is a toy CLI project that:

- Accepts a [Ripchord](https://trackbout.com/ripchord) XML file (.rpc)
  - At their core, Ripchord XML files are a set of mappings between 1 MIDI note and _n_ MIDI notes
  - I used the standard `encoding/xml` module to convert the XML to a struct
- Using the [gomidi midi module](https://gitlab.com/gomidi/midi), finds MIDI I/O
- Listens for incoming MIDI, maps single incoming MIDI notes to the _n_ notes specified in the Ripchord file, and sends those _n_ notes back out via MIDI
  - It tries to intelligently merge and split notes, since playing two mappings at once may result in duplicate notes

This means you can play chords using single notes. For instance if you had mappings of root notes to major triads:

- You play C4 (MIDI note 60)
- This program would send out C4 (60), E4 (64), and G4 (67)

However there are no limitations to the mappings. You could have C4 just play C#4 if you wanted to.

## Use

```
$ git clone https://github.com/handeyeco/go-ripchord-cli.git
$ cd go-ripchord-cli
$ go run . -f "presets/test.rpc" -i "NAME OF MIDI IN DEVICE" -o "NAME OF MIDI OUT DEVICE"
```
