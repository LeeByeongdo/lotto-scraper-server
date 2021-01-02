package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"strconv"
)

func main() {
	addr := ":7171"
	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type pageInfo struct {
	StatusCode int
	Result     [5][6]int
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing url argument")
		return
	}

	c := colly.NewCollector()

	p := &pageInfo{}

	c.OnHTML("div.tbl_basic table tbody", func(e *colly.HTMLElement) {
		log.Println("onHTML")
		var result [5][6]int
		e.ForEach("tr", func(index int, tbody *colly.HTMLElement) {
			var arr [6]int
			tbody.ForEach("td span", func(idx int, span *colly.HTMLElement) {
				i, err := strconv.Atoi(span.Text)
				if err != nil {
					fmt.Println("parsing err")
					fmt.Println(err)
				} else {
					arr[idx] = i
				}
			})
			result[index] = arr
		})
		p.Result = result
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", URL)
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(URL)

	// dump results
	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	log.Println(p)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
