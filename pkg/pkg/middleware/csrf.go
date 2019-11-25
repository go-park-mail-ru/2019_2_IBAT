package middleware

import (
	auth "2019_2_IBAT/pkg/app/auth"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"log"

	"encoding/json"
	"net/http"
)

func CSRFMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		log.Printf("req.RequestURI = %s \n", req.RequestURI)
		log.Printf("req.Method = %s \n", req.Method)

		if (req.RequestURI == "/auth" && (req.Method == http.MethodPost || req.Method == http.MethodGet)) ||
			(req.Method == http.MethodPost && (req.RequestURI == "/seeker" ||
				req.RequestURI == "/employer")) || req.Method == http.MethodGet {
		} else {
			token := req.Header.Get("X-Csrf-Token")

			authInfo, ok := FromContext(req.Context())
			if !ok {
				res.Header().Set("Content-Type", "application/json; charset=UTF-8")
				res.WriteHeader(http.StatusUnauthorized)
				errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg}) //token msg
				res.Write([]byte(errJSON))
				return
			}

			cookie, _ := req.Cookie(auth.CookieName)

			_, err := csrf.Tokens.Check(authInfo.ID.String(), cookie.Value, token)
			if err != nil {
				res.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
