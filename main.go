package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(0, 0, 7).Format(http.TimeFormat)
)

func main() {

	originServerURL, err := url.Parse(os.Getenv("FORESTRY_URL"))
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Cache-Control", "max-age:290304000, public")
		rw.Header().Set("Last-Modified", cacheSince)
		rw.Header().Set("Expires", cacheUntil)

		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

		originServerResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, originServerResponse.Body)
	})

	log.Fatal(http.ListenAndServe(os.Getenv("FORESTRY_LISTEN"), reverseProxy))

}
