package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type WikiService struct {
	http     *http.Client
	restBase string
}

type RandomPageTitleResponseBody struct {
	Items []struct {
		Title string `json:"title"`
	} `json:"items"`
}

func NewWikiService() *WikiService {
	return &WikiService{
		http:     &http.Client{},
		restBase: `https://en.wikipedia.org/api/rest_v1`,
	}
}

func (service *WikiService) GetCoolName() string {
	return "Lucas"
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
