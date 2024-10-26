package output

import (
	"html/template"
	"net/http"
)

func RenderPage(templates *template.Template, w http.ResponseWriter, template string, pageData any) {
	// Write header to output
	err := templates.ExecuteTemplate(w, "header.html", pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write content template to output
	err = templates.ExecuteTemplate(w, template, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write footer template to output
	err = templates.ExecuteTemplate(w, "footer.html", pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
