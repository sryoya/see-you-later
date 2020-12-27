package syl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

type reservation struct {
	siteInfo *siteInfo
	openTime time.Duration
}

// OptionFlags has option to execute
type OptionFlags struct {
	repeats bool
}

// siteInfo is an info of the website to open
type siteInfo struct {
	title string
	url   string
}

// Run parses the request and executes
func Run(url string, strDur string, opts *OptionFlags) {
	d, err := time.ParseDuration(strDur)
	if err != nil {
		printRed(errInvalidDuration)
		return
	}

	resp, err := http.Get(url)
	fmt.Println(resp.StatusCode)
	if err != nil {

	}
	defer resp.Body.Close()
	title, ok := getHtmlTitle(resp.Body)
	if ok {
		fmt.Println(title)
	}

	select {
	case <-time.After(d):
		openURLWithBrowser(url)
	}

}

func openURLWithBrowser(url string) {
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
	if err != nil {
		log.Fatal(err)
	}
}

func getHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
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

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}
