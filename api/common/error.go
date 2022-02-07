package common

import "fmt"

type ErrReqParam struct {
	Name string
	Msg  string
}

func (e *ErrReqParam) Error() string {
	return fmt.Sprintf("param name '%s' with error '%s'", e.Name, e.Msg)
}
