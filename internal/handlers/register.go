package handlers

import "net/http"

func RegisterPage(w http.ResponseWriter, r *http.Request) error {
	w.Write(tmpls["register.html"])
	return nil
}

func Register(w http.ResponseWriter, r *http.Request) error {
	return nil
}
