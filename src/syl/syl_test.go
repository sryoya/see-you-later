package syl

import (
	"bytes"
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
		"non http protocol": {
			durStr: "1s",
			url:    "www.google.com",

			wantOutput: "unsupported protocol scheme \"\"",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			buffer.Reset()
			Run(tc.durStr, tc.url, tc.opts)
			if !strings.Contains(buffer.String(), tc.wantOutput) {
				t.Errorf("unexpected output message, want:%s, got: %s", tc.wantOutput, buffer.String())
			}
		})
	}
}
