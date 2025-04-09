package handlers

import (
	_ "embed"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	w.Write(tmpls["index.html"])
	return nil
}
