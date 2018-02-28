package common

import ()

type Adapter interface {
	Find(collection string, id string, result *interface{}) error
	Create(collection string, doc interface{}) error
	Close()
}
