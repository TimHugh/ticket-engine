package mock

type ErrorReporter struct {
	Store error
}

func (r *ErrorReporter) Error(err error) {
	r.Store = err
}
