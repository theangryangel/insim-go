package facts

import (
	"errors"
)

type Track struct {
	Name    string
	Code    string
	License string
}

var ErrUnknownTrack = errors.New("Unknown Track")

var Tracks = []Track{

	// TODO: Add reverse, etc.

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

	{Name: "Kyoto Oval", Code: "KY1", License: "S2"},
	{Name: "Kyoto National", Code: "KY2", License: "S2"},
	{Name: "Kyoto GP Long", Code: "KY3", License: "S2"},

	{Name: "Westhill National", Code: "WE1", License: "S2"},
	{Name: "Westhill International", Code: "WE2", License: "S2"},
	{Name: "Westhill Car Park", Code: "WE3", License: "S2"},
	{Name: "Westhill Karting", Code: "WE4", License: "S2"},
	{Name: "Westhill Karting National", Code: "WE5", License: "S2"},

	{Name: "Aston Cadet", Code: "AS1", License: "S2"},
	{Name: "Aston	Club", Code: "AS2", License: "S2"},
	{Name: "Aston National", Code: "AS3", License: "S2"},
	{Name: "Aston Historic", Code: "AS4", License: "S2"},
	{Name: "Aston Grand Prix", Code: "AS5", License: "S2"},
	{Name: "Aston Grand Touring", Code: "AS6", License: "S2"},
	{Name: "Aston North", Code: "AS7", License: "S2"},

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
