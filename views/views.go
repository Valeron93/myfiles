package views

import (
	"embed"
	"net/http"

	"github.com/gofiber/template/html/v2"
)

//go:embed *.html
var viewsFS embed.FS
var Engine = html.NewFileSystem(http.FS(viewsFS), ".html")
