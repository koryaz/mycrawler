package main

import (
	"fmt"
	"net/http"
	"net/url"

	"../mycrawler/mycrawler"
)

func main() {

	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	f1 := mycrawler.New(mycrawler.HandlerFunc(handler))
	go f1.Start(c1, 1)
	f2 := mycrawler.New(mycrawler.HandlerFunc(handler))
	go f2.Start(c2, 2)
	f3 := mycrawler.New(mycrawler.HandlerFunc(handler))
	go f3.Start(c3, 3)

	var urlStr string
	chNum := 0
	for {
		fmt.Scanf("%s", &urlStr)
		switch chNum {
		case 0:
			c1 <- urlStr
		case 1:
			c2 <- urlStr
		case 2:
			c3 <- urlStr
		}
		chNum++
		if chNum >= 3 {
			chNum = 0
		}
	}

}

func handler(method string, url url.URL, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, method, url.String())
}
