package common

import ()

type Adapter interface {
	Find(collection string, id string) (*Document, error)
	Create(collection string, document *Document) error
}

type Document struct {
	Data map[string]interface{}
}
