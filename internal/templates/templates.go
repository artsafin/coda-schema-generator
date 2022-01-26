package templates

import (
	"embed"
)

//go:embed *.tmpl types
var FS embed.FS
