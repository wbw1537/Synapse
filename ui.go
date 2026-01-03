package synapse

import "embed"

// UI holds our static web server content.
//go:embed web/dist/*
var UI embed.FS
