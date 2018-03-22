package rollbar

import (
	"fmt"

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

func (r Reporter) Error(err error) {
	var emptyCustom map[string]string
	_, reportErr := r.Client.Error(err, emptyCustom)
	if reportErr != nil {
		fmt.Printf("ERROR: Unable to report error to rollbar!\nOriginal Error: %s\nRollbar error: %s", err, reportErr)
	}
}
