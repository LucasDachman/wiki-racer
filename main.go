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

	// wiki.ListLinks("Albert_Einstein", WikiContinue{})
	// fmt.Println(links, err)

	title1, title2 := "Albert_Einstein", "Molecule"
	race(title1, title2)
	// fmt.Println(visited)
}

type Job struct {
	op func()
}

func worker(id int, jobs <-chan Job) {
	var jobNum = 0
	for job := range jobs {
		jobNum++
		// fmt.Println("Starting Job: ", jobNum)
		job.op()
		// fmt.Println("Finished Job: ", jobNum)
	}
}

func crawl(title1, title2 string, jobs chan<- Job, done chan<- bool) {
	fmt.Printf("Visiting %v\n", title1)
	links, err := wiki.ListLinks(title1, WikiContinue{})
	// fmt.Printf("Finished %v\n", title1)
	// fmt.Printf("Found %v\n", links)
	if err != nil {
		panic("Failed visiting title1: " + err.Error())
		// fmt.Println(err.Error())
		// return
	}
	for _, nestedLink := range links {
		// fmt.Printf("Inspecting %v\n", nestedLink)
		if nestedLink == "" {
			fmt.Printf("Found empty string\n")
			continue
		}
		if nestedLink == title2 {
			fmt.Printf("Found %v on %v\n", title2, title1)
			done <- true
			return
		}
		if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
			// fmt.Printf("Creating new job for %v\n", nestedLink)
			newTitle := nestedLink
			go (func() {
				jobs <- Job{func() {
					crawl(newTitle, title2, jobs, done)
				}}
			})()
		}
		// else {
		// 	// fmt.Printf("Found visited link %v\n", nestedLink)
		// }
	}
	// fmt.Println("End of crawl function")
}

const MaxJobs = 1

func race(title1, title2 string) {
	jobs := make(chan Job, MaxJobs)
	done := make(chan bool)

	for i := 0; i < MaxJobs; i++ {
		go worker(i, jobs)
	}

	jobs <- Job{func() {
		visited.Store(title1, true)
		crawl(title1, title2, jobs, done)
		// fmt.Println("Executing initial job")
		// done <- true
	}}

	<-done
}
