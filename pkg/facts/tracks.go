package facts

type Track struct {
	Name string
	Code string
}

var unknownTrack = Track{
	Name: "Unknown",
	Code: "???",
}

var Tracks = []Track{
	{Name: "Blackwood GP", Code: "BL1"},
	{Name: "Blackwood Historic", Code: "BL2"},
}

func TrackFromCode(code string) *Track {
	for _, t := range Tracks {
		if t.Code == code {
			return &t
		}
	}

	return &unknownTrack
}
