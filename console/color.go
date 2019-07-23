package console

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func ConvertToString(input interface{}) string {
	var output string
	switch val := input.(type) {
	case int:
		output = strconv.Itoa(val)
	case string:
		output = val
	case byte:
		output = string(val)
	case bool:
		if val == true {
			output = "true"
		} else {
			output = "false"
		}
	case error:
		output = val.Error()
	}
	return output
}

func BlackStr(info interface{}) string {
	return textColor(TextBlack, ConvertToString(info))
}
func Black(info ...interface{}) {
	fmt.Print(BlackStr(fmt.Sprintln(info...)))
}

func RedStr(info interface{}) string {
	return textColor(TextRed, ConvertToString(info))
}
func Red(info ...interface{}) {
	fmt.Print(RedStr(fmt.Sprintln(info...)))
}

func GreenStr(info interface{}) string {
	return textColor(TextGreen, ConvertToString(info))
}
func Green(info ...interface{}) {
	fmt.Print(GreenStr(fmt.Sprintln(info...)))
}
func YellowStr(info interface{}) string {
	return textColor(TextYellow, ConvertToString(info))
}
func Yellow(info ...interface{}) {
	fmt.Print(YellowStr(fmt.Sprintln(info...)))
}

func BlueStr(info interface{}) string {
	return textColor(TextBlue, ConvertToString(info))
}
func Blue(info ...interface{}) {
	fmt.Print(BlueStr(fmt.Sprintln(info...)))
}
func MagentaStr(info interface{}) string {
	return textColor(TextMagenta, ConvertToString(info))
}
func Magenta(info ...interface{}) {
	fmt.Print(MagentaStr(fmt.Sprintln(info...)))
}
func CyanStr(info interface{}) string {
	return textColor(TextCyan, ConvertToString(info))
}
func Cyan(info ...interface{}) {
	fmt.Print(CyanStr(fmt.Sprintln(info...)))
}
func WhiteStr(info interface{}) string {
	return textColor(TextWhite, ConvertToString(info))
}
func White(info ...interface{}) {
	fmt.Print(WhiteStr(fmt.Sprintln(info...)))
}
func textColor(color int, str string) string {
	if runtime.GOOS == "windows" {
		return str
	}
	switch color {
	case TextBlack:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlack, str)
	case TextRed:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextRed, str)
	case TextGreen:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextGreen, str)
	case TextYellow:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextYellow, str)
	case TextBlue:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlue, str)
	case TextMagenta:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextMagenta, str)
	case TextCyan:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextCyan, str)
	case TextWhite:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextWhite, str)
	default:
		return str
	}
}
