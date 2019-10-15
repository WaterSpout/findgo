package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	if !isTerminalStdin() {
		return
	}
	total := getTotalGoCount()
	fmt.Printf("Total: %v\n", total)
}

func getTotalGoCount() (total int) {
	var url string
	var urls []string
	var divided [][]string
	k := 5
	var totalCount int
	for {
		_, err := fmt.Fscan(os.Stdin, &url)
		if len(urls) == k || (len(urls) != 0 && err != nil) {
			divided = append(divided, urls)
			urls = nil
		}
		if err != nil {
			break
		}
		totalCount++
		urls = append(urls, url)
	}
	all := make(chan int, totalCount)
	var wg sync.WaitGroup

	for _, urls := range divided {
		wg.Add(len(urls))
		for _, url := range urls {
			go search(url, &wg, all)
		}
		wg.Wait()
	}
	close(all)
	for cnt := range all {
		total += cnt
	}
	return total
}

func isTerminalStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func search(url string, wg *sync.WaitGroup, all chan<- int) {
	defer wg.Done()
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	reg := regexp.MustCompile("Go")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	matches := reg.FindAllStringIndex(string(body), -1)
	fmt.Printf("Count for %s: %v\n", url, len(matches))
	all <- len(matches)
}
