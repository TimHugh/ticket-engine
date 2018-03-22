package rollbar

import (
	"github.com/stvp/roll"
)

type Reporter struct {
	Client roll.Client
}

func New(token, env string) Reporter {
	return Reporter{
		roll.New(token, env),
	}
}

func (r Reporter) Error(err error) error {
	var emptyCustom map[string]string
	_, reportErr := r.Client.Error(err, emptyCustom)
	return reportErr
}
