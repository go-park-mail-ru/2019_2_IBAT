package middleware

import (
	"2019_2_IBAT/internal/pkg/auth"
	"net/http"
)

type key string

const authRecord key = "AuthRecord"

func AuthMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie(auth.CookieName)

		if cookie.Value == "" {
			h.ServeHTTP(res, req)
		} else {

		}

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

		// h.ServeHTTP(res, req)
	}

	return mw
}

// type AuthService struct {
// 	Storage AuthStorage
// }

// const CookieName = "session-id"

// func (h *AuthService) CreateSession(body io.ReadCloser, usS UserStorage) (http.Cookie, string, error) {
// 	bytes, err := ioutil.ReadAll(body)
// 	if err != nil {
// 		// log.Printf("error while reading body: %s", err)
// 		// err = errors.Wrap(err, "reading body error")
// 		return http.Cookie{}, "", errors.New("Invalid body, transfer error")
// 	}

// 	userAuthInput := new(UserAuthInput)
// 	err = json.Unmarshal(bytes, userAuthInput)
// 	if err != nil {
// 		// log.Printf("Error while unmarshaling: %s", err)
// 		// err = errors.Wrap(err, "error while unmarshaling")
// 		return http.Cookie{}, "", errors.New("Invalid json")
// 	}

// 	id, class, ok := usS.CheckUser(userAuthInput.Email, userAuthInput.Password)
// 	if !ok {
// 		// log.Printf("No such user error")
// 		return http.Cookie{}, "", errors.New("Invalid password or email")
// 	}

// 	authInfo, cookieValue, err := h.Storage.Set(id, class) //possible return authInfo

// 	if err != nil {
// 		// log.Printf("Error while unmarshaling: %s", err)
// 		// err = errors.Wrap(err, "error while unmarshaling")
// 		return http.Cookie{}, "", errors.New("Creating session error")
// 	}

// 	expiresAt, _ := time.Parse(TimeFormat, authInfo.Expires)

// 	cookie := http.Cookie{
// 		Name:    CookieName,
// 		Value:   cookieValue,
// 		Expires: expiresAt,
// 	}

// 	return cookie, authInfo.Role, nil
// }

// func (h *AuthService) DeleteSession(cookie *http.Cookie) bool {
// 	_, ok := h.Storage.Get(cookie.Value)
// 	if !ok {
// 		log.Printf("No such session")
// 		return false
// 	}

// 	ok = h.Storage.Delete(cookie.Value)
// 	if !ok {
// 		return false
// 	}
// 	cookie.Expires = time.Now().AddDate(0, 0, -1)

// 	return true
// }
