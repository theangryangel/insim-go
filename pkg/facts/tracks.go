package facts

import "errors"

type Track struct {
	Name string
	Code string
}

var ErrUnknownTrack = errors.New("Unknown Track")

var Tracks = []Track{
	{Name: "Blackwood GP", Code: "BL1"},
	{Name: "Blackwood Historic", Code: "BL2"},
	{Name: "Blackwood RallyX", Code: "BL3"},
	{Name: "Blackwood Car Park", Code: "BL4"},

	{Name: "South City Classic", Code: "SO1"},
	{Name: "South City Sprint 1", Code: "SO2"},
	{Name: "South City Sprint 2", Code: "SO3"},
	{Name: "South City Long", Code: "SO4"},
	{Name: "South City Town", Code: "SO5"},
	{Name: "South City Chicane", Code: "SO6"},

	{Name: "Fern Bay Club", Code: "FE1"},
	{Name: "Fern Bay Green", Code: "FE2"},
	{Name: "Fern Bay Gold", Code: "FE3"},
	{Name: "Fern Bay Black", Code: "FE4"},
	{Name: "Fern Bay RallyX ", Code: "FE5"},
	{Name: "Fern Bay RallyX Green", Code: "FE6"},

	{Name: "AutoX", Code: "AU1"},
	{Name: "Skid Pad", Code: "AU2"},
	{Name: "Drag Strip", Code: "AU3"},
	{Name: "Drag Strip (8 Lane)", Code: "AU4"},

	{Name: "Rockingham ISSC", Code: "RO1"},
	{Name: "Rockingham National", Code: "RO2"},
	{Name: "Rockingham Oval", Code: "RO3"},
	{Name: "Rockingham ISSC Long", Code: "RO4"},
	{Name: "Rockingham Lake", Code: "RO5"},
	{Name: "Rockingham Handling", Code: "RO6"},
	{Name: "Rockingham International", Code: "RO7"},
	{Name: "Rockingham Historic", Code: "RO8"},
	{Name: "Rockingham Historic Short", Code: "RO9"},
	{Name: "Rockingham International Long", Code: "RO10"},
	{Name: "Rockingham Sportscar", Code: "RO11"},
}

func TrackFromCode(code string) (*Track, error) {
	for _, t := range Tracks {
		if t.Code == code {
			return &t, nil
		}
	}

	return nil, ErrUnknownTrack
}
