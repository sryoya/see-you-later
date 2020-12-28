package syl

import (
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/html"
)

var goOS = runtime.GOOS
var startCmd = func(c *exec.Cmd) error { return c.Start() }
var exit = func() { os.Exit(1) }

type reservation struct {
	siteInfo *siteInfo
	openTime time.Duration
}

// OptionFlags has option to execute
// TODO: use it
type OptionFlags struct {
	repeats    bool
	ignores404 bool
}

// siteInfo is an info of the website to open
type siteInfo struct {
	title string
	url   string
}

// Run parses the request and executes
func Run(durStr, url string, opts *OptionFlags) {
	d, err := time.ParseDuration(durStr)
	if err != nil {
		printRed(errInvalidDuration)
		return
	}

	// check provided URL and get link
	res, err := http.Get(url)
	if err != nil {
		printRed(err)
		return
	}
	defer res.Body.Close()
	title, ok := getHTMLTitle(res)
	if !ok {
		title = url
	}
	print("Hope to see you later! %s 👋\n", title)

	// prepare for exit from a client-side action
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-time.After(d):
		err := openURLWithBrowser(url)
		if err != nil {
			printRed("Oops! We cannot open your link \n: %v", err)
			return
		}
		printGreen("Happy to see you! I hope you enjoy 🎉")
		exit()
	case <-exitCh:
		printGreen("Goodbye 👋👋👋")
		exit()
	}
}

func openURLWithBrowser(url string) error {
	cmd, err := prepareCommand(url)
	if err != nil {
		return err
	}

	return startCmd(cmd)
}

func prepareCommand(url string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return nil, errUnsupportedOS
	}

	return cmd, nil
}

func getHTMLTitle(r *http.Response) (string, bool) {
	if r.StatusCode != http.StatusOK {
		return "", false
	}

	doc, err := html.Parse(r.Body)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func traverse(n *html.Node) (string, bool) {
	// check is title element
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}
