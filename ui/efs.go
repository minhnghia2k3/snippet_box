package ui

import "embed"

// Store file in `ui/html` and `ui/static` is an embedded filesystem.
//
//go:embed "html" "static"
var Files embed.FS
