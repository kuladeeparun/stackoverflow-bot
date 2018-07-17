package main

import (
	"net/http"
	"os"
	"text/template"
)

type context struct {
	Title string
}

func main() {
	templates := populateTemplates()

	home := new(context)
	home.Title = "Ask Me Anything"

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("<h1> Hello Harold </h1>"))
		w.Header().Add("Content Type", "text/html")
		templates.Lookup("index.html").Execute(w, home)

	})

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("../public"))))
	/* http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := "../public" + r.URL.Path
		fmt.Println(path)
		f, err := os.Open(path)
		if err == nil {
			br := bufio.NewReader(f)
			if strings.HasSuffix(path, ".css") {
				w.Header().Add("Content Type", "text/css")
			} else if strings.HasSuffix(path, ".gif") {
				w.Header().Add("Content Type", "image/gif")
			} else if strings.HasSuffix(path, ".png") {
				w.Header().Add("Content Type", "image/png")
			}

			br.WriteTo(w)
		} else {
			fmt.Println(err)
			w.WriteHeader(404)
		}
	}) */

	http.ListenAndServe(":8000", nil)
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "../templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)

	return result
}
