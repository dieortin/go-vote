package main

import (
	"fmt"
	"html/template"
	"io"
)

const TemplateDir = "views"

type PageMap struct {
	pages        map[string]*template.Template
	baseTemplate string
}

func TemplatePath(filename string) string {
	return fmt.Sprintf("%s/%s", TemplateDir, filename)
}

func NewPageMap(base string, pages []string) PageMap {
	pageMap := PageMap{
		pages:        make(map[string]*template.Template),
		baseTemplate: base,
	}
	for _, name := range pages {
		pageMap.pages[name] = template.Must(template.ParseFiles(TemplatePath(name), TemplatePath(base)))
	}
	return pageMap
}

func (pm PageMap) ExecuteTemplate(wr io.Writer, tmpl string, data any) error {
	err := pm.pages[tmpl].ExecuteTemplate(wr, "base", data)
	if err != nil {
		return err
	}
	return nil
}
