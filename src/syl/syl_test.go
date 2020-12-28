package syl

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/fatih/color"
)

var buffer *bytes.Buffer

func init() {
	buffer = &bytes.Buffer{}
	color.Output = buffer
	writer = buffer
	startCmd = func(c *exec.Cmd) error { return nil }
}

func TestRun(t *testing.T) {
	cases := map[string]struct {
		durStr string
		url    string
		opts   *OptionFlags

		wantOutput string
	}{
		"success": {
			durStr: "1s",
			url:    "https://www.google.com/",

			wantOutput: "Hope to see you later! Google 👋\nHappy to see you! I hope you enjoy 🎉\n",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			buffer.Reset()
			Run(tc.durStr, tc.url, tc.opts)
			if bytes.Compare(buffer.Bytes(), []byte(tc.wantOutput)) != 0 {
				t.Errorf("unexpected output message, want:%s, got: %s", tc.wantOutput, buffer.String())
			}
		})
	}
}
