package dbuser

import (
	"embed"
	"io/fs"
)

//go:embed scripts/*.sql
var migrationsAssets embed.FS

var MigrationAssets, _ = fs.Sub(migrationsAssets, "scripts")
