package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/satchel/auth"
)

func renderTemplateWithStatus(c *gin.Context, templateName string, data gin.H, status int) {
	c.Status(status)
	x := c.GetHeader("HX-Request")

	isHTMX := x != ""

	data["IsAuthenticated"] = auth.IsAuthenticated(c.Request)

	if isHTMX {
		tmpl.ExecuteTemplate(c.Writer, templateName, data)
	} else {
		data["Body"] = template.HTML(renderTemplateToString(templateName, data))
		tmpl.ExecuteTemplate(c.Writer, "layout", data)
	}
}

func renderForbidden(c *gin.Context, data gin.H) {
	renderTemplateWithStatus(c, "forbidden", data, http.StatusForbidden)
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
