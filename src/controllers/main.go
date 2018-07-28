package controllers

import (
	"text/template"
)

type context struct {
	Title string
}

type answer struct {
	Answerer string `json:"answerer"`
	/* Question string `json:"question"` */
	Body string `json:"answer"`
}

func Register(templates *template.Template) {
	index(templates)
	ask()
}
