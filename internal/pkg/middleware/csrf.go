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
		log.Printf("req.RequestURI = %s \n", req.RequestURI)
		log.Printf("req.Method = %s \n", req.Method)

		if (req.RequestURI == "/auth" && (req.Method == http.MethodPost || req.Method == http.MethodGet)) ||
			(req.Method == http.MethodPost && (req.RequestURI == "/seeker" || req.RequestURI == "/employer")) || req.Method == http.MethodGet {
		} else {
			token := req.Header.Get("X-CSRF-Token")

			authInfo, ok := context.Get(req, AuthRec).(AuthStorageValue)
			if !ok {
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
