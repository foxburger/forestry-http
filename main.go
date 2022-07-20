package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	url, _ := url.Parse(os.Getenv("FORESTRY_URL"))
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(r *http.Response) error {
		if r.StatusCode == http.StatusNotFound {
			return errors.New("NOT_FOUND")
		}

		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}

	http.HandleFunc("/http-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.Host = req.URL.Host

		proxy.ServeHTTP(w, req)
	})

	http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 0

	http.NotFoundHandler()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
