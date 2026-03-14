package framework

import "fmt"

type Code int

const (
	CodeSuccess Code = iota
	CodeUnschedulable
	CodeWait
	CodeError
)

type Status struct {
	code    Code
	reasons []string
	err     error
}

func NewStatus(code Code, reasons ...string) *Status {
	return &Status{
		code:    code,
		reasons: reasons,
	}
}

func AsError(err error) *Status {
	if err == nil {
		return nil
	}
	return &Status{
		code: CodeError,
		err:  err,
	}
}

func (s *Status) Code() Code {
	if s == nil {
		return CodeSuccess
	}
	return s.code
}

func (s *Status) Reasons() []string {
	if s == nil {
		return nil
	}
	return s.reasons
}

func (s *Status) Err() error {
	if s == nil {
		return nil
	}
	if s.err != nil {
		return s.err
	}
	if len(s.reasons) == 0 {
		return nil
	}
	return fmt.Errorf("%v", s.reasons)
}

func (s *Status) IsSuccess() bool {
	return s == nil || s.code == CodeSuccess
}
