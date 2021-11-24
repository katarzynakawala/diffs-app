package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/katarzynakawala/diffs-app/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//map in the role of a cache
	cache := map[string]*template.Template{}

	//slice of pages with the same extension
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	//iterating over pages
	for _, page := range pages {
		//extract file name from path
		name := filepath.Base(page)

		//parse the page template file to a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//add 'layout' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		//add 'partial' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		//Add template set to the cache, using the name of the page as key
		cache[name] = ts
	}
	return cache, nil
}