package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/fatih/color"
)

var methodColor = color.New(color.FgBlue).SprintFunc()
var durationColor = color.New(color.FgRed).SprintFunc()

func LogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		for k, v := range rec.HeaderMap {
			for _, x := range v {
				w.Header().Add(k, x)
			}
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())

		startFormatted := start.Format("2006-01-02 15:04:05.9999")
		duration := durationColor(fmt.Sprintf("%.2fms", float64(time.Now().Sub(start))/float64(time.Millisecond)))
		method := methodColor(fmt.Sprintf("%s", r.Method))
		url := r.URL.String()
		respCode := rec.Code
		length := rec.Body.Len()

		log.Println(fmt.Sprintf("[%v] method: %v url: %v status: %v duration: %v length: %v", startFormatted, method, url, respCode, duration, length))
	})
}
