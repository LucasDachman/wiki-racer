package main

import (
	"fmt"
	"sync"
)

var visited sync.Map
var wiki *WikiService

func main() {
	wiki = NewWikiService()

	title1 := wiki.RandomPageTitle()
	title2 := wiki.RandomPageTitle()

	fmt.Printf("Title1: %v, Title2: %v\n", title1, title2)

	race(title1, title2)
	// fmt.Println(visited)
}

func race(title1, title2 string) {
	fmt.Printf("Visiting %v\n", title1)
	links := wiki.ListLinks(title1, "")
	fmt.Printf("Got links for %v\n", title1)
	for _, link := range links {
		if link == title2 {
			fmt.Printf("Found on %v\n", title1)
			return
		}
		if _, ok := visited.LoadOrStore(link, true); !ok {
			race(link, title2)
		}
	}
}
