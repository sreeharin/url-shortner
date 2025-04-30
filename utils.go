package main

import "github.com/jxskiss/base62"

func convertURL(original string) URL {
	var converted URL
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
