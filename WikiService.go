package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type WikiService struct {
	http       *http.Client
	actionBase string
	restBase   string
}

type RandomPageTitleResponseBody struct {
	Titles struct {
		Canonical string `json:"canonical"`
	} `json:"titles"`
}

type QueryLinksResponseBody struct {
	Query struct {
		Pages map[string]struct {
			Title        string `json:"title"`
			CanonicalUrl string `json:"canonicalUrl"`
		} `json:"pages"`
	} `json:"query"`
	Continue WikiContinue `json:"continue"`
}

type WikiContinue struct {
	Gplcontinue string `json:"gplcontinue"`
	Continue    string `json:"continue"`
}

func NewWikiService() *WikiService {
	return &WikiService{
		http:       &http.Client{},
		restBase:   `https://en.wikipedia.org/api/rest_v1`,
		actionBase: `https://en.wikipedia.org/w/api.php`,
	}
}

//curl -X GET "https://en.wikipedia.org/api/rest_v1/page/random/title" -H  "accept: application/problem+json"
func (service *WikiService) RandomPageTitle() string {
	resp, err := service.http.Get(service.restBase + `/page/random/summary`)
	if err != nil {
		panic(`Unable to fetch random page`)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(`Unable to read response body for random page`)
	}

	var jsonBody RandomPageTitleResponseBody
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		panic(`Unable to unmarshal json response`)
	}
	return jsonBody.Titles.Canonical
}

// https://en.wikipedia.org/w/api.php?action=query&titles=Albert%20Einstein&prop=links
func (service WikiService) ListLinks(title string, cont WikiContinue) ([]string, error) {
	myUrl, err := url.Parse(service.actionBase)
	if err != nil {
		panic(`Cannot parse action base`)
	}
	q := myUrl.Query()
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("generator", "links")
	q.Add("titles", title)
	q.Add("redirects", "")
	q.Add("prop", "info")
	q.Add("inprop", "url")
	q.Add("gpllimit", "max")

	if cont.Continue != "" {
		q.Add("gplcontinue", cont.Gplcontinue)
		q.Add("continue", cont.Continue)
	}
	myUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, myUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-Wiki-Race/0.0.1 (lucas.dachman@gmail.com)")

	resp, err := service.http.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(body))
	var jsonBody QueryLinksResponseBody
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println(string(body))
		return nil, err
	}

	// fmt.Println(jsonBody)
	links := make([]string, len(jsonBody.Query.Pages))

	// parse response body and aggregate links
	i := 0
	for _, page := range jsonBody.Query.Pages {
		reg := regexp.MustCompile(`\/wiki\/`)
		url := reg.Split(page.CanonicalUrl, 2)[1]
		links[i] = url
		i++
	}

	if jsonBody.Continue.Continue != "" {
		moreLinks, err := service.ListLinks(title, jsonBody.Continue)
		if err != nil {
			return nil, err
		}
		links = append(links, moreLinks...)
	}

	return links, nil
}
