package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/bbalet/stopwords"
	"github.com/tidwall/gjson"
)

type context struct {
	Title string
}

type answer struct {
	Answerer string `json:"answerer"`

	Body string `json:"answer"`
}

func main() {

	templates := populateTemplates()

	home := new(context)
	home.Title = "Ask Me Anything"

	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		//w.Write([]byte("<h1> Hello Harold </h1>"))
		w.Header().Add("Content Type", "text/html")
		templates.Lookup("index.html").Execute(w, home)

	})

	http.HandleFunc("/ask", func(w http.ResponseWriter, req *http.Request) {
		question := req.FormValue("question")
		cleanContent := stopwords.CleanString(question, "en", true)
		println(cleanContent)

		//tagSearch := strings.Join(strings.Split(cleanContent, " "), ";")
		tagSearch := strings.Replace(cleanContent, " ", ";", -1)
		println(tagSearch)

		head := "https://api.stackexchange.com/2.2/tags/"
		tail := "/info?order=desc&sort=popular&site=stackoverflow"
		tresp, err := http.Get(head + tagSearch + tail)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(404)
		}

		tagsJson, err := ioutil.ReadAll(tresp.Body)
		//tags := make([]string, 10)
		var tags strings.Builder
		if err != nil {

			log.Fatal(err)
			w.WriteHeader(404)
		} else {
			items := gjson.Get(string(tagsJson), "items").Array()
			if len(items) > 0 {
				for _, item := range items {
					if item.Get("count").Int() > 20000 {
						tags.WriteString(item.Get("name").String())
						tags.WriteString(";")
					}
				}
				println(tags.String())

				qhead := "https://api.stackexchange.com/2.2/search/advanced?order=desc&sort=activity&q=%s"
				qtail := "&tagged=%s&accepted=True&site=stackoverflow"
				qurl := fmt.Sprint(fmt.Sprintf(qhead, strings.Replace(cleanContent, " ", "%20", -1)), fmt.Sprintf(qtail, tags.String()))
				println(qurl)

				qresp, err := http.Get(qurl)
				qJson, err := ioutil.ReadAll(qresp.Body)
				if err != nil {
					log.Fatal(err)
					w.WriteHeader(404)
				} else {
					items := gjson.Get(string(qJson), "items").Array()
					if len(items) > 0 {
						dat := new([]answer)
						w.Header().Add("Content Type", "application/json")
						ids := new([]string)
						for _, ans := range items {
							/* a := new(answer)
							a.Answerer = ans.Get("owner.display_name").String()
							a.Title = ans.Get("title").String()
							a.Link = ans.Get("link").String()

							*dat = append(*dat, *a) */
							*ids = append(*ids, ans.Get("accepted_answer_id").String())
						}
						ahead := "https://api.stackexchange.com/2.2/answers/%s"
						atail := "?order=desc&sort=activity&site=stackoverflow&filter=!9Z(-wzu0T"

						aurl := fmt.Sprint(fmt.Sprintf(ahead, strings.Join(*ids, ";")), atail)
						println(aurl)
						aresp, err := http.Get(aurl)
						if err != nil {
							log.Fatal(err)
							w.WriteHeader(404)
						} else {
							ajson, err := ioutil.ReadAll(aresp.Body)

							if err != nil {
								log.Fatal(err)
								w.WriteHeader(404)
							} else {
								items := gjson.Get(string(ajson), "items").Array()

								for _, ans := range items {
									a := new(answer)
									a.Answerer = ans.Get("owner.display_name").String()
									a.Body = ans.Get("body").String()

									*dat = append(*dat, *a)
								}
							}
						}
						json.NewEncoder(w).Encode(dat)
					} else {
						println("No answers found")
						w.WriteHeader(404)
					}
				}

			} else {
				println("No tags found for the question")
				w.WriteHeader(204)
			}
		}

		//w.Header().Add("Content Type", "plain/text")
		//w.Write([]byte(cleanContent))
	})

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("../public"))))

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
