package utils

import (
	"github.com/jxskiss/base62"
)

// convertURL takes uint and converts it to base62 format.
func ConvertID(id uint) string {
	bytes := base62.FormatUint(uint64(id))

	return base62.EncodeToString(bytes)
}
