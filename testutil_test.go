package testutil_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/lag13/testutil"
)

// TestMustNewHTTPRequestAndFriends tests that the request gets
// created as expected and uses some other helper functions along the
// way!
func TestMustNewHTTPRequestAndFriends(t *testing.T) {
	wantReq := testutil.HTTPRequest{
		Header: http.Header{
			"Header1": {"hey there"},
			"Header2": {"pretty momma"},
		},
		Method: "GET-OUTTA-HERE",
		URL:    "http://hello.com/woweee?hello=world",
		Body:   "hello",
	}
	req := testutil.MustNewHTTPRequest(wantReq.Method+"asd", wantReq.URL+"fdsa", strings.NewReader(wantReq.Body+"asf"))
	req.Header.Add("Header1", "hey there")
	req.Header.Add("Header2", "pretty momma")
	req.Header.Add("Extra", "an extra header")
	if diff := testutil.CompareHTTPRequests(testutil.HTTPReqToTestutilHTTPReq(req), wantReq); diff != "" {
		t.Error(diff)
	}
	t.Error("another error")
}
