package middleware

import (
	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	. "2019_2_IBAT/pkg/pkg/interfaces"

	"log"
	"net/http"

	"github.com/google/uuid"
)

func AuthMiddlewareGenerator(authServ session.ServiceClient) (mw func(http.Handler) http.Handler) {

	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			log.Println("AuthMiddleware: started")
			cookie, err := req.Cookie(auth.CookieName)

			ctx := req.Context()
			if err != nil {
				log.Println("AuthMiddleware: No cookie detected")
			} else {
				sessionInfo, err := authServ.GetSession(ctx, &session.Cookie{
					Cookie: cookie.Value,
				})
				if err != nil {
					log.Printf("sessionInfo fetch error: %s\n", err)
				} else {
					if sessionInfo.Ok {
						ctx = NewContext(req.Context(), AuthStorageValue{
							ID:      uuid.MustParse(sessionInfo.ID),
							Role:    sessionInfo.Role,
							Expires: sessionInfo.Expires,
						})
					}
				}
			}

			reqWithCxt := req.WithContext(ctx)

			h.ServeHTTP(res, reqWithCxt)
		})
	}
	return
}
