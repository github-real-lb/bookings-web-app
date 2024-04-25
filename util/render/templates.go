package render

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

type Renderer interface {
	LoadTemplates(path, basePattern, tmplPattern string) error
	RenderTemplate(name string, data any) ([]byte, error)
}

type SmartRenderer struct {
	Templates map[string]*template.Template
}

// NewTemplateRenderer creates a new empty instance of TemplateRenderer.
// path is the templates directory path
func NewSmartRenderer() *SmartRenderer {
	return &SmartRenderer{
		Templates: make(map[string]*template.Template),
	}
}

// loadTemplatesFromDirectory loads specific templates from specific sub-directory in TemplatesPath.
// If any of the arguments are empty ("") function aboarts and returns nil.
func (tr *SmartRenderer) LoadTemplates(path, basePattern, tmplPattern string) error {
	if path == "" || basePattern == "" || tmplPattern == "" {
		return nil
	}

	path = strings.TrimSuffix(path, "/")
	basePattern = fmt.Sprintf("%s/%s", path, basePattern)
	tmplPattern = fmt.Sprintf("%s/%s", path, tmplPattern)

	// get the names of all the files matching pagePattern
	tmpls, err := filepath.Glob(tmplPattern)
	if err != nil {
		return err
	}

	// range thruogh all the pagePattern files
	for _, tmpl := range tmpls {
		// extracting the filename itself from the full path
		name := filepath.Base(tmpl)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(tmpl)
		if err != nil {
			return err
		}

		ts, err = ts.ParseFiles(basePattern)
		if err != nil {
			return err
		}

		tr.Templates[name] = ts
	}

	return nil
}

// RenderTemplate execute a template
func (tr *SmartRenderer) RenderTemplate(name string, data any) ([]byte, error) {
	if tr.Templates == nil {
		return nil, fmt.Errorf("templates cache is nil. please initiate first using NewSmartRenderer function")
	}

	var err error

	// checks if gohtml template exist in cache
	t, ok := tr.Templates[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find %s in template cache", name)
	}

	// check for error in template execution before passing it to w (http.ResponseWriter)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
