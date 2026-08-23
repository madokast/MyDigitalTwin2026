package lib

import "maps"

func MapsMerge[K comparable, V any, M ~map[K]V](ms ...M) map[K]V {
	var r = make(map[K]V)
	for _, m := range ms {
		maps.Copy(r, m)
	}
	return r
}
