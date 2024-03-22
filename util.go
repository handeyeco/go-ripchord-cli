package main

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
