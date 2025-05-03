package utils

import (
	"strings"

	"github.com/jxskiss/base62"
	"github.com/sreeharin/url-shortner/internal/models"
)

// convertURL takes a string and converts it to base62 format.
func ConvertURL(original string) models.URL {
	var converted models.URL
	if !strings.HasPrefix(original, "http") {
		original = "http://" + original
	}
	converted.Original = original

	LIMIT := 6
	encoded := base62.EncodeToString([]byte(original))

	if len(encoded) > LIMIT {
		converted.Shortened = encoded[len(encoded)-LIMIT:]
	} else {
		converted.Shortened = encoded
	}

	return converted
}
