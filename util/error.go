package util

import (
	"github.com/unchartedsoftware/prism/util/color"
)

// FormatErr removes any coloring from the error string.
func FormatErr(err error) string {
	return color.RemoveColor(err.Error())
}
