package lib

import "strings"

// SQLEscapeLike 将用户输出变成纯粹的 contains 查询
func EscapeSQLLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
