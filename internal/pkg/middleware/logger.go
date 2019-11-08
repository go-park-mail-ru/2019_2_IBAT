package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type AccessLogger struct {
	StdLogger *log.Logger
}

func (ac AccessLogger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		fmt.Printf("FMT [%s] %s, %s %s\n\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		log.Printf("LOG [%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		ac.StdLogger.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}
