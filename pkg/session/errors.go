package session

import (
	"errors"
)

var (
	ErrNotEnough   = errors.New("Not enough data")
	ErrUnknownType = errors.New("Unknown Packet Type")
	ErrNoPacket    = errors.New("No packet returned from Unmarshal")
	ErrTimeout     = errors.New("Timeout after 70 seconds")
)
