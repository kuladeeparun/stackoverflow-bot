package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"nlp"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/tidwall/gjson"
)

var qId2TitleMap map[string]string
var qId2BodyMap map[string]string
var qTitle2IdMap map[string]string

func ask() {
	http.HandleFunc("/ask", func(w http.ResponseWriter, req *http.Request) {
		query := req.FormValue("query")
		cleanQuery := stopwords.CleanString(query, "en", true)
		println(cleanQuery)

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panic occured: ", r)
			}
		}()

		tags := getTags(cleanQuery, w)

		qId2TitleMap = make(map[string]string)
		qId2BodyMap = make(map[string]string)
		qTitle2IdMap = make(map[string]string)
		_ = getQuestions(cleanQuery, tags, w)

		relevantAnswerIDs, _ := nlp.MatchingQuestions(query, qTitle2IdMap)

		dat := getAnswers(relevantAnswerIDs, w)
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
		panic("No tags")
	}
	return tags.String()
}

func getQuestions(cleanQuery, tags string, w http.ResponseWriter) []string {
	qhead := "https://api.stackexchange.com/2.2/search/advanced?order=desc&sort=activity&q=%s"
	qtail := "&tagged=%s&accepted=True&filter=!9Z(-wwYGT&site=stackoverflow"
	qurl := fmt.Sprint(fmt.Sprintf(qhead, strings.Replace(cleanQuery, " ", "%20", -1)), fmt.Sprintf(qtail, tags))

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
			id := ans.Get("accepted_answer_id").String()
			title := ans.Get("title").String()
			body := ans.Get("body").String()
			*ids = append(*ids, id)
			qId2TitleMap[id] = title
			qId2BodyMap[id] = body
			qTitle2IdMap[title] = id

		}
	} else {
		println("No answers found for the question")
		w.WriteHeader(204)
		panic("No answers")
	}
	return *ids
}

func getAnswers(ids []string, w http.ResponseWriter) []answer {
	ahead := "https://api.stackexchange.com/2.2/answers/%s"
	atail := "?order=desc&sort=activity&site=stackoverflow&filter=!9Z(-wzu0T"

	aurl := fmt.Sprint(fmt.Sprintf(ahead, strings.Join(ids, ";")), atail)

	dat := new([]answer)

	ajson := fetch(aurl)

	items := gjson.Get(string(ajson), "items").Array()

	//fmt.Println(qId2TitleMap)

	for _, ansID := range ids {
		for _, ans := range items {

			qid := ans.Get("answer_id").String()
			if ansID == qid {
				a := new(answer)
				a.Answerer = ans.Get("owner.display_name").String()
				a.Body = ans.Get("body").String()

				//fmt.Println(qid)
				a.Question = qId2TitleMap[qid]
				a.QBody = qId2BodyMap[qid]
				a.AID = qid
				*dat = append(*dat, *a)
			}
		}
	}
	return *dat

}

func fetch(url string) string {
	//fmt.Println(url)
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
