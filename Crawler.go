package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var Max_Results = 2

func init() {
	if val, ok := os.LookupEnv("wiki_race_matches"); ok {
		num, err := strconv.Atoi(val)
		if err == nil {
			Max_Results = num
		} else {
			panic("Cannot convert wiki_race_matches to int: " + err.Error())
		}
	}
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

type Crawler struct {
	match   string
	pool    *WorkPool
	result  chan *PathNode
	visited sync.Map
}

func NewCrawler(match string, pool *WorkPool) *Crawler {
	return &Crawler{
		match:  match,
		pool:   pool,
		result: make(chan *PathNode),
	}
}

func (crawler *Crawler) WaitForResult() []*PathNode {
	defer crawler.pool.Stop()

	var results = make([]*PathNode, Max_Results)
	for i := 0; i < Max_Results; i++ {
		results[i] = <-crawler.result
		log.Printf("Result %v: %v\n", i+1, results[i])
	}

	return results
}

func (crawler *Crawler) Start(title string) *Crawler {
	crawler.pool.Start()
	crawler.pool.AddJob(Job{func() {
		crawler.visited.Store(title, true)
		crawl(crawler, (&PathNode{}).New(title))
	}})
	return crawler
}

func crawl(crawler *Crawler, node *PathNode) {
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
			log.Println("Found empty string")
			continue
		}
		if nestedLink == crawler.match {
			fmt.Printf("Found %v on %v\n", crawler.match, node.name)
			crawler.result <- node.New(nestedLink)
			return
		}
		if _, loaded := crawler.visited.LoadOrStore(nestedLink, true); !loaded {
			// fmt.Printf("Creating new job for %v\n", nestedLink)
			newTitle := nestedLink
			go crawler.pool.AddJob(Job{func() {
				crawl(crawler, node.New(newTitle))
			}})
		}
		// else {
		// 	// fmt.Printf("Found visited link %v\n", nestedLink)
		// }
	}
	// fmt.Println("End of crawl function")
}
