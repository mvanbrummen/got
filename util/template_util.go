package util

import (
	"io/ioutil"
	"text/template"
)

type Templates map[string]*template.Template

func InitTemplates(dir string) (Templates, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	templates := make(Templates, 0)

	for _, f := range files {
		templates[f.Name()] = template.Must(template.ParseFiles(dir + f.Name()))
	}
	return templates, nil
}
