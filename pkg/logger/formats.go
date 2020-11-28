package logger

import "github.com/fatih/color"

func Underline(message string) string {
	return color.New(color.Underline).Sprint(message)
}

func Bold(message string) string {
	return color.New(color.Bold).Sprint(message)
}

func Italic(message string) string {
	return color.New(color.Italic).Sprint(message)
}
