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

type IdStudyList struct {
	Header struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"header"`
	Esearchresult struct {
		Count          string   `json:"count"`
		Retmax         string   `json:"retmax"`
		Retstart       string   `json:"retstart"`
		Idlist         []string `json:"idlist"`
		Translationset []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"translationset"`
		Querytranslation string `json:"querytranslation"`
	} `json:"esearchresult"`
}

type StudyLinks struct {
	Header struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"header"`
	Linksets []struct {
		Dbfrom    string `json:"dbfrom"`
		Idurllist []struct {
			ID      string `json:"id"`
			Objurls []struct {
				URL struct {
					Value string `json:"value"`
				} `json:"url"`
				Iconurl struct {
					Value string `json:"value"`
				} `json:"iconurl"`
				Subjecttypes []interface{} `json:"subjecttypes"`
				Categories   []string      `json:"categories"`
				Attributes   []string      `json:"attributes"`
				Provider     struct {
					Name     string `json:"name"`
					Nameabbr string `json:"nameabbr"`
					ID       string `json:"id"`
				} `json:"provider"`
			} `json:"objurls"`
		} `json:"idurllist"`
	} `json:"linksets"`
}

type PubmedArticleSet struct {
	XMLName       xml.Name `xml:"PubmedArticleSet"`
	Text          string   `xml:",chardata"`
	PubmedArticle struct {
		Text            string `xml:",chardata"`
		MedlineCitation struct {
			Text   string `xml:",chardata"`
			Status string `xml:"Status,attr"`
			Owner  string `xml:"Owner,attr"`
			PMID   struct {
				Text    string `xml:",chardata"`
				Version string `xml:"Version,attr"`
			} `xml:"PMID"`
			DateCompleted struct {
				Text  string `xml:",chardata"`
				Year  string `xml:"Year"`
				Month string `xml:"Month"`
				Day   string `xml:"Day"`
			} `xml:"DateCompleted"`
			DateRevised struct {
				Text  string `xml:",chardata"`
				Year  string `xml:"Year"`
				Month string `xml:"Month"`
				Day   string `xml:"Day"`
			} `xml:"DateRevised"`
			Article struct {
				Text     string `xml:",chardata"`
				PubModel string `xml:"PubModel,attr"`
				Journal  struct {
					Text string `xml:",chardata"`
					ISSN struct {
						Text     string `xml:",chardata"`
						IssnType string `xml:"IssnType,attr"`
					} `xml:"ISSN"`
					JournalIssue struct {
						Text        string `xml:",chardata"`
						CitedMedium string `xml:"CitedMedium,attr"`
						Volume      string `xml:"Volume"`
						Issue       string `xml:"Issue"`
						PubDate     struct {
							Text  string `xml:",chardata"`
							Year  string `xml:"Year"`
							Month string `xml:"Month"`
						} `xml:"PubDate"`
					} `xml:"JournalIssue"`
					Title           string `xml:"Title"`
					ISOAbbreviation string `xml:"ISOAbbreviation"`
				} `xml:"Journal"`
				ArticleTitle string `xml:"ArticleTitle"`
				Pagination   struct {
					Text       string `xml:",chardata"`
					StartPage  string `xml:"StartPage"`
					EndPage    string `xml:"EndPage"`
					MedlinePgn string `xml:"MedlinePgn"`
				} `xml:"Pagination"`
				ELocationID struct {
					Text    string `xml:",chardata"`
					EIdType string `xml:"EIdType,attr"`
					ValidYN string `xml:"ValidYN,attr"`
				} `xml:"ELocationID"`
				Abstract struct {
					Text                 string `xml:",chardata"`
					AbstractText         string `xml:"AbstractText"`
					CopyrightInformation string `xml:"CopyrightInformation"`
				} `xml:"Abstract"`
				AuthorList struct {
					Text       string `xml:",chardata"`
					CompleteYN string `xml:"CompleteYN,attr"`
					Author     struct {
						Text            string `xml:",chardata"`
						ValidYN         string `xml:"ValidYN,attr"`
						LastName        string `xml:"LastName"`
						ForeName        string `xml:"ForeName"`
						Initials        string `xml:"Initials"`
						AffiliationInfo struct {
							Text        string `xml:",chardata"`
							Affiliation string `xml:"Affiliation"`
						} `xml:"AffiliationInfo"`
					} `xml:"Author"`
				} `xml:"AuthorList"`
				Language            string `xml:"Language"`
				PublicationTypeList struct {
					Text            string `xml:",chardata"`
					PublicationType []struct {
						Text string `xml:",chardata"`
						UI   string `xml:"UI,attr"`
					} `xml:"PublicationType"`
				} `xml:"PublicationTypeList"`
				ArticleDate struct {
					Text     string `xml:",chardata"`
					DateType string `xml:"DateType,attr"`
					Year     string `xml:"Year"`
					Month    string `xml:"Month"`
					Day      string `xml:"Day"`
				} `xml:"ArticleDate"`
			} `xml:"Article"`
			MedlineJournalInfo struct {
				Text        string `xml:",chardata"`
				Country     string `xml:"Country"`
				MedlineTA   string `xml:"MedlineTA"`
				NlmUniqueID string `xml:"NlmUniqueID"`
				ISSNLinking string `xml:"ISSNLinking"`
			} `xml:"MedlineJournalInfo"`
			ChemicalList struct {
				Text     string `xml:",chardata"`
				Chemical []struct {
					Text            string `xml:",chardata"`
					RegistryNumber  string `xml:"RegistryNumber"`
					NameOfSubstance struct {
						Text string `xml:",chardata"`
						UI   string `xml:"UI,attr"`
					} `xml:"NameOfSubstance"`
				} `xml:"Chemical"`
			} `xml:"ChemicalList"`
			CitationSubset  string `xml:"CitationSubset"`
			MeshHeadingList struct {
				Text        string `xml:",chardata"`
				MeshHeading []struct {
					Text           string `xml:",chardata"`
					DescriptorName struct {
						Text         string `xml:",chardata"`
						UI           string `xml:"UI,attr"`
						MajorTopicYN string `xml:"MajorTopicYN,attr"`
					} `xml:"DescriptorName"`
					QualifierName []struct {
						Text         string `xml:",chardata"`
						UI           string `xml:"UI,attr"`
						MajorTopicYN string `xml:"MajorTopicYN,attr"`
					} `xml:"QualifierName"`
				} `xml:"MeshHeading"`
			} `xml:"MeshHeadingList"`
			CoiStatement string `xml:"CoiStatement"`
		} `xml:"MedlineCitation"`
		PubmedData struct {
			Text    string `xml:",chardata"`
			History struct {
				Text          string `xml:",chardata"`
				PubMedPubDate []struct {
					Text      string `xml:",chardata"`
					PubStatus string `xml:"PubStatus,attr"`
					Year      string `xml:"Year"`
					Month     string `xml:"Month"`
					Day       string `xml:"Day"`
					Hour      string `xml:"Hour"`
					Minute    string `xml:"Minute"`
				} `xml:"PubMedPubDate"`
			} `xml:"History"`
			PublicationStatus string `xml:"PublicationStatus"`
			ArticleIdList     struct {
				Text      string `xml:",chardata"`
				ArticleId []struct {
					Text   string `xml:",chardata"`
					IdType string `xml:"IdType,attr"`
				} `xml:"ArticleId"`
			} `xml:"ArticleIdList"`
		} `xml:"PubmedData"`
	} `xml:"PubmedArticle"`
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
				title := pubmedArticleSet.PubmedArticle.MedlineCitation.Article.ArticleTitle
				urlArticle := fmt.Sprintf(
					"https://pubmed.ncbi.nlm.nih.gov/%s/",
					idStudyList.Esearchresult.Idlist[0],
				)
				//handle authors
				abstract := pubmedArticleSet.PubmedArticle.MedlineCitation.Article.Abstract.AbstractText
				return &StudyStruct{
					Title:    title,
					Url:      urlArticle,
					Authors:  "",
					Abstract: abstract,
				}, true

			}
		} else {
			return nil, false
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
			idStrings := strings.Join(idStudyList.Esearchresult.Idlist, ",")
			urlLinks := fmt.Sprintf("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/elink.fcgi?dbfrom=pubmed&id=%s&retmode=json&mindate%s&maxdate=2024&cmd=prlinks", idStrings, minYear)

			req, err := http.NewRequest("GET", urlLinks, nil)
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
				var studyLinks StudyLinks
				err := json.NewDecoder(resp.Body).Decode(&studyLinks)
				if err != nil {
					log.Printf("error translating json links response from PMC API %v", err)
					return nil, false
				}
				if len(studyLinks.Linksets) != 0 {
					for _, value := range studyLinks.Linksets[0].Idurllist {
                        if len(value.Objurls) != 0 {
                            log.Printf("Id %s URL %s",value.ID,  value.Objurls[0].URL.Value)
                        }
					}
				} else {
                    return nil, false
                }

			} else {
				return nil, false
			}

		} else {
			return nil, false
		}

	}

	return nil, false
}
