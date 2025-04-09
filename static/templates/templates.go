package templates

import "embed"

//go:embed *.html
var embedTmplFS embed.FS

func GetTmplFS() embed.FS {
	return embedTmplFS
}
