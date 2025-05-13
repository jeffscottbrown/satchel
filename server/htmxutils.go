package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderTemplateWithStatus(c *gin.Context, templateName string, data gin.H, status int) {
	c.Status(status)
	x := c.GetHeader("HX-Request")

	isHTMX := x != ""

	if isHTMX {
		tmpl.ExecuteTemplate(c.Writer, templateName, data)
	} else {
		tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
			"Body": template.HTML(renderTemplateToString(templateName, data)),
		})
	}
}

// func renderBadRequest(c *gin.Context, data gin.H) {
// 	renderTemplateWithStatus(c, "error", data, http.StatusBadRequest)
// }

// func renderNotFound(c *gin.Context, data gin.H) {
// 	renderTemplateWithStatus(c, "error", data, http.StatusNotFound)
// }

func renderTemplate(c *gin.Context, templateName string, data gin.H) {
	renderTemplateWithStatus(c, templateName, data, http.StatusOK)
}

func renderTemplateToString(name string, data any) string {
	var buf []byte
	writer := &buffer{&buf}
	_ = tmpl.ExecuteTemplate(writer, name, data)
	return string(*writer.buf)
}

type buffer struct {
	buf *[]byte
}

func (w *buffer) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}
