// Package testutil contains testing utilities.
package testutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckErrHasMsg checks that the received error contains the message
// we want.
func CheckErrHasMsg(err error, wantMsg string) string {
	if wantMsg == "" && err != nil {
		return fmt.Sprintf("got non-nil error: %v", err)
	} else if got, want := fmt.Sprintf("%v", err), wantMsg; wantMsg != "" && !strings.HasPrefix(got, want) {
		return fmt.Sprintf("got error message:\n  %s\nwant error message to start with the string:\n  %s", got, want)
	}
	return ""
}

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
// where they differ or "" if they don't. Useful for when two large
// strings need to be compared.
func CompareStrings(got string, want string) string {
	for i := range want {
		if i > len(got)-1 {
			return fmt.Sprintf("got a shorter string than what we wanted (characters match otherwise) and the missing characters are: %s", want[i:])
		}
		if got[i] != want[i] {
			return fmt.Sprintf("strings differ at index %d, from that index on:\n##### got string #####\n%s\n##### want string #####\n%s", i, got[i:], want[i:])
		}
	}
	if len(want) < len(got) {
		return fmt.Sprintf("got a longer string than what we wanted (characters match otherwise) and the extra characters are: %s", got[len(want):])
	}
	return ""
}

// HTTPRequest represents the fields of a HTTP request which I think
// are most important for checking in a unit test. It also can be
// marshalled to JSON the intent being that you can use it to check if
// your API sent the expected requests to another API in an end-to-end
// test. In practice that would look something like:
//
// 	- In test code trigger an endpoint on your API which talks to
// 	  a mock API XYZ.
//      - When XYZ receives the request it will record it.
//      - After your API is done the test code will call an endpoint
//        on XYZ to return the requests it received. Then your test
//        code checks to make sure that the expected requests were
//        sent.
type HTTPRequest struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

// CheckHTTPRequest checks to make sure that a http.Request has the
// fields we're looking for.
func CheckHTTPRequest(got *http.Request, want HTTPRequest) string {
	diffs := []string{}
	for headerName := range want.Header {
		if got, want := got.Header.Get(headerName), want.Header.Get(headerName); got != want {
			diffs = append(diffs, fmt.Sprintf("header %q got value %q, want %q", headerName, got, want))
		}
	}
	if got, want := got.Method, want.Method; got != want {
		diffs = append(diffs, fmt.Sprintf("got method %q, want %q", got, want))
	}
	if got, want := got.URL.String(), want.URL; got != want {
		diffs = append(diffs, fmt.Sprintf("got url:\n  %q\nwant:\n  %q", got, want))
	}
	if diff := CompareStrings(MustReadAll(got.Body), want.Body); diff != "" {
		diffs = append(diffs, "body is not expected, "+diff)
	}
	if len(diffs) > 0 {
		return "request does not match what is expected:\n" + strings.Join(diffs, "\n")
	}
	return ""
}

// HTTPResponse contains the fields on a http.Response we are
// interested in checking.
type HTTPResponse struct {
	StatusCode int
	Header     http.Header
	Body       string
}

// CheckHTTPResponse compares two *http.Responses for equailty. It
// only checks fields which I feel actually matter namely:
//
// 	- Status code
// 	- Header
// 	- Body
//
// It will probably get used in end-to-end tests to make sure that a
// response received from an API is expected.
func CheckHTTPResponse(gotResp *http.Response, wantResp HTTPResponse) string {
	diffs := []string{}
	if got, want := gotResp.StatusCode, wantResp.StatusCode; got != want {
		diffs = append(diffs, fmt.Sprintf("got status code %d, want %d", got, want))
	}
	for headerName := range wantResp.Header {
		if got, want := gotResp.Header.Get(headerName), wantResp.Header.Get(headerName); got != want {
			diffs = append(diffs, fmt.Sprintf("header %q got value %q, want %q", headerName, got, want))
		}
	}
	if diff := CompareStrings(MustReadAll(gotResp.Body), wantResp.Body); diff != "" {
		diffs = append(diffs, "body is not expected, "+diff)
	}
	if len(diffs) > 0 {
		return "response does not match what is expected:\n" + strings.Join(diffs, "\n")
	}
	return ""
}
