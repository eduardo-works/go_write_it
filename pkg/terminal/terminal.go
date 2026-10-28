package terminal

import (
	"fmt"
	"os"
	"strings"
)

const (
	Reset = "\033[0m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
)

var colors = map[string]string{
	"Reset":         Reset,
	"Black":         Black,
	"Red":           Red,
	"Green":         Green,
	"Yellow":        Yellow,
	"Blue":          Blue,
	"Magenta":       Magenta,
	"Cyan":          Cyan,
	"White":         White,
	"BrightBlack":   BrightBlack,
	"BrightRed":     BrightRed,
	"BrightGreen":   BrightGreen,
	"BrightYellow":  BrightYellow,
	"BrightBlue":    BrightBlue,
	"BrightMagenta": BrightMagenta,
	"BrightCyan":    BrightCyan,
	"BrightWhite":   BrightWhite,
	"Bold":          Bold,
	"Dim":           Dim,
	"Italic":        Italic,
	"Underline":     Underline,
}

func Format(message string) string {
	return replace(message + "{Reset}")
}

func FormatLine(message string) string {
	return Format(message) + "\n"
}

func Write(message string) {
	fmt.Print(Format(message))
}

func WriteLine(message string) {
	fmt.Print(FormatLine(message))
}

func Error(message string) {
	fmt.Fprint(os.Stderr, FormatLine("{BrightRed}{Italic}"+message))
}

func Warning(message string) {
	fmt.Fprint(os.Stderr, FormatLine("{BrightYellow}{Italic}"+message))
}

func Info(message string) {
	WriteLine("{BrightGreen}{Italic}" + message + "{Reset}")
}

func replace(message string) string {
	for name, code := range colors {
		message = strings.ReplaceAll(message, "{"+name+"}", code)
	}

	return message
}
