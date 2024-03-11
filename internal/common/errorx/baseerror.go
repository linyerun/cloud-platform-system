package errorx

var defaultCode = 10010

type BaseError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BaseError) Error() string {
	return e.Msg
}

func NewBaseError(code int, msg string) error {
	return &BaseError{Code: code, Msg: msg}
}
func NewDefaultError(msg string) error {
	return NewBaseError(defaultCode, msg)
}
