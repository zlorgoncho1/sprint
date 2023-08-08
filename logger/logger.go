package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Logger struct{}

var colorMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
}

func (l Logger) Print(message interface{}, moduleName string) {
	fmt.Println(message)
}

func (l Logger) Color(__color string) func(a ...interface{}) string {
	colorValue, ok := colorMap[strings.ToLower(__color)]
	if !ok {
		colorValue = color.FgWhite
	}
	return color.New(colorValue).SprintFunc()
}

func (l Logger) Debug(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("white")("[Sprint] [Beta - v0.0.0]"),
			l.Color("magenta")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("white")("DEBUG"),
			l.Color("magenta")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("white")(message)),
	)
}

func (l Logger) Log(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("blue")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")("LOG"),
			l.Color("green")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("blue")(message)),
	)
}

func (l Logger) Warn(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("yellow")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("yellow")("WARN"),
			l.Color("white")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("yellow")(message)),
	)
}

func (l Logger) Error(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("red")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("red")("ERROR"),
			l.Color("white")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("red")(message)),
	)
}

func (l Logger) Plog(message interface{}, elapsed time.Duration, moduleName string, statusCode string, statusMessage string) {
	var formattedMessage string
	switch statusCode {
	case "0":
		formattedMessage = fmt.Sprintf("%s %s %s %s %s %s",
			l.Color("blue")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")("LOG"),
			l.Color("green")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("blue")(message),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	case "2":
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("green")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	case "3":
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	default:
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Beta - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("red")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	}
	fmt.Println(formattedMessage)
}

func (l Logger) Reload() {
	fmt.Printf(l.Color("white")(fmt.Sprintf("[Sprint] [Beta - v0.0.0] - [%s] - Server Reloading ...\n", time.Now().Format("02/01/2006, 15:04:05"))))
}
