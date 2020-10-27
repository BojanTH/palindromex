package service

import (
	"net/http"
	"html/template"
	"strings"
	"os"
	"path/filepath"
)

const viewsDir = "web/template/views/"
const layoutsDir = "web/template/layouts/"

// Template is a wraper for templates
type Template struct {
	Templ *template.Template
	Flash *Flash
}

// Execute wraps the execution of the template.Execute
// currently this is used to simplify rendering of flash messages
func (t *Template) Execute(response http.ResponseWriter, request *http.Request, data interface{}) error {
	if nil == data {
		data = make(map[string]interface{})
	}
	if d, ok := data.(Flasher); ok {
		d.SetFlashes(t.Flash.GetFlashes(response, request))
	}

	return t.Templ.Execute(response, data)
}

// GetTemplates returns parsed templates from viewsDir
func GetTemplates(flash *Flash) (map[string]*Template, error) {
	templates := make(map[string]*Template)

	var layoutPaths []string
	err := filepath.Walk(layoutsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			layoutPaths = append(layoutPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(viewsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// templatePath will be [subDirInViewsDir/]templateName.html
			templatePath := strings.Replace(path, viewsDir, "", 1)
			baseName := filepath.Base(templatePath)
			tmpl := template.New(baseName).Funcs(getFunctionsMap())
			tmpl, err := tmpl.ParseFiles(append(layoutPaths, path)...)
		    if err != nil {
		    	return err
			}
			
			templates[templatePath] = &Template{Templ: tmpl, Flash: flash}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return templates, nil
}
