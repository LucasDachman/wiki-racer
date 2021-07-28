package main

import (
	"fmt"
	"os"
	"sync"
)

const MaxJobs = 4

var visited sync.Map
var wiki *WikiService

func main() {
	wiki = NewWikiService()

	var title1, title2 string
	if len(os.Args) == 1 {
		title1 = wiki.RandomPageTitle()
		title2 = wiki.RandomPageTitle()
	}
	if len(os.Args) == 2 {
		title1 = wiki.RandomPageTitle()
		title2 = os.Args[1]
	}
	if len(os.Args) == 3 {
		title1, title2 = os.Args[1], os.Args[2]
	}

	fmt.Printf("Starting on: %v, looking for: %v\n", title1, title2)
	race(title1, title2)
}

type Job struct {
	op func()
}

type PathNode struct {
	name   string
	parent *PathNode
	len    int
}

func (node *PathNode) New(name string) *PathNode {
	return &PathNode{
		name:   name,
		parent: node,
		len:    node.len + 1,
	}
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
	}
	for _, nestedLink := range links {
		// fmt.Printf("Inspecting %v\n", nestedLink)
		if nestedLink == "" {
			fmt.Printf("Found empty string\n")
			continue
		}
		if nestedLink == match {
			fmt.Printf("Found %v on %v\n", match, node.name)
			done <- node.New(nestedLink)
			return
		}
		if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
			// fmt.Printf("Creating new job for %v\n", nestedLink)
			newTitle := nestedLink
			go (func() {
				jobs <- Job{func() {
					crawl(node.New(newTitle), match, jobs, done)
				}}
			})()
		}
		// else {
		// 	// fmt.Printf("Found visited link %v\n", nestedLink)
		// }
	}
	// fmt.Println("End of crawl function")
}

func race(title1, title2 string) {
	jobs := make(chan Job, MaxJobs)
	done := make(chan *PathNode)

	for i := 0; i < MaxJobs; i++ {
		go worker(i, jobs)
	}

	jobs <- Job{func() {
		visited.Store(title1, true)
		parent := &PathNode{}
		next := parent.New(title1)
		crawl(next, title2, jobs, done)
	}}

	lastNode := <-done
	printPath(lastNode)
}

func printPath(node *PathNode) {
	ptr := node
	defer fmt.Printf("Found in %v visits\n", node.len)
	for ptr != nil {
		defer fmt.Print(ptr.name + " -> ")
		ptr = ptr.parent
	}
}
