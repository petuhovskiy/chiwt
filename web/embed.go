package web

import "embed"

const resourcesPrefix = "resources/"

//go:embed resources
var resources embed.FS

const templatesPrefix = "templates/"

//go:embed templates
var templates embed.FS
