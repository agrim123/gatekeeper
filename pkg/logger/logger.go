package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	blue   = color.New(color.FgBlue)
	cyan   = color.New(color.FgCyan)
	green  = color.New(color.FgGreen)
)

type log struct {
	root bool
}

func L() log {
	return log{}
}

func (l log) P(privileged bool) log {
	l.root = privileged
	return l
}

func (l log) Infof(message string, attributes ...interface{}) {
	if l.root {
		blue.PrintFunc()("[INFO]  ⬆  | ")
	} else {
		blue.PrintFunc()("[INFO]     | ")
	}
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func Info(message string) {
	blue.PrintFunc()("[INFO]     | ")
	fmt.Println(message)
}

func Error(message string) {
	red.PrintFunc()("[ERROR]    | ")
	fmt.Println(message)
}

func Warn(message string) {
	yellow.PrintFunc()("[WARNING]  | ")
	fmt.Println(message)
}

func Fatal(message string) {
	red.PrintFunc()("[FATAL]    | ")
	fmt.Println(message)
	os.Exit(1)
}

func Success(message string) {
	green.PrintFunc()("[SUCCESS]  | ")
	fmt.Println(message)
}

func Warnf(message string, attributes ...interface{}) {
	yellow.PrintFunc()("[WARNING]  | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func Errorf(message string, attributes ...interface{}) {
	red.PrintFunc()("[ERROR]    | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func Fatalf(message string, attributes ...interface{}) {
	red.PrintFunc()("[FATAL]    | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
	os.Exit(1)
}

func Successf(message string, attributes ...interface{}) {
	green.PrintFunc()("[SUCCESS]  | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func Notifier(message string) {
	cyan.PrintFunc()("[NOTIFIER] | ")
	fmt.Println(message)
}

func Infof(message string, attributes ...interface{}) {
	blue.PrintFunc()("[INFO]     | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func InfofP(message string, attributes ...interface{}) {
	blue.PrintFunc()("[INFO]  ⬆  | ")
	fmt.Println(fmt.Sprintf(message, attributes...))
}

func InfoScan(message string) string {
	blue.PrintFunc()("[INFO]     | ")
	fmt.Print(message)
	var input string
	fmt.Scanln(&input)
	return input
}
