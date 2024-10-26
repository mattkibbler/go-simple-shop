package output

import (
	"html/template"
	"io"
)

type PageData struct {
	Title string
	Data  any
}

func RenderPage(templates *template.Template, w io.Writer, template string, pageData PageData) error {
	// Write header to output
	err := templates.ExecuteTemplate(w, "header.html", pageData)
	if err != nil {
		return err
	}
	// Write content template to output
	err = templates.ExecuteTemplate(w, template, pageData)
	if err != nil {
		return err
	}
	// Write footer template to output
	err = templates.ExecuteTemplate(w, "footer.html", pageData)
	if err != nil {
		return err
	}
	return nil
}
