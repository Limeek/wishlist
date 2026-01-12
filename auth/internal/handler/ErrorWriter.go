package handler

import (
	"fmt"
	"net/http"
)

func writeError(code int, err error, w http.ResponseWriter) {
	w.WriteHeader(code)

	_, writeErr := w.Write([]byte(err.Error()))

	if writeErr != nil {
		fmt.Printf("Error writing response: %s", writeErr)
	}
}
