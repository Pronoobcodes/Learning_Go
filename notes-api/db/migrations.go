package db

import "embed"

// Migrations contains the SQL migrations applied by the API at startup.
//
//go:embed migrations/*.sql
var Migrations embed.FS
