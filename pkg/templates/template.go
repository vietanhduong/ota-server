package templates

import (
	"github.com/labstack/echo/v4"
	"io"
	"text/template"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
