package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	corsData = CorsData{
		AllowOrigins: []string{
			// "localhost:8080",
			// "20192ibat-cyb91y0rs.now.sh",
			"*",
		},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type", "X-Content-Type-Options", "X-Csrf-Token"},
		AllowCredentials: true,
	}
)

// https://20192ibat-cyb91y0rs.now.sh/
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
			// var flag bool
			// for _, header := range corsData.AllowOrigins {
			// 	if val[0] == header {
			// 		flag = true
			// 	}
			// }
			flag := true
			if flag {
				res.Header().Set("Access-Control-Allow-Origin", val[0])
				res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
				log.Println("Access-Control-Allow-Origin and Access-Control-Allow-Credentials headers were set")
			} else {
				log.Println("StatusForbidden: headers not set")
				res.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if req.Method == "OPTIONS" {
			log.Println("Access-Control-Allow-Methods and Access-Control-Allow-Headers headers were set")
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			return
		}

		h.ServeHTTP(res, req)
	}

	return mw
}
