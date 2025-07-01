package server

import (
	"embed"
	"io/fs"

	"github.com/jeffscottbrown/satchel/repository"
	"github.com/markbates/goth/gothic"

	"html/template"
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

func init() {
	tmpl = template.Must(template.New("").ParseFS(embeddedHTMLFiles, "html/*.html"))
}

func configureRoutes(router *gin.Engine) {
	staticFiles, _ := fs.Sub(embeddedAssets, "assets")
	router.StaticFS("/static", http.FS(staticFiles))

	router.GET("/", rootHandler)
	router.GET("/employee/:employeeName", employeeHandler)
}

func forbiddenHandler(c *gin.Context) {
	renderForbidden(c, gin.H{})
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
