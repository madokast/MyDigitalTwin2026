package httpx

import (
	"fmt"
	"net/http"
	"strconv"
)

type PathParams struct {
	*http.Request
}

func (p PathParams) String(name string) (string, *Error) {
	val := p.PathValue(name)
	if val == "" {
		return "", NewBadRequestError(fmt.Sprintf(
			"parameter %s in path %s is not set",
			name, p.Pattern,
		))
	}

	return val, nil
}

func (p PathParams) Int64(name string) (int64, *Error) {
	valStr, httpErr := p.String(name)
	if httpErr != nil {
		return 0, httpErr
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, NewBadRequestError(fmt.Sprintf(
			"parameter %s in path %s expected an integer, but got: '%s'",
			name, p.Pattern, valStr,
		))
	}

	return val, nil
}
