package synapse

import "embed"

// Content holds our static web server content.
//go:embed web/dist/*
var Content embed.FS
