package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	"2019_2_IBAT/internal/pkg/auth/session"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"sync"

	// "2019_2_IBAT/internal/pkg/pool"
	"2019_2_IBAT/internal/pkg/users"

	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const publicDir = "/static"

const maxUploadSize = 2 * 1024 * 1024 // 2 mb

type Handler struct {
	// Pool        *pool.Pool //
	WsConnects map[string]Connections

	InternalDir string
	AuthService session.ServiceClient
	UserService users.Service
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

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

	authInfo, ok := FromContext(r.Context())

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

	sessionBool, err := h.AuthService.DeleteSession(r.Context(), &session.Cookie{
		Cookie: cookie.Value,
	})
	if !sessionBool.Ok {
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

	authInfo, ok := FromContext(r.Context())

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

		authInfo, ok := FromContext(r.Context())

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := FromContext(r.Context())
	// if !ok {
	// 	log.Println("Notifications Handler: unauthorized")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		fmt.Println("Failed to fetch cookie")
	}
	fmt.Printf("Cookie: %s\n", cookie.Value)

	fmt.Println(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	node, ok := h.WsConnects[authInfo.ID.String()] //TODO can be deleted after check
	if !ok {
		node := Connections{
			Conns: []*websocket.Conn{ws},
			Mu:    &sync.Mutex{},
		}
		h.WsConnects[authInfo.ID.String()] = node
	} else {
		node.Mu.Lock()
		node.Conns = append(node.Conns, ws) //careful with mu
		h.WsConnects[authInfo.ID.String()] = node
		node.Mu.Unlock()
	}
	// go sendNewMsgNotifications(ws)
	fmt.Println(h.WsConnects)
}

func sendNewMsgNotifications(client *websocket.Conn) {
	// ticker := time.NewTicker(10 * time.Second)
	// for {
	// 	w, err := client.NextWriter(websocket.TextMessage)
	// 	if err != nil {
	// 		ticker.Stop()
	// 		break
	// 	}

	// 	msg := newMessage()
	// 	w.Write(msg)
	// 	w.Close()
	// 	<-ticker.C
	// }
}
