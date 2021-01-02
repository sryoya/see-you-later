# see-you-later

A CLI tool that keeps the site you want to see and opens after the specified duration 👋

![syl-example](https://user-images.githubusercontent.com/14014676/103462484-d034f080-4d68-11eb-8b42-5b309f855b00.gif)

## ❓ What is this for?

This tool is assumued to be used, for example, when you have a article to view but you are busy on your current task and cannot view that now.
In the same way, you can also use this for like a Slack message you want to reply later, a Google Spreadsheet you want to work on later...etc.

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
$ syl 30min https://play.golang.org/
```

