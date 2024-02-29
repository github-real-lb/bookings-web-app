package render

import (
	"log"
	"net/http"
	"text/template"
)

// renderTemplate parse a template from a file and execute the template.
func RenderTemplate(w http.ResponseWriter, gohtml string) {
	parsedTemplate, err := template.ParseFiles("./templates/" + gohtml)
	if err != nil {
		log.Printf("error parsing template %s: %v\n", gohtml, err)
		return
	}

	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Printf("error parsing template %s: %v\n", gohtml, err)
	}
}
