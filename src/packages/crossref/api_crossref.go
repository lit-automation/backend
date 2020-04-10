package crossref

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wimspaargaren/slr-automation/src/packages/httpclient"
)

type ArticleInfo struct {
	Indexed struct {
		DateParts [][]int   `json:"date-parts"`
		DateTime  time.Time `json:"date-time"`
		Timestamp int64     `json:"timestamp"`
	} `json:"indexed"`
	ReferenceCount int    `json:"reference-count"`
	Publisher      string `json:"publisher"`
	Issue          string `json:"issue"`
	License        []struct {
		URL   string `json:"URL"`
		Start struct {
			DateParts [][]int   `json:"date-parts"`
			DateTime  time.Time `json:"date-time"`
			Timestamp int64     `json:"timestamp"`
		} `json:"start"`
		DelayInDays    int    `json:"delay-in-days"`
		ContentVersion string `json:"content-version"`
	} `json:"license"`
	ContentDomain struct {
		Domain               []interface{} `json:"domain"`
		CrossmarkRestriction bool          `json:"crossmark-restriction"`
	} `json:"content-domain"`
	ShortContainerTitle []string `json:"short-container-title"`
	PublishedPrint      struct {
		DateParts [][]int `json:"date-parts"`
	} `json:"published-print"`
	DOI      string `json:"DOI"`
	Abstract string `json:"abstract"`
	Type     string `json:"type"`
	Created  struct {
		DateParts [][]int   `json:"date-parts"`
		DateTime  time.Time `json:"date-time"`
		Timestamp int64     `json:"timestamp"`
	} `json:"created"`
	Page                string   `json:"page"`
	Source              string   `json:"source"`
	IsReferencedByCount int      `json:"is-referenced-by-count"`
	Title               []string `json:"title"`
	Prefix              string   `json:"prefix"`
	Volume              string   `json:"volume"`
	Author              []struct {
		Given       string        `json:"given"`
		Family      string        `json:"family"`
		Sequence    string        `json:"sequence"`
		Affiliation []interface{} `json:"affiliation"`
	} `json:"author"`
	Member         string   `json:"member"`
	ContainerTitle []string `json:"container-title"`
	Language       string   `json:"language"`
	Link           []struct {
		URL                 string `json:"URL"`
		ContentType         string `json:"content-type"`
		ContentVersion      string `json:"content-version"`
		IntendedApplication string `json:"intended-application"`
	} `json:"link"`
	Deposited struct {
		DateParts [][]int   `json:"date-parts"`
		DateTime  time.Time `json:"date-time"`
		Timestamp int64     `json:"timestamp"`
	} `json:"deposited"`
	Score  float64 `json:"score"`
	Issued struct {
		DateParts [][]int `json:"date-parts"`
	} `json:"issued"`
	ReferencesCount int `json:"references-count"`
	JournalIssue    struct {
		PublishedPrint struct {
			DateParts [][]int `json:"date-parts"`
		} `json:"published-print"`
		Issue string `json:"issue"`
	} `json:"journal-issue"`
	AlternativeID []string `json:"alternative-id"`
	URL           string   `json:"URL"`
	ISSN          []string `json:"ISSN"`
	IssnType      []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"issn-type"`
}

// WorksResult works result contains a list of item for given query
type WorksResult struct {
	Status         string `json:"status"`
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        struct {
		Facets struct {
		} `json:"facets"`
		TotalResults int           `json:"total-results"`
		Items        []ArticleInfo `json:"items"`
		ItemsPerPage int           `json:"items-per-page"`
		Query        struct {
			StartIndex  int    `json:"start-index"`
			SearchTerms string `json:"search-terms"`
		} `json:"query"`
	} `json:"message"`
}

type WorksResultDOI struct {
	Status         string      `json:"status"`
	MessageType    string      `json:"message-type"`
	MessageVersion string      `json:"message-version"`
	Message        ArticleInfo `json:"message"`
}

type Client interface {
	QueryWorks(query string) (*WorksResult, error)
	GetOnDOI(doi string) (*WorksResultDOI, error)
}

type APIClient struct {
	URL        string
	HTTPClient httpclient.HTTPClient
}

const (
	WorksPrefix string = "works"
)

func (a *APIClient) GetOnDOI(doi string) (*WorksResultDOI, error) {
	doiURLEncoded := url.QueryEscape(doi)
	resp, err := a.HTTPClient.Get(fmt.Sprintf("%s%s/%s", a.URL, WorksPrefix, doiURLEncoded))
	if err != nil {
		return nil, err
	}
	var result WorksResultDOI
	err = decodeBody(resp, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *APIClient) QueryWorks(query string) (*WorksResult, error) {
	queryURLEncoded := url.QueryEscape(query)
	resp, err := a.HTTPClient.Get(fmt.Sprintf("%s%s?query=%s", a.URL, WorksPrefix, queryURLEncoded))
	if err != nil {
		return nil, err
	}
	var result WorksResult
	err = decodeBody(resp, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func decodeBody(resp *http.Response, data interface{}) error {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.WithError(err).Errorf("unable to close body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &data)
}

func NewCrossRefClient(httpClient httpclient.HTTPClient) Client {
	return &APIClient{
		URL:        "https://api.crossref.org/",
		HTTPClient: httpClient,
	}
}
