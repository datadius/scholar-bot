package apihandlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type StudyEmbed struct {
	Title    string
	Url      string
	Authors  string
	Abstract string
}

func QueryFirstGs(query string) *StudyEmbed {
	// Define the URL of the Google Scholar search page
	url := fmt.Sprintf("https://scholar.google.com/scholar?hl=en&q=%s", query)

	// Define a User-Agent header
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}

	// Send a GET request to the URL with the User-Agent header
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Check if the request was successful (status code 200)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		// Parse the HTML content of the page using goquery
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Find all the search result blocks with class "gs_ri"
		s := doc.Find(".gs_ri").First()
		// Extract the title and URL
		titleElem := s.Find("h3.gs_rt")
		title := titleElem.Text()
		url, _ := titleElem.Find("a").Attr("href")

		// Extract the authors and publication details
		authorsElem := s.Find("div.gs_a")
		authors := authorsElem.Text()

		// Extract the abstract or description
		abstractElem := s.Find("div.gs_rs")
		abstract := abstractElem.Text()

		return &StudyEmbed{Title: title, Url: url, Authors: authors, Abstract: abstract}
	} else {
		fmt.Println("Failed to retrieve the page. Status code:", resp.StatusCode)
	}

	return nil
}
