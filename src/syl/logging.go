package syl

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var writer io.Writer = os.Stdout

var print = func(format string, a ...interface{}) { fmt.Fprintf(writer, format, a...) }
var printRed = color.New(color.FgRed).Println
var printYellow = color.New(color.FgYellow).Println
var printBlue = color.New(color.FgBlue).Println
var printGreen = color.New(color.FgGreen).Println
