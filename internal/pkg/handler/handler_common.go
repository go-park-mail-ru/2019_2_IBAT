package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/context"
)

const publicDir = "/static"

const maxUploadSize = 2 * 1024 * 1024 // 2 mb

type Handler struct {
	InternalDir string
	AuthService auth.Service
	UserService users.Service
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)

	if !ok {
		log.Println("GetUser Handler: unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if authInfo.Role == SeekerStr {
		seeker, err := h.UserService.GetSeeker(authInfo.ID)

		if err != nil {
			log.Println("GetUser Handler: failed to get seeker")
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
			w.Write([]byte(errJSON))
			return
		}

		answer := UserSeekAnswer{
			Role:   SeekerStr,
			Seeker: seeker,
		}
		answerJSON, _ := json.Marshal(answer)

		w.Write([]byte(answerJSON))
	} else if authInfo.Role == EmployerStr {
		employer, err := h.UserService.GetEmployer(authInfo.ID)

		if err != nil {
			log.Println("GetUser Handler: failed to get employer")
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
			w.Write([]byte(errJSON))
			return
		}

		answer := UserEmplAnswer{
			Role:     EmployerStr,
			Employer: employer,
		}
		answerJSON, _ := json.Marshal(answer)

		w.Write([]byte(answerJSON))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("GetUser Handler: unauthorized")
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	err := h.UserService.DeleteUser(authInfo)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	cookie, _ := r.Cookie(auth.CookieName) //костыль

	ok = h.AuthService.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	defer r.Body.Close()

	authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if authInfo.Role == SeekerStr {
		err := h.UserService.PutSeeker(r.Body, authInfo.ID)
		if err != nil {
			w.WriteHeader(http.StatusForbidden) //should add invalid email case
			errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
			w.Write([]byte(errJSON))
			return
		}
	} else if authInfo.Role == EmployerStr {
		err := h.UserService.PutEmployer(r.Body, authInfo.ID)
		if err != nil {
			w.WriteHeader(http.StatusForbidden) //should add invalid email case
			errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
			w.Write([]byte(errJSON))
			return
		}
	}
}

func (h *Handler) UploadFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
			w.Write([]byte(errJSON))
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: "Invalid size"})
			w.Write([]byte(errJSON))
			return
		}

		// parse and validate file and post parameters
		file, _, err := r.FormFile("my_file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: "Invalid form key"})
			w.Write([]byte(errJSON))
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: "Bad file"})
			w.Write([]byte(errJSON))
			return
		}

		filetype := http.DetectContentType(fileBytes)

		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: "Invalid extension"})
			w.Write([]byte(errJSON))
			return
		}

		fileName := uuid.New().String()
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{Message: "Invalid extension"})
			w.Write([]byte(errJSON))
			return
		}

		internalPath := filepath.Join(h.InternalDir, fileName+fileEndings[0])

		//fmt.Println(internalPath)
		newFile, err := os.Create(internalPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errJSON, _ := json.Marshal(Error{Message: "Failed to set image"})
			w.Write([]byte(errJSON))
			return
		}

		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
			w.Write([]byte(errJSON))
			return
		}

		publicPath := filepath.Join(publicDir, fileName+fileEndings[0])
		h.UserService.SetImage(authInfo.ID, authInfo.Role, publicPath)
	})
}
