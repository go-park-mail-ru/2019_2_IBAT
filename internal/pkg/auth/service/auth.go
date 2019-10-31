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
	// UsSt    *usRep.UserStorage
}

const CookieName = "session-id"

func (h AuthService) CreateSession(id uuid.UUID, class string) (http.Cookie, string, error) {
	// id, class, ok := usS.CheckUser(userAuthInput.Email, userAuthInput.Password)
	// if !ok {
	// 	// log.Printf("No such user error")
	// 	return http.Cookie{}, "", errors.New("Invalid password or email")
	// }

	authInfo, cookieValue, err := h.Storage.Set(id, class) //possible return authInfo

	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "error while unmarshaling")
		return http.Cookie{}, "", errors.New("Creating session error")
	}

	expiresAt, _ := time.Parse(TimeFormat, authInfo.Expires)

	cookie := http.Cookie{
		Name:    CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	return cookie, authInfo.Role, nil
}

func (h AuthService) DeleteSession(cookie *http.Cookie) bool {
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

func (auth AuthService) AuthMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie(CookieName)

		record, ok := auth.Storage.Get(cookie.Value)
		if ok {
			gcontext.Set(req, AuthRec, record)
		}

		h.ServeHTTP(res, req)
	}

	return mw
}

func (auth AuthService) GetSession(cookie string) (AuthStorageValue, bool) {
	return auth.Storage.Get(cookie)
}

func (auth AuthService) SetRecord(id uuid.UUID, class string) (AuthStorageValue, string, error) {
	return auth.Storage.Set(id, SeekerStr)
}
