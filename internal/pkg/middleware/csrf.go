package middleware

import (
	auth "2019_2_IBAT/internal/pkg/auth"
	csrf "2019_2_IBAT/internal/pkg/csrf"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"log"

	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

func CSRFMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		// fmt.Println(req.Context())
		// fmt.Println("Request was accepted")
		// val, ok := req.Header["Origin"]
		// if ok {
		// 	res.Header().Set("Access-Control-Allow-Origin", val[0])
		// 	res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
		// }

		// if req.Method == "OPTIONS" {
		// 	res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
		// 	res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
		// 	return
		// }
		log.Printf("req.RequestURI = %s \n", req.RequestURI)
		log.Printf("req.Method = %s \n", req.Method)

		if (req.RequestURI == "/auth" && (req.Method == http.MethodPost || req.Method == http.MethodGet)) ||
			(req.Method == http.MethodPost && (req.RequestURI == "/seeker" || req.RequestURI == "/employer")) || req.Method == http.MethodGet {
		} else {
			token := req.Header.Get("X-CSRF-Token")

			// token, ok := req.Header["X-CSRF-Token"]

			authInfo, ok := context.Get(req, AuthRec).(AuthStorageValue)
			if !ok {
				// res.Header()
				res.WriteHeader(http.StatusUnauthorized)
				errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg}) //token msg
				res.Write([]byte(errJSON))
				return
			}

			cookie, _ := req.Cookie(auth.CookieName)

			_, err := csrf.Tokens.Check(authInfo.ID.String(), cookie.Value, token)
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg}) //token msg
				res.Write([]byte(errJSON))
				return
			}
		}
		h.ServeHTTP(res, req)
	}

	return mw
}
