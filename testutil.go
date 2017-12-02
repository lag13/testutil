// Package testutil contains testing utilities.
package testutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// MustNewHTTPRequest creates a new HTTP request suitable for sending
// as opposed to httptest.NewRequest which is only suitable for
// passing into a http.Handler. It panic's if the request cannot be
// created.
func MustNewHTTPRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return req
}

// MustSendHTTPRequest sends a HTTP request using the default HTTP
// client and panic's if the send fails.
func MustSendHTTPRequest(r *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustReadAll reads everything from a Reader and panic's if reading
// fails.
func MustReadAll(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// CompareStrings compares two strings and returns a string detailing
// where they differ. Useful for when two large strings need to be
// compared.
func CompareStrings(got string, want string) string {
	return ""
}

// HTTPRequest represents the fields of a HTTP request which I think
// are most important for comparing. It also can be marshalled to JSON
// with the intent that you can use it to check that your API sent the
// expected requests to another API in an end-to-end test. In practice
// that should look something like:
//
// 	- In test code trigger an endpoint on your API which talks to
// 	  a mock API XYZ.
//      - When XYZ receives the request it will record it.
//      - After your API is done the test code will hit an endpoint on
//        XYZ to return the requests it received. Then your test code
//        checks to make sure that the expected requests were sent.
type HTTPRequest struct {
	Header http.Header `json:"header"`
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Body   string      `json:"body"`
}

// HTTPReqToTestutilHTTPReq converts an *http.Request to a
// HTTPRequest. It exists because I started writing a comparison
// function for http.Request's and HTTPRequest's and thought that was
// silly since they were so similar.
func HTTPReqToTestutilHTTPReq(req *http.Request) HTTPRequest {
	return HTTPRequest{
		Header: req.Header,
		Method: req.Method,
		URL:    req.URL.String(),
		Body:   MustReadAll(req.Body),
	}
}

// CompareHTTPRequests checks two HTTPRequest types for equality (the
// exception being the HTTP Header where we just check that we
// produced the specific headers we're interested in and ignore any
// extras).
func CompareHTTPRequests(got HTTPRequest, want HTTPRequest) string {
	diffs := []string{}
	for headerName := range want.Header {
		if got, want := got.Header.Get(headerName), want.Header.Get(headerName); got != want {
			diffs = append(diffs, fmt.Sprintf("for header %q got value %q, want %q", headerName, got, want))
		}
	}
	if got, want := got.Method, want.Method; got != want {
		diffs = append(diffs, fmt.Sprintf("got method %q, want %q", got, want))
	}
	// TODO: Maybe use compare string function here?
	if got, want := got.URL, want.URL; got != want {
		diffs = append(diffs, fmt.Sprintf("got url:\n  %q\nwant:\n  %q", got, want))
	}
	// TODO: Use compare string function here.
	if got, want := got.Body, want.Body; got != want {
		diffs = append(diffs, fmt.Sprintf("got body:\n  %s\nwant:\n  %s", got, want))
	}
	if len(diffs) > 0 {
		return "request is not expected:\n" + strings.Join(diffs, "\n")
	}
	return ""
}
