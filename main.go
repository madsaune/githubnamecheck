package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type HttpStatus struct {
	StatusCode int
	Status     string
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "path", "./urls.txt", "path to urls.txt")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	urls, err := getUrls(filePath)
	if err != nil {
		log.Fatalf("could not get urls from file: %v", err)
	}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// Sleep for 0-5 seconds to avoid rate limiting
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			status, err := ping(url)
			if err != nil {
				log.Printf("err: failed to ping %s: %v\n", url, err)
			}

			if status.StatusCode == 404 {
				log.Printf("404: %s was not found\n", url)
				return
			}

			log.Printf("%s: %s\n", status.Status, url)
		}(url)
	}

	wg.Wait()
}

func ping(url string) (*HttpStatus, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return &HttpStatus{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}, nil
}

func getUrls(path string) ([]string, error) {
	var urls []string

	f, err := os.Open("./urls.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		urls = append(urls, s.Text())
	}

	return urls, nil
}
