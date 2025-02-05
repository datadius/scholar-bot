package apihandlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StudyStruct struct {
	Title    string
	Url      string
	Authors  string
	Abstract string
}

func QueryFirstGs(query string, minYear string) (*StudyStruct, bool) {
	// Define the URL of the Google Scholar search page
	urlQuery := fmt.Sprintf(
		"https://scholar.google.com/scholar?hl=en&q=%s&as_ylo=%s",
		url.QueryEscape(query), minYear,
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

func QueryTopTenGs(query string, minYear string) (*[]StudyStruct, bool) {
	// Define the URL of the Google Scholar search page
	urlQuery := fmt.Sprintf(
		"https://scholar.google.com/scholar?hl=en&q=%s&as_ylo=%s",
		url.QueryEscape(query), minYear,
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

func QueryFirstPMC(query string, minYear string) (*StudyStruct, bool) {
	//https://www.ncbi.nlm.nih.gov/books/NBK25499/#_chapter4_ESearch_
	urlQuery := fmt.Sprintf(
		"https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=%s&retmode=json&sort=relevance&retmax=1&mindate=%s&maxdate=2024",
		url.QueryEscape(query),
		minYear,
	)
	fmt.Println(urlQuery)

	// Define a User-Agent header
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Content-Type": "application/json",
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
		var idStudyList IdStudyList
		err := json.NewDecoder(resp.Body).Decode(&idStudyList)
		if err != nil {
			log.Printf("error translating json response from PMC API %v", err)
			return nil, false
		}
		if len(idStudyList.Esearchresult.Idlist) != 0 {
			idStudy := idStudyList.Esearchresult.Idlist[0]
			urlStudy := fmt.Sprintf(
				"https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=%s",
				idStudy,
			)
			log.Println(urlStudy)
			req, err := http.NewRequest("GET", urlStudy, nil)
			if err != nil {
				log.Printf("failed to form new request to get study details %v", err)
				return nil, false
			}
			headers := map[string]string{
				"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
				"Content-Type": "application/xml",
			}
			for key, value := range headers {
				req.Header.Add(key, value)
			}

			// Check if the request was successful (status code 200)
			studyResp, err := client.Do(req)
			if err != nil {
				log.Printf("error in executing the request %v", err)
				return nil, false
			}
			defer studyResp.Body.Close()
			log.Println(studyResp.StatusCode)
			if studyResp.StatusCode == 200 {
				var pubmedArticleSet PubmedArticleSet
				err := xml.NewDecoder(studyResp.Body).Decode(&pubmedArticleSet)
				if err != nil {
					log.Printf("Error decoding url  data for pmc study %v", err)
					return nil, false
				}
				if len(pubmedArticleSet.PubmedArticle) != 0 {
					title := pubmedArticleSet.PubmedArticle[0].MedlineCitation.Article.ArticleTitle
					urlArticle := fmt.Sprintf(
						"https://pubmed.ncbi.nlm.nih.gov/%s/",
						idStudyList.Esearchresult.Idlist[0],
					)
					//handle authors
					abstract := pubmedArticleSet.PubmedArticle[0].MedlineCitation.Article.Abstract.AbstractText
					return &StudyStruct{
						Title:    title,
						Url:      urlArticle,
						Authors:  "",
						Abstract: abstract,
					}, true
				}
			}
		}
	}

	return nil, false
}

func QueryTopTenPMC(query string, minYear string) (*[]StudyStruct, bool) {
	//https://www.ncbi.nlm.nih.gov/books/NBK25499/#_chapter4_ESearch_
	urlQuery := fmt.Sprintf(
		"https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=%s&retmode=json&sort=relevance&retmax=10&mindate=%s&maxdate=2024",
		url.QueryEscape(query),
		minYear,
	)
	fmt.Println(urlQuery)

	// Define a User-Agent header
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Content-Type": "application/json",
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
		var idStudyList IdStudyList
		err := json.NewDecoder(resp.Body).Decode(&idStudyList)
		if err != nil {
			log.Printf("error translating json response from PMC API %v", err)
			return nil, false
		}
		if len(idStudyList.Esearchresult.Idlist) != 0 {
			idStudies := strings.Join(idStudyList.Esearchresult.Idlist, ",")
			urlStudy := fmt.Sprintf(
				"https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=%s",
				idStudies,
			)
			log.Println(urlStudy)
			req, err := http.NewRequest("GET", urlStudy, nil)
			if err != nil {
				log.Printf("failed to form new request to get study details %v", err)
				return nil, false
			}
			headers := map[string]string{
				"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
				"Content-Type": "application/xml",
			}
			for key, value := range headers {
				req.Header.Add(key, value)
			}

			// Check if the request was successful (status code 200)
			studyResp, err := client.Do(req)
			if err != nil {
				log.Printf("error in executing the request %v", err)
				return nil, false
			}
			defer studyResp.Body.Close()
			log.Println(studyResp.StatusCode)
			if studyResp.StatusCode == 200 {
				var pubmedArticleSet PubmedArticleSet
				err := xml.NewDecoder(studyResp.Body).Decode(&pubmedArticleSet)
				if err != nil {
					log.Printf("Error decoding url  data for pmc study %v", err)
					return nil, false
				}
				var studyStructSlice []StudyStruct
				for _, value := range pubmedArticleSet.PubmedArticle {
					title := value.MedlineCitation.Article.ArticleTitle
					urlArticle := fmt.Sprintf(
						"https://pubmed.ncbi.nlm.nih.gov/%s/",
						value.MedlineCitation.PMID.Text,
					)
					studyStructSlice = append(studyStructSlice, StudyStruct{Title: title, Url: urlArticle, Authors: "", Abstract: ""})

				}
				if len(studyStructSlice) != 0 {
					return &studyStructSlice, true
				}
			}
		}
	}

	return nil, false
}
