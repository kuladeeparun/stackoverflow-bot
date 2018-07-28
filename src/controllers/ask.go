package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/tidwall/gjson"
)

func ask() {
	http.HandleFunc("/ask", func(w http.ResponseWriter, req *http.Request) {
		query := req.FormValue("query")
		cleanQuery := stopwords.CleanString(query, "en", true)
		println(cleanQuery)

		tags := getTags(cleanQuery, w)

		answerIDs := getQuestions(cleanQuery, tags, w)

		dat := getAnswers(answerIDs, w)
		json.NewEncoder(w).Encode(dat)
	})
}

func getTags(cleanQuery string, w http.ResponseWriter) string {
	tagSearch := strings.Replace(cleanQuery, " ", ";", -1)
	println(tagSearch)

	head := "https://api.stackexchange.com/2.2/tags/"
	tail := "/info?order=desc&sort=popular&site=stackoverflow"

	tagsJson := fetch(head + tagSearch + tail)
	var tags strings.Builder
	items := gjson.Get(tagsJson, "items").Array()
	if len(items) > 0 {
		for _, item := range items {
			if item.Get("count").Int() > 20000 {
				tags.WriteString(item.Get("name").String())
				tags.WriteString(";")
			}
		}
		println(tags.String())

	} else {
		println("No tags found for the question")
		w.WriteHeader(204)
	}
	return tags.String()
}

func getQuestions(cleanQuery, tags string, w http.ResponseWriter) []string {
	qhead := "https://api.stackexchange.com/2.2/search/advanced?order=desc&sort=activity&q=%s"
	qtail := "&tagged=%s&accepted=True&site=stackoverflow"
	qurl := fmt.Sprint(fmt.Sprintf(qhead, strings.Replace(cleanQuery, " ", "%20", -1)), fmt.Sprintf(qtail, tags))
	println(qurl)

	qJson := fetch(qurl)
	items := gjson.Get(string(qJson), "items").Array()
	ids := new([]string)
	if len(items) > 0 {

		w.Header().Add("Content Type", "application/json")

		for _, ans := range items {
			/* a := new(answer)
			a.Answerer = ans.Get("owner.display_name").String()
			a.Title = ans.Get("title").String()
			a.Link = ans.Get("link").String()

			*dat = append(*dat, *a) */
			*ids = append(*ids, ans.Get("accepted_answer_id").String())
		}
	} else {
		println("No answers found for the question")
		w.WriteHeader(204)
	}
	return *ids
}

func getAnswers(ids []string, w http.ResponseWriter) []answer {
	ahead := "https://api.stackexchange.com/2.2/answers/%s"
	atail := "?order=desc&sort=activity&site=stackoverflow&filter=!9Z(-wzu0T"

	aurl := fmt.Sprint(fmt.Sprintf(ahead, strings.Join(ids, ";")), atail)
	println(aurl)

	dat := new([]answer)

	ajson := fetch(aurl)

	items := gjson.Get(string(ajson), "items").Array()

	for _, ans := range items {
		a := new(answer)
		a.Answerer = ans.Get("owner.display_name").String()
		a.Body = ans.Get("body").String()

		*dat = append(*dat, *a)
	}

	return *dat

}

func fetch(url string) string {
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
		//w.WriteHeader(404)
		return ""
	}

	stringResponse, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		//w.WriteHeader(404)
		return ""
	}

	return string(stringResponse)
}
