package common

import ()

type Location struct {
	ID           string
	SignatureKey string
}

type LocationRepository interface {
	Store(Location)
	Find(string) *Location
}
