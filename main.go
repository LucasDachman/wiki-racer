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
	c := make(chan string, 2)
	visited.Store(title1, true)
	c <- title1
	for link := range c {
		link := link
		go (func() {
			fmt.Printf("Visiting %v\n", link)
			links, err := wiki.ListLinks(link, "")
			fmt.Printf("Finished %v\n", link)
			if err != nil {
				panic("Failed visiting link: " + err.Error())
				// fmt.Println(err.Error())
				// return
			}
			for _, nestedLink := range links {
				if nestedLink == title2 {
					fmt.Printf("Found on %v\n", link)
					close(c)
					return
				}
				if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
					c <- nestedLink
				}
			}
		})()
	}
}
