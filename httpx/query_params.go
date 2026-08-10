package httpx

import (
	"dt2026/lib"
	"fmt"
	"net/url"
	"strconv"

	"dt2026/lib/optional"
)

type QueryParams url.Values

func (q QueryParams) GetOptionalSingleString(key string) (optional.Optional[string], *Error) {
	rawValues := q[key]

	if len(rawValues) == 0 {
		return optional.None[string](), nil
	} else if len(rawValues) == 1 {
		return optional.Some(rawValues[0]), nil
	} else {
		return optional.Optional[string]{}, NewBadRequestError(fmt.Sprintf(
			"query parameter %s expected a single value, but got %d values: (%s)",
			key, len(rawValues), lib.SliceToString(rawValues),
		))
	}
}

func (q QueryParams) GetOptionalSingleInt64(key string) (optional.Optional[int64], *Error) {
	opt, httpErr := q.GetOptionalSingleString(key)
	if httpErr != nil {
		return optional.Optional[int64]{}, httpErr
	}

	s, ok := opt.Get()
	if !ok {
		return optional.None[int64](), nil
	}

	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return optional.Optional[int64]{}, NewBadRequestError(fmt.Sprintf(
			"query parameter %s expected an integer, but got: '%s'",
			key, s,
		))
	}
	return optional.Some(value), nil
}

func (q QueryParams) GetOptionalStrings(key string) []string {
	return q[key]
}
