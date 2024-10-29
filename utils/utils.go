package gUtils

import (
	"strings"
)

func IsJSONFile(filename string) bool {
	return strings.HasSuffix(filename, ".json")
}
