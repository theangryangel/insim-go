package session

import (
	"errors"
)

// ErrNotEnough ...
var (
	// ErrNotEnough is not enough data
	ErrNotEnough = errors.New("Not enough data")
	// ErrUnknownType is unknown type
	ErrUnknownType = errors.New("Unknown Packet Type")
	// ErrNoPacket is when no packet comes from an unmarshal call
	ErrNoPacket = errors.New("No packet returned from Unmarshal")
	// ErrTimeout is when a timeout occurs after 70s
	ErrTimeout = errors.New("Timeout after 70 seconds")
)
