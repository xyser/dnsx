package asset

import (
	"embed"
)

//go:embed ui
var UI embed.FS

//go:embed sql
var SQL embed.FS
