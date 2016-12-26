package main

import (
	"fmt"
	"net/http"
	"net/url"

	"../mycrawler/mycrawler"
)

func main() {
	f := mycrawler.New(mycrawler.HandlerFunc(handler))
	f.Start("http://golang.org")
}

func handler(method string, url url.URL, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, method, url.String())
}
