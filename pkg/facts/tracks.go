package facts

import "errors"

type Track struct {
	Name    string
	Code    string
	License string
}

var ErrUnknownTrack = errors.New("Unknown Track")

var Tracks = []Track{
	{Name: "Blackwood GP", Code: "BL1", License: "S1"},
	{Name: "Blackwood Historic", Code: "BL2", License: "S1"},
	{Name: "Blackwood RallyX", Code: "BL3", License: "S1"},
	{Name: "Blackwood Car Park", Code: "BL4", License: "S1"},

	{Name: "South City Classic", Code: "SO1", License: "S1"},
	{Name: "South City Sprint 1", Code: "SO2", License: "S1"},
	{Name: "South City Sprint 2", Code: "SO3", License: "S1"},
	{Name: "South City Long", Code: "SO4", License: "S1"},
	{Name: "South City Town", Code: "SO5", License: "S1"},
	{Name: "South City Chicane", Code: "SO6", License: "S1"},

	{Name: "Fern Bay Club", Code: "FE1", License: "S1"},
	{Name: "Fern Bay Green", Code: "FE2", License: "S1"},
	{Name: "Fern Bay Gold", Code: "FE3", License: "S1"},
	{Name: "Fern Bay Black", Code: "FE4", License: "S1"},
	{Name: "Fern Bay RallyX ", Code: "FE5", License: "S1"},
	{Name: "Fern Bay RallyX Green", Code: "FE6", License: "S1"},

	{Name: "AutoX", Code: "AU1", License: "S1"},
	{Name: "Skid Pad", Code: "AU2", License: "S1"},
	{Name: "Drag Strip", Code: "AU3", License: "S1"},
	{Name: "Drag Strip (8 Lane)", Code: "AU4", License: "S1"},

	// Missing all the S2 tracks

	{Name: "Rockingham ISSC", Code: "RO1", License: "S3"},
	{Name: "Rockingham National", Code: "RO2", License: "S3"},
	{Name: "Rockingham Oval", Code: "RO3", License: "S3"},
	{Name: "Rockingham ISSC Long", Code: "RO4", License: "S3"},
	{Name: "Rockingham Lake", Code: "RO5", License: "S3"},
	{Name: "Rockingham Handling", Code: "RO6", License: "S3"},
	{Name: "Rockingham International", Code: "RO7", License: "S3"},
	{Name: "Rockingham Historic", Code: "RO8", License: "S3"},
	{Name: "Rockingham Historic Short", Code: "RO9", License: "S3"},
	{Name: "Rockingham International Long", Code: "RO10", License: "S3"},
	{Name: "Rockingham Sportscar", Code: "RO11", License: "S3"},
}

func TrackFromCode(code string) (*Track, error) {
	for _, t := range Tracks {
		if t.Code == code {
			return &t, nil
		}
	}

	return nil, ErrUnknownTrack
}
