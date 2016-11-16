package http

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/prism/util/color"
)

const (
	indent = "    "
)

func handleErr(w http.ResponseWriter, err error) {
	// write error header
	w.WriteHeader(500)
	// error string
	str := fmt.Sprintf("{\n%s\"success\": \"false\",\n%s\"error\": \"%s\"\n}",
		indent,
		indent,
		color.RemoveColor(err.Error()))
	// write error
	fmt.Fprint(w, str)
}
