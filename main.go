package main

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "regexp"
)

func main() {
    target, err := url.Parse("http://localhost:8000")
    if err != nil {
        panic(err)
    }

    proxy := httputil.NewSingleHostReverseProxy(target)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // If the URL contains "/ref=", then proxy to the target URL
        refRegex := regexp.MustCompile(`^/ref=(.*)`)
        if refRegex.MatchString(r.URL.Path) {
            refURL := refRegex.ReplaceAllString(r.URL.Path, "$1")
            refTarget, err := url.Parse(refURL)
            if err != nil {
                http.Error(w, fmt.Sprintf("invalid ref URL: %s", refURL), http.StatusBadRequest)
                return
            }
            proxy.Director = func(req *http.Request) {
                req.URL.Scheme = refTarget.Scheme
                req.URL.Host = refTarget.Host
                req.URL.Path = refTarget.Path
            }
        } else {
            // Serve the original page
            proxy.Director = nil
        }
        proxy.ServeHTTP(w, r)
    })

    if err := http.ListenAndServe(":80", nil); err != nil {
        panic(err)
    }
}