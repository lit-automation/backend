package bibtex

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/wimspaargaren/slr-automation/src/packages/httpclient"
)

type Client interface {
	GetBibTex(doi string) (string, error)
}

type DXDOIClient struct {
	URL        string
	HTTPClient httpclient.HTTPClient
}

func (d *DXDOIClient) GetBibTex(doi string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", d.URL, doi), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/x-bibtex;q=1")
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.WithError(err).Errorf("unable to close body")
		}
	}()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	} else {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func NewDXDOIClient(httpClient httpclient.HTTPClient) Client {
	return &DXDOIClient{
		URL:        "http://dx.doi.org/",
		HTTPClient: httpClient,
	}
}
