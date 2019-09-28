package handler

import (
	"encoding/json"
	"hh_workspace/2019_2_IBAT/internal/pkg/auth"
	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"
	"hh_workspace/2019_2_IBAT/internal/pkg/users"

	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	AuthHandler   auth.Handler
	UserControler users.Controler
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)

	cookie, err := h.AuthHandler.CreateSession(r.Body, h.UserControler.Storage)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, &cookie)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthHandler.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)
	uuid, err := h.UserControler.HandleCreateSeeker(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, cookieValue := h.AuthHandler.Storage.Set(uuid, SeekerStr) //possible return authInfo

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
}

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)

	uuid, err := h.UserControler.HandleCreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, cookieValue := h.AuthHandler.Storage.Set(uuid, EmployerStr) //possible return authInfo

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
}

func (h *Handler) CreateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserControler.HandleCreateResume(r.Body, cookie.Value, h.AuthHandler.Storage)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Message{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) DeleteResume(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	strId := mux.Vars(r)["id"]
	resId, err := uuid.Parse(strId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserControler.HandleDeleteResume(resId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resume, err := h.UserControler.HandleGetResume(resId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resumeJSON, err := json.Marshal(resume)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(resumeJSON))
}

func (h *Handler) PutResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserControler.HandlePutResume(resId, r.Body, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetSeeker(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker, err := h.UserControler.HandleGetSeeker(cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{"No correct session cookie detected"})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserControler.HandleDeleteUser(cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthHandler.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) GetEmployer(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer, err := h.UserControler.HandleGetEmployer(cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, ok := h.AuthHandler.Storage.Get(cookie.Value) //impossible error, should use only Set method
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if authInfo.Class == SeekerStr {
		err = h.UserControler.HandlePutSeeker(r.Body, authInfo.ID)
	} else if authInfo.Class == EmployerStr {
		err = h.UserControler.HandlePutEmployer(r.Body, authInfo.ID)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

//should test method
func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker, _ := h.UserControler.Storage.GetSeeker(seekId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	seeker.Password = "" //danger
	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	emplId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer, _ := h.UserControler.Storage.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	employer.Password = "" //danger
	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	employers := h.UserControler.Storage.GetEmployers()

	for i, item := range employers {
		item.Password = ""
		employers[i] = item
	}

	employerJSON, err := json.Marshal(employers)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetResumes(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	resumes := h.UserControler.Storage.GetResumes()

	resumesJSON, _ := json.Marshal(resumes)

	w.Write([]byte(resumesJSON))

}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	vacancies := h.UserControler.Storage.GetVacancies()

	vacanciesJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(vacanciesJSON))

}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserControler.HandleCreateVacancy(r.Body, cookie.Value, h.AuthHandler.Storage)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Message{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacancy, err := h.UserControler.HandleGetVacancy(vacId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacancyJSON, err := json.Marshal(vacancy)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(vacancyJSON))
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserControler.HandleDeleteVacancy(vacId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setDefaultHeaders(w)
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserControler.HandlePutVacancy(vacId, r.Body, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}
