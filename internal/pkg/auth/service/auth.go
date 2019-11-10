package service

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"context"

	"log"
	"net/http"

	"github.com/google/uuid"
	// context "github.com/gorilla/context"

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

func (h *AuthService) DeleteSession(cookie string) bool {
	_, ok := h.Storage.Get(cookie)
	if !ok {
		log.Printf("No such session")
		return false
	}

	ok = h.Storage.Delete(cookie)
	if !ok {
		return false
	}

	return true
}

func (auth *AuthService) AuthMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		log.Println("AuthMiddleware: started")
		cookie, err := req.Cookie(CookieName)

		ctx := context.TODO()
		if err != nil {
			log.Println("AuthMiddleware: No cookie detected")
		} else {
			log.Println(auth.Storage)
			record, ok := auth.Storage.Get(cookie.Value)

			if ok {
				// context.Set(req, AuthRec, record)
				ctx = NewContext(req.Context(), record)
				log.Println("AuthMiddleware: auth_record was setted")
			} else {
				log.Println("AuthMiddleware: failed to set auth_record")
			}
		}

		reqWithCxt := req.WithContext(ctx)
		log.Println("CTX")
		log.Println(reqWithCxt)
		log.Println("AuthMiddleware: passing to serve")

		h.ServeHTTP(res, reqWithCxt)
	}

	return mw
}

func (auth *AuthService) GetSession(cookie string) (AuthStorageValue, bool) {
	return auth.Storage.Get(cookie)
}

func (auth *AuthService) SetRecord(id uuid.UUID, class string) (AuthStorageValue, string, error) {
	return auth.Storage.Set(id, SeekerStr)
}
