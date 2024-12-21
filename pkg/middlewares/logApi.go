package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func LogRequestResponse(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// Read request body
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Response recorder to capture the response
		rec := &responseRecorder{ResponseWriter: res}

		// Process the request
		next.ServeHTTP(rec, req)

		// Calculate duration
		duration := time.Since(start)

		// Get IP address
		ip := req.RemoteAddr
		if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		// Format the log message
		logMessage := fmt.Sprintf("%s %s | %s | %s | %s | IP - %s | PATH - %s | STATUS CODE - %d | TIME TAKEN - %v",
			time.Now().Format("02"),       // day
			time.Now().Format("15:04:05"), // time
			"INFO",                        // level
			"Request processed",           // message
			"Internal",                    // service name
			ip,
			req.URL.Path,
			rec.statusCode,
			duration,
		)

		log.Println(logMessage)

		if len(bodyBytes) > 0 {
			log.Printf("Request body: %s", string(bodyBytes))
		}
		if rec.body.Len() > 0 {
			log.Printf("Response body: %s", rec.body.String())
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *responseRecorder) Write(p []byte) (n int, err error) {
	rec.body.Write(p)
	return rec.ResponseWriter.Write(p)
}
