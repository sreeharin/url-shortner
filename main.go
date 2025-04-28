package main

import (
	"fmt"

	"github.com/jxskiss/base62"
)

func convertURL(original string) string {
	LIMIT := 6
	encoded := base62.EncodeToString([]byte(original))
	if len(encoded) > LIMIT {
		return encoded[len(encoded)-LIMIT:]
	}

	return encoded
}

func main() {
	var url string
	fmt.Println("Enter the URL to shorten:")
	fmt.Scanln(&url)

	fmt.Println(convertURL(url))
}
