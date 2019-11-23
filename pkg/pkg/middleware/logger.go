package middleware

import (
	// "fmt"
	"log"
	"net/http"
	"time"
	"os"
)

type Logger struct {
	StdLogger *log.Logger
	f *os.File
}

func (ac Logger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("LOG START [%s] %s, %s \n",
		r.Method, r.RemoteAddr, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("LOG END [%s] %s, %s %s\n\n\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

func NewLogger() Logger {
	loger := Logger{}
	
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	loger.f = f


	log.SetOutput(loger.f)
	return loger
}