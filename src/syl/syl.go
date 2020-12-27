package syl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/html"
)

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
func Run(strDur, url string, opts *OptionFlags) {
	d, err := time.ParseDuration(strDur)
	if err != nil {
		printRed(errInvalidDuration)
		return
	}

	resp, err := http.Get(url)
	fmt.Println(resp.StatusCode)
	if err != nil {
		printRed(err)
		return
	}
	defer resp.Body.Close()
	title, ok := getHTMLTitle(resp.Body)
	if !ok {
		title = url
	}
	fmt.Printf("See you later! %s 👋", title)

	// prepare for exit from a client side
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-time.After(d):
		err := openURLWithBrowser(url)
		if err != nil {
			printRed("Oops! We cannot open your link \n: %v", err)
			return
		}
		printGreen("Happy to see you! I hope you enjoy 🎉")
	case <-exit:
		printGreen("Goodbye 👋")
		os.Exit(1)
	}

}

func openURLWithBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func getHTMLTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
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
