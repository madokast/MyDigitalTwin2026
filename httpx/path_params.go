package httpx

import (
	"fmt"
	"net/http"
	"strconv"
)

type PathParams struct {
	*http.Request
}

func (p PathParams) Int64(name string) (int64, *Error) {
	valStr := p.PathValue(name)
	if valStr == "" {
		return 0, NewBadRequestError(fmt.Sprintf(
			"parameter %s in path %s is not set",
			name, p.Pattern,
		))
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, NewBadRequestError(fmt.Sprintf(
			"parameter %s in path %s expected an integer, but got %s",
			name, p.Pattern, valStr,
		))
	}

	return val, nil
}
