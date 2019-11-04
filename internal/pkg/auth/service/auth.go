package service

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"

	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	gcontext "github.com/gorilla/context"
	"github.com/pkg/errors"
)

type AuthService struct {
	Storage auth.Repository
}

const CookieName = "session-id"

func (h *AuthService) CreateSession(id uuid.UUID, class string) (AuthStorageValue, string, error) {

	authInfo, cookieValue, err := h.Storage.Set(id, class)

	if err != nil {
		log.Printf("Error while unmarshaling: %s\n", err)
		err = errors.Wrap(err, "error while unmarshaling")
		return authInfo, cookieValue, errors.New("Creating session error")
	}

	return authInfo, cookieValue, nil
}

func (h *AuthService) DeleteSession(cookie *http.Cookie) bool {
	_, ok := h.Storage.Get(cookie.Value)
	if !ok {
		log.Printf("No such session")
		return false
	}

	ok = h.Storage.Delete(cookie.Value)
	if !ok {
		return false
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)

	return true
}

func (auth *AuthService) AuthMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		log.Println("AuthMiddleware: started")
		cookie, err := req.Cookie(CookieName)

		if err != nil {
			log.Println("AuthMiddleware: No cookie detected")
		} else {
			log.Println(auth.Storage)
			log.Println("AuthMiddleware: Get cookie start")
			record, ok := auth.Storage.Get(cookie.Value)

			log.Println("AuthMiddleware: Get cookie end")

			if ok {
				gcontext.Set(req, AuthRec, record)
				log.Println("AuthMiddleware: auth_record was setted")
			}
		}

		log.Println("AuthMiddleware: passing to serve")
		h.ServeHTTP(res, req)
	}

	return mw
}

func (auth *AuthService) GetSession(cookie string) (AuthStorageValue, bool) {
	return auth.Storage.Get(cookie)
}

func (auth *AuthService) SetRecord(id uuid.UUID, class string) (AuthStorageValue, string, error) {
	return auth.Storage.Set(id, SeekerStr)
}
