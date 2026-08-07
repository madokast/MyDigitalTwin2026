package lib

import (
	"fmt"
	"strings"
)

func SliceToString[E any, S ~[]E](s S) string {
	var ss []string
	for _, e := range s {
		ss = append(ss, fmt.Sprintf("'%v'", e))
	}
	return strings.Join(ss, ", ")
}
