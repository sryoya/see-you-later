package syl

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

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
			durStr: "1s",

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
			durStr: "1s",
			url:    "www.google.com", // important

			wantOutput: []string{"Hope to see you later!",
				"Happy to see you! I hope you enjoy",
				"unsupported protocol scheme", // important
			},
		},
		"suceess but non title is gotten by non existing url": {
			durStr: "1s",
			url:    "http://not-existing-exsting",
			mockTargetSiteRes: &mockHTTPResponse{
				statusCode:  200,
				contentType: "text/html; charset=utf-8",
				response:    testHTML,
			},
			wantOutput: []string{"Hope to see you later! http://not-existing-exsting 👋",
				"Happy to see you! I hope you enjoy",
				"no such host", // important
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			buffer.Reset()

			// prepare mock for external URL
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeMockedHTTPResponse(t, w, tc.mockTargetSiteRes)
			}))
			defer ts.Close()

			url := tc.url
			if tc.url == "" {
				url = ts.URL
			}

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
