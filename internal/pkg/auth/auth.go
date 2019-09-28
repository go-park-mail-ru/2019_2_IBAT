package auth

import (
	"encoding/json"
	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type AuthService struct {
	Storage AuthStorage
}

const CookieName = "session-id"

func (h *AuthService) CreateSession(body io.ReadCloser, usS UserStorage) (http.Cookie, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return http.Cookie{}, err
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "Error while unmarshaling")
		return http.Cookie{}, err
	}

	id, class, ok := usS.CheckUser(userAuthInput.Login, userAuthInput.Password)
	if !ok {
		// log.Printf("No such user error")
		return http.Cookie{}, errors.New("No such user error")
	}

	authInfo, cookieValue := h.Storage.Set(id, class) //possible return authInfo

	// authInfo, _ := h.Storage.Get(cookieValue) //impossible error, should use only Set method

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)

	if err != nil {
		log.Printf("Error while time conversing: %s", err)
		err = errors.Wrap(err, "Error while time conversing")
		return http.Cookie{}, err
	} //impossible error

	cookie := http.Cookie{
		Name:    CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	return cookie, nil
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
