package apihandlers

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
	PubmedArticle []struct {
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
				Abstract struct {
					Text                 string `xml:",chardata"`
					AbstractText         string `xml:"AbstractText"`
					CopyrightInformation string `xml:"CopyrightInformation"`
				} `xml:"Abstract"`
				AuthorList struct {
					Text       string `xml:",chardata"`
					CompleteYN string `xml:"CompleteYN,attr"`
					Author     []struct {
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
			} `xml:"Article"`
			MedlineJournalInfo struct {
				Text        string `xml:",chardata"`
				Country     string `xml:"Country"`
				MedlineTA   string `xml:"MedlineTA"`
				NlmUniqueID string `xml:"NlmUniqueID"`
				ISSNLinking string `xml:"ISSNLinking"`
			} `xml:"MedlineJournalInfo"`
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

