package main

type LimitReachedError struct {
	msg string
}

func (self LimitReachedError) Error() string {
	return self.msg
}

func NewLimitReachedError() LimitReachedError {
	return LimitReachedError{"Limit reached."}
}
