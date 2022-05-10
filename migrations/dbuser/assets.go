package dbuser

import (
	"embed"
	"io/fs"
)

//go:embed *
var MmigrationsAssets embed.FS

var MigrationAssets, _ = fs.Sub(MmigrationsAssets, "scripts")
