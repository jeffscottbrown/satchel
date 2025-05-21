package server

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/satchel/auth"
	"github.com/jeffscottbrown/satchel/repository"
	"github.com/markbates/goth/gothic"
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

	router.GET("/", rootHandler)
	router.GET("/employee/:employeeName", auth.AuthRequired, employeeHandler)
	router.GET("/unauthorized", unauthorizedHandler)
	auth.ConfigureAuthorizationHandlers(router)
}

func unauthorizedHandler(c *gin.Context) {
	renderUnauthorized(c, gin.H{})
}

func rootHandler(c *gin.Context) {
	employees, err := repository.GetEmployees()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving employees: %v", err)
		return
	}
	user, _ := gothic.GetFromSession("authenticatedUser", c.Request)
	renderTemplate(c, "main", gin.H{
		"Employees":         employees,
		"AuthenticatedUser": user,
	})
}

func employeeHandler(c *gin.Context) {
	employeeName := c.Param("employeeName")
	employee, err := repository.GetEmployeeByName(employeeName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving employee: %v", err)
		return
	}
	renderTemplate(c, "card", gin.H{
		"Employee": employee,
	})
}

var tmpl *template.Template
