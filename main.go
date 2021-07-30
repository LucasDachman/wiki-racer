package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var Num_Workers = 8

var wiki *WikiService

func init() {
	if val, ok := os.LookupEnv("wiki_race_workers"); ok {
		num, err := strconv.Atoi(val)
		if err == nil {
			Num_Workers = num
		} else {
			panic("Cannot convert wiki_race_workers to int: " + err.Error())
		}
	}
}

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
	log.Println("Workers", Num_Workers)

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

func race(title1, title2 string) {
	pool := NewWorkPool(Num_Workers)

	crawler := NewCrawler(title2, &pool).Start(title1)
	results := crawler.WaitForResult()

	fmt.Println()
	for _, r := range results {
		printPath(r)
	}
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
