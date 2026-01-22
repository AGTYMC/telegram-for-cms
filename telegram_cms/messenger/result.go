package messenger

type Result struct {
	Success bool
	Err     error
	Data    any
	Message string
	Code    int32
}
