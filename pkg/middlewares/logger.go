package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func LogRequestResponse(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// -Read the request body (if any)
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
			//  Reset the request body so it can be read again by handler
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		//  Log request
		log.Printf("Received %s request for %s at %s", req.Method, req.URL.Path, start.Format(time.RFC3339))
		if len(bodyBytes) > 0 {
			log.Printf("Request body : %s", string(bodyBytes))
		}
		// Response recoder to capture the response body
		rec := &responseRecorder{ResponseWriter: res}
		next.ServeHTTP(rec, req)
		duration := time.Since(start)
		log.Printf("Responded with status %d in %v", rec.statusCode, duration)
		log.Printf("Response body: %s", rec.body.String())

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
