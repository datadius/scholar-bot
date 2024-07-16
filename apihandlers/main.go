package apihandlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type StudyStruct struct {
	Title    string
	Url      string
	Authors  string
	Abstract string
}

func QueryFirstGs(query string) (*StudyStruct, bool) {
	// Define the URL of the Google Scholar search page
	urlQuery := fmt.Sprintf(
		"https://scholar.google.com/scholar?hl=en&q=%s",
		url.QueryEscape(query),
	)
	fmt.Println(urlQuery)

	// Define a User-Agent header
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}

	// Send a GET request to the URL with the User-Agent header
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		log.Printf("error getting the request %v", err)
		return nil, false
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Check if the request was successful (status code 200)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error in executing the request %v", err)
		return nil, false
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	if resp.StatusCode == 200 {
		// Parse the HTML content of the page using goquery
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("error in parsing html %v", err)
			return nil, false
		}

		// Find all the search result blocks with class "gs_ri"
		s := doc.Find(".gs_ri").First()
		// Extract the title and URL
		titleElem := s.Find("h3.gs_rt")
		title := titleElem.Text()
		urlResp, _ := titleElem.Find("a").Attr("href")

		// Extract the authors and publication details
		authorsElem := s.Find("div.gs_a")
		authors := authorsElem.Text()

		// Extract the abstract or description
		abstractElem := s.Find("div.gs_rs")
		abstract := abstractElem.Text()

		return &StudyStruct{Title: title, Url: urlResp, Authors: authors, Abstract: abstract}, true
	} else {
		fmt.Println("Failed to retrieve the page. Status code:", resp.StatusCode)
		return nil, false
	}
}

func QueryTopTenGs(query string) (*[]StudyStruct, bool) {
	// Define the URL of the Google Scholar search page
	urlQuery := fmt.Sprintf(
		"https://scholar.google.com/scholar?hl=en&q=%s",
		url.QueryEscape(query),
	)
	fmt.Println(urlQuery)

	// Define a User-Agent header
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}

	// Send a GET request to the URL with the User-Agent header
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		log.Printf("error getting the request %v", err)
		return nil, false
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Check if the request was successful (status code 200)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error in executing the request %v", err)
		return nil, false
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	if resp.StatusCode == 200 {
		// Parse the HTML content of the page using goquery
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("error in parsing html %v", err)
			return nil, false
		}

		var studySlice []StudyStruct

		// Find all the search result blocks with class "gs_ri"
		doc.Find(".gs_ri").Each(func(i int, s *goquery.Selection) {
			// Extract the title and URL
			titleElem := s.Find("h3.gs_rt")
			title := titleElem.Text()
			urlResp, _ := titleElem.Find("a").Attr("href")

			studySlice = append(
				studySlice,
				StudyStruct{Title: title, Url: urlResp, Authors: "", Abstract: ""},
			)

		})
		if len(studySlice) != 0 {
			return &studySlice, true
		} else {
			return nil, false
		}
	} else {
		fmt.Println("Failed to retrieve the page. Status code:", resp.StatusCode)
		return nil, false
	}
}
