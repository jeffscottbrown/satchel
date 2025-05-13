package server

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed assets/**
var embeddedAssets embed.FS

//go:embed html/*.html
var embeddedHTMLFiles embed.FS

func Run() {
	router := createRouter()

	router.Run()
}

func createRouter() *gin.Engine {
	router := gin.Default()
	configureRoutes(router)
	return router
}

func configureRoutes(router *gin.Engine) {
	staticFiles, _ := fs.Sub(embeddedAssets, "assets")
	router.StaticFS("/static", http.FS(staticFiles))

	tmpl = template.Must(template.New("").Funcs(router.FuncMap).ParseFS(embeddedHTMLFiles, "html/*.html"))

	router.GET("/", renderRoot)
}

func renderRoot(c *gin.Context) {
	renderTemplate(c, "main", gin.H{})
}

var tmpl *template.Template
