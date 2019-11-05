package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type AccessLogger struct {
	StdLogger *log.Logger
	// ZapLogger    *zap.SugaredLogger
	// LogrusLogger *logrus.Entry
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

		// ac.ZapLogger.Info(r.URL.Path,
		// 	zap.String("method", r.Method),
		// 	zap.String("remote_addr", r.RemoteAddr),
		// 	zap.String("url", r.URL.Path),
		// 	zap.Duration("work_time", time.Since(start)),
		// )

		// ac.LogrusLogger.WithFields(logrus.Fields{
		// 	"method":      r.Method,
		// 	"remote_addr": r.RemoteAddr,
		// 	"work_time":   time.Since(start),
		// }).Info(r.URL.Path)
	})
}
