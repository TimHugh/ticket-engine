package common

import ()

type Adapter interface {
	Find(collection string, id string) (*struct{}, error)
	Create(collection string, document *struct{}) error
	Close()
}
