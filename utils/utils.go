package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zlorgoncho1/sprint/core"
)

func DictToJson(__dict map[string]string) string {
	var items []string
	for k, v := range __dict {
		items = append(items, fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	return "{" + strings.Join(items, ",") + "}"
}

func FormatStatusResponse(statusCode int, statusText string, protocol string) string {
	if protocol == "" {
		protocol = "HTTP/1.1"
	}
	return fmt.Sprintf("%s %d %s", protocol, statusCode, statusText)
}

func DictToHTTPHeadersResponse(headers map[string]string) string {
	var lines []string
	for k, v := range headers {
		lines = append(lines, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(lines, "\n")
}

func GetDefaultHeader(content string, contentType string) map[string]string {
	return map[string]string{
		"Content-Type":   contentType,
		"Content-Length": fmt.Sprintf("%d", len(content)),
		"Connection":     "close",
		"Date":           time.Now().String(),
	}
}

func FormatHTTPResponse(status, headers, content string) []byte {
	var response strings.Builder
	response.WriteString(status + "\n")
	response.WriteString(headers + "\n\n" + content)
	return []byte(response.String())
}

func HandleHTML(response *core.Response) {

}

func HandleJSON(response *core.Response) {
	response.ContentType = "application/json"
	jsonBytes, err := json.Marshal(response.Content)
	if err != nil {
		panic(err)
	}
	response.Content = string(jsonBytes)
}

func HandlePlainText(response *core.Response) {
	response.ContentType = "text/plain"
	response.Content = (response.Content).(string)
}
