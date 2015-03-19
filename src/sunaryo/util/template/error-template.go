package template

import (
	"fmt"
	"net/http"
	"strconv"
)

func IfError(err error, w http.ResponseWriter, status int, elseFn func()) {
	if err != nil {
		fmt.Fprintf(w, ("{\"status\":" + strconv.Itoa(status) + "}"))
		panic(err)
	} else {
		elseFn()
	}
}
