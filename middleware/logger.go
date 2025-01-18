package middleware

import (
	"fmt"
	"github.com/gorilla/handlers"
	"io"
	"time"
)

func Logger(writer io.Writer, params handlers.LogFormatterParams) {
	clientIP := params.Request.RemoteAddr
	method := params.Request.Method
	statusCode := params.StatusCode
	latency := time.Since(params.TimeStamp)
	path := params.URL.Path

	logLine := fmt.Sprintf("[API] %v | %3d | %13v | %15s | %-7s %s\n",
		params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusCode,
		latency,
		clientIP,
		method,
		path,
	)
	writer.Write([]byte(logLine))
}
