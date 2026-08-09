package tags

import (
	"dt2026/httpx"
	"fmt"
	"strings"
)

func NormalizeTag(tag string) (string, *httpx.Error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", httpx.NewBadRequestError("tag cannot be empty or contain only whitespace")
	}
	return tag, nil
}

func NormalizeTags(tags []string) ([]string, *httpx.Error) {
	var normalizedTags []string
	var seen = make(map[string]bool)
	for _, tag := range tags {
		normalizedTag, err := NormalizeTag(tag)
		if err != nil {
			return nil, err
		}
		if seen[normalizedTag] {
			return nil, httpx.NewBadRequestError(fmt.Sprintf(
				"duplicate tags are not allowed: '%s'",
				normalizedTag,
			))
		}
		seen[normalizedTag] = true
		normalizedTags = append(normalizedTags, normalizedTag)
	}
	return normalizedTags, nil
}
