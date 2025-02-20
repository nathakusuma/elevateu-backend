package mailtmpl

import "embed"

//go:embed *.html
var Templates embed.FS
