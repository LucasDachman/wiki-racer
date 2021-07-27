package main

import "fmt"

// https://en.wikipedia.org/w/api.php?action=query&titles=Albert%20Einstein&prop=links
func main() {
	wiki := NewWikiService()
	title := wiki.RandomPageTitle()
	links := wiki.ListLinks(title, "")
	fmt.Println(title)
	fmt.Println(links)
}
