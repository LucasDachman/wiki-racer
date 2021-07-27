package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WikiService struct {
	http       *http.Client
	actionBase string
	restBase   string
}

type RandomPageTitleResponseBody struct {
	Items []struct {
		Title string `json:"title"`
	} `json:"items"`
}

type QueryLinksResponseBody struct {
	Query struct {
		Pages map[string]struct {
			Title string `json:"title"`
			Links []struct {
				Title string `json:"title"`
			} `json:"Links"`
		} `json:"pages"`
	} `json:"query"`
	Continue struct {
		Plcontinue string `json:"plcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
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
	resp, err := service.http.Get(service.restBase + `/page/random/title`)
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
	if len(jsonBody.Items) == 0 {
		panic(`No items found in response`)
	}
	return jsonBody.Items[0].Title
}

// https://en.wikipedia.org/w/api.php?action=query&titles=Albert%20Einstein&prop=links
func (service WikiService) ListLinks(title string, plcontinue string) []string {
	myUrl, err := url.Parse(service.actionBase)
	if err != nil {
		panic(`Cannot parse action base`)
	}
	q := myUrl.Query()
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("prop", "links")
	q.Add("titles", title)
	if plcontinue != "" {
		q.Add("plcontinue", plcontinue)
	}
	myUrl.RawQuery = q.Encode()

	resp, err := service.http.Get(myUrl.String())
	if err != nil {
		panic(`Error fetching title links`)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(`Error reading links body`)
	}

	// fmt.Println(string(body))
	var jsonBody QueryLinksResponseBody
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		panic(`Error parsing body to json ` + err.Error())
	}

	// fmt.Println(jsonBody)
	links := make([]string, len(jsonBody.Query.Pages))

	// parse response body and aggregate links
	for _, page := range jsonBody.Query.Pages {
		for _, link := range page.Links {
			links = append(links, link.Title)
		}
	}

	if jsonBody.Continue.Plcontinue != "" {
		moreLinks := service.ListLinks(title, jsonBody.Continue.Plcontinue)
		links = append(links, moreLinks...)
	}

	return links
}
