package tm_http_redirect

import (
	"net/http"
)

func redirect(res http.ResponseWriter, destination string, statusCode int) {
	res.Header().Set(DefaultRedirectionHeader, destination)
	res.WriteHeader(statusCode)
}

func (t *TmHttpRedirect) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url := req.URL.RequestURI()
	for _, r := range *t.rules {
		if r.FromMatcher.MatchString(url) {
			if destination := r.Handle(url); destination != nil {
				redirect(res, *destination, r.Code)
				return
			}
		}
	}
	// Default -> Go on!
	t.next.ServeHTTP(res, req)
}
