package logging

import (
	"fmt"
)

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	italic    = "\033[3m"
	underline = "\033[4m"
)

func printPrefix() string {
	color := string(config.Color)
	return color + bold + "[" + string(config.Role) + "]" + reset + color
}

func Println(a ...any) {
	println(printPrefix(), fmt.Sprint(a...), reset)
}

func Printf(format string, a ...any) {
	println(printPrefix() + fmt.Sprintf(format, a...) + reset)
}

func PrintErrStr(a ...any) {
	println(printPrefix() + " Error: " + bold + fmt.Sprint(a...) + reset)
}

func PrintErr(err error) {
	println(printPrefix() + " Error: " + bold + err.Error() + reset)
}
