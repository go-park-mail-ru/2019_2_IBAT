package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	corsData = CorsData{
		AllowOrigins: []string{
			// "localhost:8080",
			"20192ibat-cyb91y0rs.now.sh",
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

// func CorsMiddleware(h http.Handler) http.Handler {
// 	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
// 		// fmt.Println(req.Context())
// 		// fmt.Println("Request was accepted")

// 		val, _ := req.Header["Origin"]
// 		// if ok {
// 		// 	res.Header().Set("Access-Control-Allow-Origin", val[0])
// 		// 	res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
// 		// }
// 		// var flag bool
// 		// for _, header := range corsData.AllowOrigins {
// 		// 	if val[0] == header {
// 		// 		flag = true
// 		// 	}
// 		// }

// 		// // val, ok := req.Header["Origin"]
// 		// // if ok {
// 		// // res.Header().Set("Access-Control-Allow-Origin", strings.Join(corsData.AllowOrigins, ", "))
// 		// // res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
// 		// // }
// 		// if !flag {
// 		// 	res.WriteHeader(http.StatusBadRequest)
// 		// 	return
// 		// }
// 		if req.Method == "OPTIONS" {
// 			// res.Header().Set("Access-Control-Allow-Origin", strings.Join(corsData.AllowOrigins, ", "))
// 			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
// 			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
// 			// res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
// 			return
//		}

// 		h.ServeHTTP(res, req)
// 	}

// 	return mw
// }

func CorsMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Context())
		fmt.Println("Request was accepted")
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
			} else {
				res.WriteHeader(http.StatusForbidden)
				return
			}

		}
		// else {
		// 	res.WriteHeader(http.StatusForbidden)
		// 	return
		// }

		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			return
		}

		h.ServeHTTP(res, req)
	}

	return mw
}
