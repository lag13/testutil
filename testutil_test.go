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

// TestCompareStrings tests that the expected diff is generated in
// different scenarios.
func TestCompareStrings(t *testing.T) {
	tests := []struct {
		name     string
		gotStr   string
		wantStr  string
		wantDiff string
	}{
		{
			name:     "strings differ at a character",
			gotStr:   "hello there",
			wantStr:  "hello theer buddy",
			wantDiff: "strings differ at index 9, from that index on:\n##### got string #####\nre\n##### want string #####\ner buddy",
		},
		{
			name:     "got string is longer but matches otherwise",
			gotStr:   "hello there!!",
			wantStr:  "hello there",
			wantDiff: "got a longer string than what we wanted (characters match otherwise) and the extra characters are: !!",
		},
		{
			name:     "want string is longer but matches otherwise",
			gotStr:   "hello there",
			wantStr:  "hello there!! Buddy!!",
			wantDiff: "got a shorter string than what we wanted (characters match otherwise) and the missing characters are: !! Buddy!!",
		},
		{
			name:     "got empty want empty",
			gotStr:   "",
			wantStr:  "",
			wantDiff: "",
		},
		{
			name:     "got empty want non-empty",
			gotStr:   "",
			wantStr:  "some non-empty string",
			wantDiff: "got a shorter string than what we wanted (characters match otherwise) and the missing characters are: some non-empty string",
		},
		{
			name:     "got non-empty want empty",
			gotStr:   "some non-empty string",
			wantStr:  "",
			wantDiff: "got a longer string than what we wanted (characters match otherwise) and the extra characters are: some non-empty string",
		},
		{
			name:     "non-empty strings match",
			gotStr:   "keep on the sunny side of life",
			wantStr:  "keep on the sunny side of life",
			wantDiff: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diff := testutil.CompareStrings(test.gotStr, test.wantStr)
			if got, want := diff, test.wantDiff; got != want {
				t.Errorf("got wrong diff:\n### GOT ###\n%s\n### WANT ###\n%s", got, want)
			}
		})
	}
}

// TestCompareHTTPRequest tests that the expected diff is generated
// when comparing HTTP requests.
func TestCompareHTTPRequest(t *testing.T) {
	gotReq := testutil.HTTPRequest{
		Header: http.Header{
			"Header1": {"some value"},
			"Header2": {"some other value"},
		},
		Method: "DELETE",
		URL:    "http://hello.com",
		Body:   "hello buddy!",
	}
	wantReq := testutil.HTTPRequest{
		Header: http.Header{
			"Header1": {"a different value"},
		},
		Method: "POST",
		URL:    "http://hello-there.com",
		Body:   "goodbye buddy!",
	}
	wantDiff := `request does not match what is expected:
header "Header1" got value "some value", want "a different value"
got method "DELETE", want "POST"
got url:
  "http://hello.com"
want:
  "http://hello-there.com"
body is not expected, strings differ at index 0, from that index on:
##### got string #####
hello buddy!
##### want string #####
goodbye buddy!`
	if diff := testutil.CompareStrings(testutil.CompareHTTPRequests(gotReq, wantReq), wantDiff); diff != "" {
		t.Error(diff)
	}
}
