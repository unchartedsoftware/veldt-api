package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/unchartedsoftware/prism/util/color"
)

const (
	indent = "    "
)

func formatErr(err error) string {
	str := color.RemoveColor(err.Error())
	return strings.Replace(str, "\"", "\\\"", -1)
}

func handleErr(w http.ResponseWriter, err error) {
	// write error header
	w.WriteHeader(500)
	// error string
	str := fmt.Sprintf("{\n%s\"success\": \"false\",\n%s\"error\": \"%s\"\n}",
		indent,
		indent,
		formatErr(err))
	// write error
	fmt.Fprint(w, str)
}
