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

	title1, title2 := "Albert_Einstein", "General_relativity"
	race(title1, title2)
	// fmt.Println(visited)
}

type Job struct {
	op func()
}

type PathNode struct {
	name   string
	parent *PathNode
	len    int
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

func crawl(node *PathNode, match string, jobs chan<- Job, done chan<- *PathNode) {
	fmt.Printf("Visiting %v\n", node.name)
	links, err := wiki.ListLinks(node.name, WikiContinue{})
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
		if nestedLink == match {
			fmt.Printf("Found %v on %v\n", match, node.name)
			done <- &PathNode{name: nestedLink, parent: node, len: node.parent.len + 1}
			return
		}
		if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
			// fmt.Printf("Creating new job for %v\n", nestedLink)
			newTitle := nestedLink
			go (func() {
				jobs <- Job{func() {
					newNode := &PathNode{name: newTitle, parent: node, len: node.len + 1}
					crawl(newNode, match, jobs, done)
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
	done := make(chan *PathNode)

	for i := 0; i < MaxJobs; i++ {
		go worker(i, jobs)
	}

	jobs <- Job{func() {
		visited.Store(title1, true)
		crawl(&PathNode{name: title1, len: 1, parent: &PathNode{}}, title2, jobs, done)
	}}

	lastNode := <-done
	printPath(lastNode)
}

func printPath(node *PathNode) {
	ptr := node
	for ptr != nil {
		defer fmt.Print(ptr.name + " -> ")
		ptr = ptr.parent
	}
}
