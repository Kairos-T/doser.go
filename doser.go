package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	url            string
	payload        string
	threads        int
	requestCounter int
	printedMsgs    []string
	waitGroup      sync.WaitGroup
	userAgents     = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 8.0.0; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.137 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:61.0) Gecko/20100101 Firefox/61.0",
	}
)

func printMsg(msg string) {
	if !contains(printedMsgs, msg) {
		fmt.Printf("\n%s after %d requests\n", msg, requestCounter)
		printedMsgs = append(printedMsgs, msg)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleStatusCodes(statusCode int) {
	requestCounter++
	fmt.Printf("\r%d requests have been sent", requestCounter)

	if statusCode == 429 {
		printMsg("You have been throttled")
	}
	if statusCode == 500 {
		printMsg("Status code 500 received")
	}
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func sendGET() {
	defer waitGroup.Done()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", getRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	handleStatusCodes(resp.StatusCode)
}

func sendPOST() {
	defer waitGroup.Done()

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	handleStatusCodes(resp.StatusCode)
}

func main() {
	flag.StringVar(&url, "g", "", "Specify GET request. Usage: -g '<url>'")
	flag.StringVar(&url, "p", "", "Specify POST request. Usage: -p '<url>'")
	flag.StringVar(&payload, "d", "", "Specify data payload for POST request")
	flag.IntVar(&threads, "t", 500, "Specify number of threads to be used")
	flag.Parse()

	if url == "" {
		flag.Usage()
		return
	}

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if flag.NFlag() == 1 {
		fmt.Println("You must specify either a GET (-g) or POST (-p) request.")
		return
	}

	waitGroup.Add(threads)

	for i := 0; i < threads; i++ {
		if url != "" {
			if payload != "" {
				go sendPOST()
			} else {
				go sendGET()
			}
		}
	}
	waitGroup.Wait()
}
