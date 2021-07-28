package main

import (
	"fmt"
	"sync"
)

var visited sync.Map
var wiki *WikiService

func main() {
	wiki = NewWikiService()

	// title1 := wiki.RandomPageTitle()
	// title2 := wiki.RandomPageTitle()

	// fmt.Printf("Title1: %v, Title2: %v\n", title1, title2)

	title1, title2 := "Albert_Einstein", "General_relativity"
	race(title1, title2)
	// fmt.Println(visited)
}

type Job struct {
	op func()
}

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		job.op()
	}
}

func crawl(title1, title2 string, jobs chan<- Job) {
	fmt.Printf("Visiting %v\n", title1)
	links, err := wiki.ListLinks(title1, "")
	fmt.Printf("Finished %v\n", title1)
	if err != nil {
		panic("Failed visiting title1: " + err.Error())
		// fmt.Println(err.Error())
		// return
	}
	for _, nestedLink := range links {
		if nestedLink == title2 {
			fmt.Printf("Found on %v\n", title1)
			close(jobs)
			return
		}
		if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
			newTitle := nestedLink
			op := func() {
				crawl(newTitle, title2, jobs)
			}
			newJob := Job{op}
			jobs <- newJob
		}
	}
}

const MaxJobs = 4

func race(title1, title2 string) {
	jobs := make(chan Job, 2)

	jobs <- Job{func() {
		visited.Store(title1, true)
		crawl(title1, title2, jobs)
	}}

	for i := 0; i < MaxJobs; i++ {
		worker(i, jobs)
	}
}
