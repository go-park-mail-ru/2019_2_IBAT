package handler

import (
	"encoding/json"
	"fmt"
	"hh_workspace/2019_2_IBAT/internal/pkg/auth"
	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"
	"hh_workspace/2019_2_IBAT/internal/pkg/users"

	"log"
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
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Start creating session")
	cookie, err := h.AuthHandler.CreateSession(r.Body, h.UserControler.Storage)
	if err != nil {
		log.Printf("error while creatig session: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error while creatig session"))
		return
	}

	http.SetCookie(w, &cookie)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	ok := h.AuthHandler.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No session detected"))
		return
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("Cookie deleted"))
}

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	uuid, err := h.UserControler.HandleCreateSeeker(r.Body)
	if err != nil {
		log.Printf("here: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	cookieValue := h.AuthHandler.Storage.Set(uuid, SeekerStr) //possible return authInfo

	authInfo, ok := h.AuthHandler.Storage.Get(cookieValue) //impossible error, should use only Set method
	if !ok {
		log.Printf("Error: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		log.Printf("Error while time conversing: %s", err)
		w.WriteHeader(http.StatusNotFound)
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
	w.Header().Set("Content-Type", "application/json")

	uuid, err := h.UserControler.HandleCreateEmployer(r.Body)
	if err != nil {
		log.Printf("Error %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while creating employer"))
		return
	}

	cookieValue := h.AuthHandler.Storage.Set(uuid, EmployerStr) //possible return authInfo

	authInfo, ok := h.AuthHandler.Storage.Get(cookieValue) //impossible error, should use only Set method
	if !ok {
		log.Printf("Error: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Session error"))
		return
	}

	expiresAt, err := time.Parse(auth.TimeFormat, authInfo.Expires)
	if err != nil {
		log.Printf("Error while time conversing: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while time conversing"))
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
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	id, err := h.UserControler.HandleCreateResume(r.Body, cookie.Value, h.AuthHandler.Storage)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while creating resume"))
		return
	}

	jsonString := `{ "name":` + `"` + id.String() + `"` + "}" //change

	w.Write([]byte(jsonString))
}

func (h *Handler) DeleteResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	strId := mux.Vars(r)["id"]
	resId, err := uuid.Parse(strId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	err = h.UserControler.HandleDeleteResume(resId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	resume, err := h.UserControler.HandleGetResume(resId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	resumeJSON, err := json.Marshal(resume)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(resumeJSON))
}

func (h *Handler) PutResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserControler.HandlePutResume(resId, r.Body, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while creating resume"))
		return
	}
}

func (h *Handler) GetSeeker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	seeker, err := h.UserControler.HandleGetSeeker(cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		w.Write([]byte("No session detected"))
		return
	}

	http.SetCookie(w, cookie)
	// w.Write([]byte("Cookie deleted"))
}

func (h *Handler) GetEmployer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	employer, err := h.UserControler.HandleGetEmployer(cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{"No correct session cookie detected"})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, ok := h.AuthHandler.Storage.Get(cookie.Value) //impossible error, should use only Set method
	if !ok {
		log.Printf("Error: %s", err)
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
		errJSON, _ := json.Marshal(Error{"Change failed"})
		w.Write([]byte(errJSON))
		return
	}
}

//should test method
func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	seeker, _ := h.UserControler.Storage.GetSeeker(seekId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	seeker.Password = "" //danger
	seekerJSON, err := json.Marshal(seeker)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(seekerJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	emplId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	employer, _ := h.UserControler.Storage.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	employer.Password = "" //danger
	employerJSON, err := json.Marshal(employer)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	employers := h.UserControler.Storage.GetEmployers()

	for i, item := range employers {
		item.Password = ""
		employers[i] = item
	}

	employerJSON, err := json.Marshal(employers)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errJSON, _ := json.Marshal(Error{
			Body: err.Error(),
		})
		w.Write([]byte(errJSON))
		return
	}

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetResumes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resumes := h.UserControler.Storage.GetResumes()

	resumesJSON, _ := json.Marshal(resumes)

	w.Write([]byte(resumesJSON))

}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vacancies := h.UserControler.Storage.GetVacancies()

	vacanciesJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(vacanciesJSON))

}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	id, err := h.UserControler.HandleCreateVacancy(r.Body, cookie.Value, h.AuthHandler.Storage)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while creating resume"))
		return
	}

	jsonString := `{ "name":` + `"` + id.String() + `"` + "}" //change

	w.Write([]byte(jsonString))
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	vacancy, err := h.UserControler.HandleGetVacancy(vacId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	vacancyJSON, err := json.Marshal(vacancy)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(vacancyJSON))
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid id"))
		return
	}

	err = h.UserControler.HandleDeleteVacancy(vacId, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserControler.HandlePutVacancy(vacId, r.Body, cookie.Value, h.AuthHandler.Storage)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error while creating vacancy"))
		return
	}
}
