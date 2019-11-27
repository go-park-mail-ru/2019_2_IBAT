package handler

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/app/users"
	. "2019_2_IBAT/pkg/pkg/models"
)

const publicDir = "/static"

const maxUploadSize = 2 * 1024 * 1024 // 2 mb

type Handler struct {
	InternalDir string
	AuthService session.ServiceClient
	UserService users.Service
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		log.Println("GetUser Handler: unauthorized")
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	if authInfo.Role == SeekerStr {
		seeker, err := h.UserService.GetSeeker(authInfo.ID)

		if err != nil {
			log.Println("GetUser Handler: failed to get seeker")
			SetError(w, http.StatusBadRequest, InternalErrorMsg)
			return
		}

		answer := UserSeekAnswer{
			Role:   SeekerStr,
			Seeker: seeker,
		}
		answerJSON, _ := answer.MarshalJSON()

		w.Write(answerJSON)
	} else if authInfo.Role == EmployerStr {
		employer, err := h.UserService.GetEmployer(authInfo.ID)

		if err != nil {
			log.Println("GetUser Handler: failed to get employer")
			SetError(w, http.StatusBadRequest, InternalErrorMsg)
			return
		}

		answer := UserEmplAnswer{
			Role:     EmployerStr,
			Employer: employer,
		}
		answerJSON, _ := answer.MarshalJSON()

		w.Write([]byte(answerJSON))
	} else {
		log.Println("GetUser Handler: unauthorized")
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	err := h.UserService.DeleteUser(authInfo)

	if err != nil {
		SetError(w, http.StatusForbidden, ForbiddenMsg)
		return
	}

	cookie, _ := r.Cookie(auth.CookieName) //костыль

	sessionBool, err := h.AuthService.DeleteSession(r.Context(), &session.Cookie{
		Cookie: cookie.Value,
	})
	if !sessionBool.Ok {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	defer r.Body.Close()

	authInfo, ok := FromContext(r.Context())

	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	if authInfo.Role == SeekerStr {
		err := h.UserService.PutSeeker(r.Body, authInfo.ID)
		if err != nil {
			SetError(w, http.StatusForbidden, ForbiddenMsg)
			return
		}
	} else if authInfo.Role == EmployerStr {
		err := h.UserService.PutEmployer(r.Body, authInfo.ID)
		if err != nil {
			SetError(w, http.StatusForbidden, ForbiddenMsg)
			return
		}
	}
}

func (h *Handler) UploadFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		authInfo, ok := FromContext(r.Context())

		if !ok {
			SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			SetError(w, http.StatusBadRequest, "Invalid size")
			return
		}

		file, _, err := r.FormFile("my_file")
		if err != nil {
			SetError(w, http.StatusBadRequest, "Invalid form key")
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			SetError(w, http.StatusBadRequest, "Bad file")
			return
		}

		filetype := http.DetectContentType(fileBytes)

		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
			break
		default:
			SetError(w, http.StatusBadRequest, "Invalid extension")
			return
		}

		fileName := uuid.New().String()
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			SetError(w, http.StatusBadRequest, "Invalid extension")
			return
		}

		pkgPath := filepath.Join(h.InternalDir, fileName+fileEndings[0])

		newFile, err := os.Create(pkgPath)
		if err != nil {
			SetError(w, http.StatusInternalServerError, "Failed to set image")
			return
		}

		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			SetError(w, http.StatusInternalServerError, InternalErrorMsg)
			return
		}

		publicPath := filepath.Join(publicDir, fileName+fileEndings[0])
		h.UserService.SetImage(authInfo.ID, authInfo.Role, publicPath)
	})
}
