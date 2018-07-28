package controllers

import (
	"net/http"
	"text/template"
)

func index(templates *template.Template) {
	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		home := new(context)
		home.Title = "Ask Me Anything"
		//w.Write([]byte("<h1> Hello Harold </h1>"))
		w.Header().Add("Content Type", "text/html")
		templates.Lookup("index.html").Execute(w, home)

	})
}
