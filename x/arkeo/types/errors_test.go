package types

import (
	"errors"
	fmt "fmt"

	. "gopkg.in/check.v1"
)

type errIsChecker struct {
	*CheckerInfo
}

var ErrIs Checker = &errIsChecker{
	&CheckerInfo{Name: "ErrIs", Params: []string{"obtained", "expected"}},
}

func (errIsChecker) Check(params []interface{}, names []string) (result bool, err string) {
	p1, ok1 := params[0].(error)
	p2, ok2 := params[1].(error)
	if !ok1 || !ok2 {
		result = false
		err = "must pass error types"
		return
	}
	result = errors.Is(p1, p2)
	if !result {
		err = fmt.Sprintf("Errors do not match!\nObtained: %s\nExpected: %s", p1, p2)
	}
	return
}
