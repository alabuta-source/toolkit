package apiError

type optionFunc func(*requestError)

type Option interface {
	Apply(*requestError)
}

func (o optionFunc) Apply(reqErr *requestError) { o(reqErr) }

func NewApiError(message string, status int, options ...Option) RequestError {
	err := requestError{
		ErrorMessage: message,
		ErrorStatus:  status,
	}

	for _, option := range options {
		option.Apply(&err)
	}
	return err
}

func WithIncommingError(err error) optionFunc {
	return func(re *requestError) {
		re.ErrorCode = err.Error()
	}
}
