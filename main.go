package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const Num_Workers = 8

var visited sync.Map
var wiki *WikiService

func main() {

	closeLogger := setupLogger()
	defer closeLogger()

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
	log.Printf("Starting on: %v, looking for: %v\n", title1, title2)

	start := time.Now()

	race(title1, title2)

	duration := time.Since(start)
	fmt.Println("Time: ", duration)
	log.Println("Time: ", duration)
}

func setupLogger() func() error {
	// Open log file
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(f)
	// log.SetOutput(os.Stdout)

	return f.Close
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

func crawl(node *PathNode, match string, pool WorkPool, result chan<- *PathNode) {
	fmt.Printf("Visiting %v\n", node.name)
	links, err := wiki.ListLinks(node.name, WikiContinue{})
	// fmt.Printf("Finished %v\n", title1)
	// fmt.Printf("Found %v\n", links)
	if err != nil {
		log.Println("Failed visiting page: " + node.name + "\n" + err.Error())
		return
	}
	for _, nestedLink := range links {
		// fmt.Printf("Inspecting %v\n", nestedLink)
		if nestedLink == "" {
			fmt.Printf("Found empty string\n")
			continue
		}
		if nestedLink == match {
			fmt.Printf("Found %v on %v\n", match, node.name)
			result <- node.New(nestedLink)
			return
		}
		if _, loaded := visited.LoadOrStore(nestedLink, true); !loaded {
			// fmt.Printf("Creating new job for %v\n", nestedLink)
			newTitle := nestedLink
			go pool.AddJob(Job{func() {
				crawl(node.New(newTitle), match, pool, result)
			}})
		}
		// else {
		// 	// fmt.Printf("Found visited link %v\n", nestedLink)
		// }
	}
	// fmt.Println("End of crawl function")
}

func race(title1, title2 string) {
	result := make(chan *PathNode)

	pool := NewWorkPool(Num_Workers)
	pool.Start()

	pool.AddJob(Job{func() {
		// todo clean up
		visited.Store(title1, true)
		parent := &PathNode{}
		next := parent.New(title1)
		crawl(next, title2, pool, result)
	}})

	lastNode := <-result
	pool.Stop()

	fmt.Println()
	printPath(lastNode)
	fmt.Println("Jobs completed: ", pool.jobCounter.num)
}

func printPath(node *PathNode) {
	ptr := node
	var str string
	for ptr != nil {
		if ptr.name != "" {
			if ptr != node {
				str = " -> " + str
			}
			str = ptr.name + str
		}
		ptr = ptr.parent
	}
	fmt.Println(str)
	fmt.Printf("Found in %v visits\n", node.len)
	log.Println(str)
	log.Printf("Found in %v visits\n", node.len)
}
