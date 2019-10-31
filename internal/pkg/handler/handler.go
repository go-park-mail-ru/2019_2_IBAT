package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"

	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const publicDir = "/static"

const maxUploadSize = 2 * 1024 * 1024 // 2 mb

type Handler struct {
	InternalDir string
	AuthService auth.Service
	UserService users.Service
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, role, ok := h.UserService.CheckUser(userAuthInput.Email, userAuthInput.Password)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	cookie, role, err := h.AuthService.CreateSession(id, role)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	// authInfo, ok := h.AuthService.Storage.Get(cookie.Value)
	authInfo, ok := h.AuthService.GetSession(cookie.Value)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	RoleJSON, _ := json.Marshal(Role{authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthService.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateSeeker(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	// authInfo, cookieValue, err := h.AuthService.Storage.Set(uuid, SeekerStr) //possible return authInfo
	authInfo, cookieValue, err := h.AuthService.SetRecord(uuid, SeekerStr) //possible return authInfo

	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "error while unmarshaling")
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	} //impossible error

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}
	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	// authInfo, cookieValue, err := h.AuthService.Storage.Set(uuid, EmployerStr) //possible return authInfo
	authInfo, cookieValue, err := h.AuthService.SetRecord(uuid, EmployerStr) //possible return authInfo

	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "error while unmarshaling")
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	} //impossible error

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) CreateResume(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateResume(r.Body, cookie.Value, h.AuthService) //.Storage
	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) DeleteResume(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	strId := mux.Vars(r)["id"]
	resId, err := uuid.Parse(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}
	err = h.UserService.DeleteResume(resId, cookie.Value, h.AuthService)

	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	resume, err := h.UserService.GetResume(resId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	resumeJSON, _ := json.Marshal(resume)

	w.Write([]byte(resumeJSON))
}

func (h *Handler) PutResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	defer r.Body.Close()
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	resId, err := uuid.Parse(mux.Vars(r)["id"])

	err = h.UserService.PutResume(resId, r.Body, cookie.Value, h.AuthService)

	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)

	if !ok {
		log.Println("GetUser Handler: unauthorized")
		// log.Println("GetUser Handler: unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if authInfo.Role == SeekerStr {
		seeker, err := h.UserService.GetSeeker(authInfo.ID)

		if err != nil {
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
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteUser(cookie.Value, h.AuthService)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthService.DeleteSession(cookie)
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

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, ok := h.AuthService.GetSession(cookie.Value) //impossible error, should use only Set method
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if authInfo.Role == SeekerStr {
		err = h.UserService.PutSeeker(r.Body, authInfo.ID)
	} else if authInfo.Role == EmployerStr {
		err = h.UserService.PutEmployer(r.Body, authInfo.ID)
	}

	if err != nil {
		w.WriteHeader(http.StatusForbidden) //should add invalid email case
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}
}

//should test method
func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	seeker, err := h.UserService.GetSeeker(seekId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	seeker.Password = "" //danger
	seekerJSON, _ := json.Marshal(seeker)

	w.Write([]byte(seekerJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	emplId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer, err := h.UserService.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer.Password = "" //danger
	employerJSON, _ := json.Marshal(employer)

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	employers, _ := h.UserService.GetEmployers()

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
	// 	w.Write([]byte(errJSON))
	// 	return
	// }

	// for i, item := range employers {
	// 	item.Password = ""
	// 	employers[i] = item
	// }

	employerJSON, _ := json.Marshal(employers)

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetResumes(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resumes, _ := h.UserService.GetResumes() //error handling

	resumesJSON, _ := json.Marshal(resumes)

	w.Write([]byte(resumesJSON))
}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vacancies, _ := h.UserService.GetVacancies() //err handle

	vacanciesJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(vacanciesJSON))

}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateVacancy(r.Body, cookie.Value, h.AuthService)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacancy, err := h.UserService.GetVacancy(vacId)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacancyJSON, _ := json.Marshal(vacancy)

	w.Write([]byte(vacancyJSON))
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: "Invalid id"})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteVacancy(vacId, cookie.Value, h.AuthService)

	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.PutVacancy(vacId, r.Body, cookie.Value, h.AuthService)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) UploadFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
			w.Write([]byte(errJSON))
			return
		}

		authInfo, ok := h.AuthService.GetSession(cookie.Value)

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

func (h *Handler) GetResponds(w http.ResponseWriter, r *http.Request) { //+
	fmt.Println("Here")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//employers, _ := h.UserService.Storage.GetEmployers()

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		// errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		// w.Write([]byte(errJSON))
		// return
		fmt.Println("No cookie")
		cookie.Value = ""
	}

	v := r.URL.Query()
	params := make(map[string]string)
	params["vacancyid"] = v.Get("vacancyid")
	params["resumeid"] = v.Get("resumeid")
	fmt.Printf("vacancyid = %s, resumeid = %s, id = %s", params["vacancyid"], params["resumeid"])

	responds, _ := h.UserService.GetResponds(cookie.Value, params, h.AuthService) //error handling

	respondsJSON, _ := json.Marshal(responds)

	w.Write([]byte(respondsJSON))

}

func (h *Handler) CreateRespond(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateRespond(r.Body, cookie.Value, h.AuthService)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{id.String()})

	w.Write([]byte(idJSON))
}
