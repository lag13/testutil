package testutil_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/lag13/testutil"
)

// TestCheckErrHasMsg checks that when we check an error for the
// expected message we get the expected diff.
func TestCheckErrHasMsg(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		msg      string
		wantDiff string
	}{
		{
			name:     "nil error empty message",
			err:      nil,
			msg:      "",
			wantDiff: "",
		},
		{
			name:     "nil error non-empty message",
			err:      nil,
			msg:      "some error happened",
			wantDiff: "got error message:\n  <nil>\nwant error message to start with the string:\n  some error happened",
		},
		{
			name:     "non-nil error empty message",
			err:      errors.New("some error"),
			msg:      "",
			wantDiff: "got non-nil error: some error",
		},
		{
			name:     "non-nil error non-empty message",
			err:      errors.New("some error: other error stuff"),
			msg:      "some error",
			wantDiff: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diff := testutil.CheckErrHasMsg(test.err, test.msg)
			if got, want := diff, test.wantDiff; got != want {
				t.Errorf("got wrong diff:\n###GOT###\n%s\n###WANT###\n%s", got, want)
			}
		})
	}
}

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
	req := testutil.MustNewHTTPRequest(wantReq.Method, wantReq.URL, strings.NewReader(wantReq.Body))
	req.Header.Add("Header1", "hey there")
	req.Header.Add("Header2", "pretty momma")
	req.Header.Add("Extra", "an extra header")
	if diff := testutil.CompareHTTPRequests(testutil.HTTPReqToTestutilHTTPReq(req), wantReq); diff != "" {
		t.Error(diff)
	}
}
