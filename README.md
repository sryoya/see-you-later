# see-you-later

A CLI tool that keeps the site you want to see and opens after the specified duration 👋

![syl-example](https://user-images.githubusercontent.com/14014676/104760657-02047900-57a5-11eb-9861-0d789298828d.gif)

## ❓ What is this for?
Have you experienced somethings like these? 
- You find the article, movies..etc that you want to view, but you have a task to do right now :dash:
- You receive the Slack notification that seems to need the response, but you are busy on the coding :man::computer:

In such cases, you can keep the link and view it later by this tool as long as there are a CLI in front of you.

## 🚀  Install

```bash
$ go install github.com/sryoya/see-you-later/cmd/syl

```

## 🏁 How to use

```bash
$ syl [duration] [URL]
```

example
```bash
$ syl 30m https://play.golang.org/
```

