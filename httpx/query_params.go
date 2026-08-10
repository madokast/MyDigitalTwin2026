package httpx

import (
	"dt2026/lib"
	"fmt"
	"net/url"
	"strconv"

	"dt2026/lib/optional"
)

type QueryParams url.Values

func (q QueryParams) GetOptionalSingleInt64(key string) (optional.Optional[int64], *Error) {
	rawValues := q[key]

	if len(rawValues) == 0 {
		return optional.Optional[int64]{}, nil
	} else if len(rawValues) == 1 {
		value, err := strconv.ParseInt(rawValues[0], 10, 64)
		if err != nil {
			return optional.Optional[int64]{}, NewBadRequestError(fmt.Sprintf(
				"query parameter %s expected an integer, but got: '%s'",
				key, rawValues[0],
			))
		}
		return optional.Some(value), nil
	} else {
		return optional.Optional[int64]{}, NewBadRequestError(fmt.Sprintf(
			"query parameter %s expected a single value, but got %d values: (%s)",
			key, len(rawValues), lib.SliceToString(rawValues),
		))
	}
}

func (q QueryParams) GetOptionalStrings(key string) []string {
	return q[key]
}
