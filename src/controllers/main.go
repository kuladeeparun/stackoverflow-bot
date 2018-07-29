package controllers

import (
	"text/template"
)

type context struct {
	Title string
}

type answer struct {
	Answerer string `json:"answerer"`
	Question string `json:"question"`
	QBody    string `json:"qbody"`
	Body     string `json:"answer"`
	AID      string `json:"aID"`
}

func Register(templates *template.Template) {
	index(templates)
	ask()
}
