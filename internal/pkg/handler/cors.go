package handler

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	corsData = CorsData{
		AllowOrigins: []string{
			"localhost:8080",
		},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}
)

type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func CorsMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		val, ok := req.Header["Origin"]
		if ok {
			res.Header().Set("Access-Control-Allow-Origin", val[0])
			res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
		}

		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			return
		}

		h.ServeHTTP(res, req)
	}

	return mw
}
