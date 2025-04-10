package handlers

import (
	"log"

	"github.com/Kry0z1/fancytasks/static/templates"
)

var (
	tmpls = make(map[string][]byte)
	paths = []string{
		"index.html",
		"register.html",
		"login.html",
	}
)

func init() {
	fs := templates.GetTmplFS()

	var err error
	for _, path := range paths {
		tmpls[path], err = fs.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read template %s: %s", path, err.Error())
		}
	}
}
