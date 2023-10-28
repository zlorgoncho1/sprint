package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Logger struct is an empty struct that serves as a receiver for logging methods.
type Logger struct{}

// colorMap maps string representations of colors to their corresponding color attributes
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

// Print simply prints a message without any additional formatting or coloring.
func (l Logger) Print(message interface{}, moduleName string) {
	fmt.Println(message)
}

// Color returns a SprintFunc that applies the specified text color.
func (l Logger) Color(__color string) func(a ...interface{}) string {
	colorValue, ok := colorMap[strings.ToLower(__color)]
	if !ok {
		colorValue = color.FgWhite // Default to white if color not found.
	}
	return color.New(colorValue).SprintFunc()
}

// Debug logs a message with the DEBUG level in a specific format, including time, level, module name, and message.
func (l Logger) Debug(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("white")("[Sprint] [Dev - v0.0.0]"),
			l.Color("magenta")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("white")("DEBUG"),
			l.Color("magenta")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("white")(message)),
	)
}

// Log logs a general message, with the LOG level.
func (l Logger) Log(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("blue")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")("LOG"),
			l.Color("green")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("blue")(message)),
	)
}

// Warn logs a message with the WARN level.
func (l Logger) Warn(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("yellow")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("yellow")("WARN"),
			l.Color("white")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("yellow")(message)),
	)
}

// Error logs a message with the ERROR level.
func (l Logger) Error(message interface{}, moduleName string) {
	fmt.Println(
		fmt.Sprintf("%s %s %s %s %s",
			l.Color("red")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("red")("ERROR"),
			l.Color("white")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("red")(message)),
	)
}

// Plog is a performance logger, logging the elapsed time along with the message, status code, and other details.
func (l Logger) Plog(message interface{}, elapsed time.Duration, moduleName string, statusCode string, statusMessage string) {
	var formattedMessage string
	switch statusCode {
	// Handling different status codes to format the message appropriately.
	case "0":
		formattedMessage = fmt.Sprintf("%s %s %s %s %s %s",
			l.Color("blue")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")("LOG"),
			l.Color("green")(fmt.Sprintf("[%s]", moduleName)),
			l.Color("blue")(message),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	case "2":
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("green")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	case "3":
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("blue")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	default:
		// Default case: log with error status code.
		formattedMessage = fmt.Sprintf("%s %s %s %s",
			l.Color("white")("[Sprint] [Dev - v0.0.0]"),
			l.Color("white")(time.Now().Format("| 02/01/2006 - 15:04:05 |")),
			l.Color("red")(fmt.Sprintf("LOG [%s] [\"%s\"] %s", moduleName, statusMessage, message)),
			l.Color("cyan")(fmt.Sprintf("+%.0f ms", elapsed.Seconds()*1000)))
	}
	fmt.Println(formattedMessage)
}

// Reload logs a server reloading message with a timestamp.
func (l Logger) Reload() {
	fmt.Printf(l.Color("white")(fmt.Sprintf("[Sprint] [Dev - v0.0.0] - [%s] - Server Reloading ...\n", time.Now().Format("02/01/2006, 15:04:05"))))
}
