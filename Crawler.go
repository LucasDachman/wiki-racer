package main

import (
	"fmt"
	"log"
	"sync"
)

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

func (crawler *Crawler) waitForResult() *PathNode {
	defer crawler.pool.Stop()
	return <-crawler.result
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
			fmt.Printf("Found empty string\n")
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
