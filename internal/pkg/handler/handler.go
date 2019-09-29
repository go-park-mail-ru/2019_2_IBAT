package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"

	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const staticDir = "/tmp/img"
const maxUploadSize = 2 * 1024 * 1024 // 2 mb

type Handler struct {
	AuthService auth.AuthService
	UserService users.UserService
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, class, err := h.AuthService.CreateSession(r.Body, h.UserService.Storage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, &cookie)
	classJSON, _ := json.Marshal(Class{class})

	w.Write([]byte(classJSON))
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, ok := h.AuthService.Storage.Get(cookie.Value)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{"Not authorized"})
		w.Write([]byte(errJSON))
		return
	}

	classJSON, _ := json.Marshal(Class{authInfo.Class})

	w.Write([]byte(classJSON))
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.CookieName)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthService.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{"Fail to delete session"})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateSeeker(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, cookieValue := h.AuthService.Storage.Set(uuid, SeekerStr) //possible return authInfo

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	} //impossible error

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}
	http.SetCookie(w, &cookie)
	classJSON, _ := json.Marshal(Class{authInfo.Class})

	w.Write([]byte(classJSON))
}

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, cookieValue := h.AuthService.Storage.Set(uuid, EmployerStr) //possible return authInfo

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	} //impossible error

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	http.SetCookie(w, &cookie)
	classJSON, _ := json.Marshal(Class{authInfo.Class})

	w.Write([]byte(classJSON))
}

func (h *Handler) CreateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateResume(r.Body, cookie.Value, h.AuthService.Storage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) DeleteResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	strId := mux.Vars(r)["id"]
	resId, err := uuid.Parse(strId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteResume(resId, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resume, err := h.UserService.GetResume(resId, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resumeJSON, err := json.Marshal(resume)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(resumeJSON))
}

func (h *Handler) PutResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	defer r.Body.Close()
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserService.PutResume(resId, r.Body, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetSeeker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker, err := h.UserService.GetSeeker(cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{"No correct session cookie detected"})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteUser(cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthService.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) GetEmployer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer, err := h.UserService.GetEmployer(cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	defer r.Body.Close()

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, ok := h.AuthService.Storage.Get(cookie.Value) //impossible error, should use only Set method
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if authInfo.Class == SeekerStr {
		err = h.UserService.PutSeeker(r.Body, authInfo.ID)
	} else if authInfo.Class == EmployerStr {
		err = h.UserService.PutEmployer(r.Body, authInfo.ID)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

//should test method
func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker, _ := h.UserService.Storage.GetSeeker(seekId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker.Password = "" //danger
	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	emplId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer, _ := h.UserService.Storage.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer.Password = "" //danger
	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	employers := h.UserService.Storage.GetEmployers()

	for i, item := range employers {
		item.Password = ""
		employers[i] = item
	}

	employerJSON, err := json.Marshal(employers)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetResumes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resumes := h.UserService.Storage.GetResumes()

	resumesJSON, _ := json.Marshal(resumes)

	w.Write([]byte(resumesJSON))

}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vacancies := h.UserService.Storage.GetVacancies()

	vacanciesJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(vacanciesJSON))

}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateVacancy(r.Body, cookie.Value, h.AuthService.Storage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacancy, err := h.UserService.GetVacancy(vacId, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacancyJSON, err := json.Marshal(vacancy)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(vacancyJSON))
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteVacancy(vacId, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserService.PutVacancy(vacId, r.Body, cookie.Value, h.AuthService.Storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) UploadFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}

		authInfo, ok := h.AuthService.Storage.Get(cookie.Value)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			errJSON, _ := json.Marshal(Error{"Not authorized"})
			w.Write([]byte(errJSON))
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			// renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}

		// parse and validate file and post parameters
		fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("my_file")
		if err != nil {
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
			// case "application/pdf":
			break
		default:
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}
		fileName := uuid.New().String()
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}
		newPath := filepath.Join(staticDir, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			errJSON, _ := json.Marshal(Error{err.Error()})
			w.Write([]byte(errJSON))
			return
		}

		h.UserService.Storage.SetImage(authInfo.ID, authInfo.Class, newPath)
	})
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
