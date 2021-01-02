package syl

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fatih/color"
)

var buffer *bytes.Buffer

func init() {
	buffer = &bytes.Buffer{}
	color.Output = buffer
	writer = buffer

	startCmd = func(c *exec.Cmd) error { return nil }

	exit = func() { return }
}

type mockHTTPResponse struct {
	statusCode  int
	contentType string
	response    string
}

func writeMockedHTTPResponse(t *testing.T, w http.ResponseWriter, res *mockHTTPResponse) {
	// case for not 200
	if res.statusCode != http.StatusOK {
		http.Error(w, res.response, res.statusCode)
	}

	w.Header().Set("Content-Type", res.contentType)
	w.Write(([]byte)(res.response))
}

var testHTML = `
<!DOCTYPE html>
<html>
    <meta charset="utf-8"/>
    <head>
        <title>SYL</title>
    </head>
</html>
`

var testHTMLNoTitle = `
<!DOCTYPE html>
<html>
    <meta charset="utf-8"/>
</html>
`

func TestRun(t *testing.T) {
	cases := map[string]struct {
		// input
		durStr string
		url    string       // optional: don't need to specify when using httptest server
		opts   *OptionFlags // not used now

		// mock setting
		mockTargetSiteRes *mockHTTPResponse

		// output
		wantOutput []string
	}{
		"success": {
			durStr: "1ms",

			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  200,
				contentType: "text/html; charset=utf-8",
				response:    testHTML,
			},

			wantOutput: []string{"Hope to see you later! SYL 👋",
				"Happy to see you! I hope you enjoy",
			},
		},
		"suceess but non title is gotten by non http protocol": {
			durStr: "1ms",
			url:    "www.google.com", // important

			wantOutput: []string{"Hope to see you later!",
				"Happy to see you! I hope you enjoy",
				"unsupported protocol scheme", // important
			},
		},
		"suceess but non title is gotten by non existing url": {
			durStr: "1ms",
			url:    "http://not-existing-exsting",
			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  200,
				contentType: "text/html; charset=utf-8",
				response:    testHTML,
			},
			wantOutput: []string{"Hope to see you later! http://not-existing-exsting 👋",
				"Happy to see you! I hope you enjoy",
			},
		},
		"success but non title is gotten by non html response": {
			durStr: "1ms",

			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  200,
				contentType: "application/json",
				response:    `{"title": "syl"}`,
			},

			wantOutput: []string{"Hope to see you later!",
				"Happy to see you! I hope you enjoy",
			},
		},
		"suceess but non title is gotten by non successfully http response code": {
			durStr: "1ms",

			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  404,
				contentType: "text/html; charset=utf-8",
				response:    testHTML,
			},

			wantOutput: []string{"Hope to see you later!",
				"Happy to see you! I hope you enjoy",
			},
		},
		"error: invalid duration string": {
			durStr: "🍣",

			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  200,
				contentType: "text/html; charset=utf-8",
				response:    testHTML,
			},

			wantOutput: []string{errInvalidDuration.Error()},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// set up the clean up after test
			defer func() {
				buffer.Reset()
				goOS = runtime.GOOS
			}()

			// prepare mock for external URL
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeMockedHTTPResponse(t, w, tc.mockTargetSiteRes)
			}))
			defer ts.Close()
			url := tc.url
			if tc.url == "" {
				url = ts.URL
			}

			// execute
			Run(tc.durStr, url, tc.opts)

			// evaluate result
			for _, s := range tc.wantOutput {
				if !strings.Contains(buffer.String(), s) {
					t.Errorf("unexpected output message, want:%s, got: %s", tc.wantOutput, buffer.String())
				}
			}
		})
	}
}

func TestPrepareCommand(t *testing.T) {
	// set up the clean up after test
	defer func() {
		goOS = runtime.GOOS
	}()

	cases := map[string]struct {
		os          string
		wantCmdArgs []string
		wantErr     error
	}{
		"linux": {
			os:          "linux",
			wantCmdArgs: []string{"xdg-open", "https://www.google.com/"},
		},
		"windows": {
			os:          "windows",
			wantCmdArgs: []string{"rundll32", "url.dll,FileProtocolHandler", "https://www.google.com/"},
		},
		"darwin": {
			os:          "darwin",
			wantCmdArgs: []string{"open", "https://www.google.com/"},
		},
		"unknown OS": {
			os:      "👽",
			wantErr: errUnsupportedOS,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// set up
			goOS = tc.os

			// execute
			res, err := prepareCommand("https://www.google.com/")

			// evaludate results
			if tc.wantCmdArgs != nil {
				if diff := cmp.Diff(tc.wantCmdArgs, res.Args); diff != "" {
					t.Errorf("response didn't match (-want / +got)\n%s", diff)
				}
			}
			if !errors.Is(tc.wantErr, err) {
				t.Errorf("unexpected error, want: %v, got: %v", tc.wantErr, err)
			}
		})
	}

}
