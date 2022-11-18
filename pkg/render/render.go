package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Gideon-isa/bookings/pkg/config"
	"github.com/Gideon-isa/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTempaltes sets the config for the template
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// tc: A declare template.Template with nil as value
	// does not contain those templates
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser")
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache is an initialize *template.Template{}. a pointer to template.Template.
	// its contents are accessible to any instance of template.Template{}
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	// myCache is map with a pointer *template.Templates as its value. all its contents are accessible
	// to any initilize template.Templates{}.
	// Which is sent/acessible to config.AppConfig.TemplateCache field which is a *template.Template
	return myCache, nil
}
