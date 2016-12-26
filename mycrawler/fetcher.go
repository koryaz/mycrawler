package mycrawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Fetcher struct {
	Handler        Handler
	HTTPClient     Doer
	urlsnonvisited []string
	urlsvisited    []string
	baseURL        *url.URL
}

// New returns an initialized Fetcher.
func New(h Handler) *Fetcher {
	return &Fetcher{
		Handler:    h,
		HTTPClient: http.DefaultClient,
	}
}

func (f *Fetcher) Start(s string) {
	fmt.Printf("Fetcher started\n")
	fmt.Printf("First url to fetch: %s\n", s)
	var err error
	f.baseURL, err = url.Parse(s)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	f.urlsnonvisited = append(f.urlsnonvisited, s)
	for {
		//Quit loop if no more urls
		if len(f.urlsnonvisited) <= 0 {
			fmt.Printf("No more urls to fetch\n")
			break
		}

		//Do request for next non visited url and add url to visited list
		nextRawURL, nextLinkExist := f.getNextLinkToRequest()
		if !nextLinkExist {
			fmt.Printf("No more links. End of execution\n")
			break
		}
		fmt.Printf("Next url to fetch: %s\n", nextRawURL)
		parsedURL, err := url.Parse(nextRawURL)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		} else {
			res, err := f.doRequest(parsedURL)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			} else {
				fmt.Printf("Url fetched sucessfully: %s\n", nextRawURL)

				//Handle response body and add new links to non visited list
				//Handler.Handle(parsedURL.String(), parsedURL, res, err)
				f.parseLinksInResponseBody(&res.Body)
			}
		}
	}
}

func (f *Fetcher) doRequest(url *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	// Do the request.
	res, err := f.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (f *Fetcher) parseLinksInResponseBody(body *io.ReadCloser) {
	fmt.Printf("parseLinksInResponse\n")
	doc, _ := goquery.NewDocumentFromReader(*body)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exist := s.Attr("href")
		if exist != false {
			parsedURL, err := url.Parse(link)
			if err != nil {
				fmt.Printf("%s\n", err)
			} else {
				if parsedURL.IsAbs() {
					fmt.Printf("[%s] added to non visited url list\n", link)
					f.urlsnonvisited = append(f.urlsnonvisited, link)
				} else {
					if f.convertRelativeUrlToAbsolute(parsedURL) {
						fmt.Printf("[%s] added to non visited url list\n", link)
						f.urlsnonvisited = append(f.urlsnonvisited, link)
					} else {
						fmt.Printf("relative url [%s] could not be converted to absolute url\n", link)
					}
				}
			}
		}
	})
}

func (f *Fetcher) getNextLinkToRequest() (nextLink string, nextLinkExist bool) {
	nextLinkExist = false
	if len(f.urlsnonvisited) <= 0 {
		nextLink = ""
	} else if len(f.urlsnonvisited) == 1 {
		nextLink = f.urlsnonvisited[0]
		f.urlsnonvisited = f.urlsnonvisited[:0]
		f.urlsvisited = append(f.urlsvisited, nextLink)
		nextLinkExist = true
	} else {
		nextLink = f.urlsnonvisited[0]
		f.urlsnonvisited = f.urlsnonvisited[1:len(f.urlsnonvisited)]
		f.urlsvisited = append(f.urlsvisited, nextLink)
		nextLinkExist = true
	}
	return nextLink, nextLinkExist
}

func (f *Fetcher) convertRelativeUrlToAbsolute(url *url.URL) (isConversionDone bool) {
	fmt.Printf("convertRelativeUrlToAbsolute\n")
	isConversionDone = false
	if url.IsAbs() {
		isConversionDone = true
	} else {
		baseUrlStr := f.baseURL.String()
		if url.Scheme == "" {
			if url.String() == "" || url.String()[0] != '/' {
				// make relative path absolute
				baseUrlDir, _ := path.Split(baseUrlStr)
				urlStr := baseUrlDir + url.String()
				absoluteURL, err := url.Parse(urlStr)
				if err != nil {
					*url = *absoluteURL
					isConversionDone = true
				}
			}
		}
	}
	return isConversionDone
}
