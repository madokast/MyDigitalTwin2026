package tags

import (
	"dt2026/httpx"
	"strings"
)

func NormalizeTag(tag string) (string, *httpx.Error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", httpx.NewBadRequestError("tag cannot be empty or contain only whitespace")
	}
	return tag, nil
}
