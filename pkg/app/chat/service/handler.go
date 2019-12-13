package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/app/auth"
	. "2019_2_IBAT/pkg/app/chat/models"
	. "2019_2_IBAT/pkg/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s Service) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	compId, err := uuid.Parse(mux.Vars(r)["companion_id"])
	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	id, err := s.CreateChat(authInfo, compId)
	if err != nil {
		var code int
		switch err.Error() {
		case InvalidIdMsg:
			code = http.StatusBadRequest
			// case UnauthorizedMsg:
			// 	code = http.StatusUnauthorized
			// case InternalErrorMsg:
			// 	code = http.StatusInternalServerError
			// default:
			// 	code = http.StatusBadRequest
		}
		SetError(w, code, err.Error())

		return
	}

	idJSON, _ := Id{Id: id.String()}.MarshalJSON()
	w.Write(idJSON)
}

func (s Service) HandleChat(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := FromContext(r.Context())
	if !ok {
		log.Println("HandleChat: unauthorized")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		fmt.Println("Failed to fetch cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Printf("Cookie: %s\n", cookie.Value)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	conn := Connect{
		Conn:   ws,
		Ch:     make(chan InChatMessage, 5),
		UserId: authInfo.ID,
	}

	s.ConnectsPool.ConsMu.Lock()
	node, ok := s.ConnectsPool.Connects[authInfo.ID]

	if !ok {
		node = &ConnectsPerUser{
			Connects: make([]*Connect, 0),
			Mu:       &sync.Mutex{},
		}
		node.Connects = append(node.Connects, &conn)
		conn.ConnIndex = 0
		s.ConnectsPool.Connects[authInfo.ID] = node

		fmt.Printf("Connection pool for user %s was created\n", authInfo.ID)
	} else {
		node.Mu.Lock()                               //?
		node.Connects = append(node.Connects, &conn) //careful with mu
		s.ConnectsPool.Connects[authInfo.ID] = node
		conn.ConnIndex = len(node.Connects) - 1
		node.Mu.Unlock() //?

		fmt.Printf("Connection pool for user %s was updated\n", authInfo.ID)
	}
	s.ConnectsPool.ConsMu.Unlock()

	mu := sync.Mutex{}
	stopCh := make(chan bool, 1)

	go s.ReadPump(&conn, authInfo, stopCh, &mu)
	go s.WritePump(&conn, stopCh, &mu)
}

func (s Service) DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	answer, _ := json.Marshal("dummy answer")

	w.Write(answer)
}

func (s Service) HandlerGetChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	chats, err := s.GetChats(authInfo)
	if err != nil {
		code := 400
		switch err.Error() {
		case InvalidIdMsg:
			code = http.StatusBadRequest
			// case UnauthorizedMsg:
			// 	code = http.StatusUnauthorized
			// case InternalErrorMsg:
			// 	code = http.StatusInternalServerError
			// default:
			// 	code = http.StatusBadRequest
		}
		SetError(w, code, err.Error())

		return
	}

	answer, _ := json.Marshal(chats)
	w.Write(answer)
}

func (s Service) HandlerGetChatHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	chatId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	chats, err := s.GetChatHistory(authInfo, chatId)
	if err != nil {
		code := 400
		switch err.Error() {
		case InvalidIdMsg:
			code = http.StatusBadRequest
			// case UnauthorizedMsg:
			// 	code = http.StatusUnauthorized
			// case InternalErrorMsg:
			// 	code = http.StatusInternalServerError
			// default:
			// 	code = http.StatusBadRequest
		}
		SetError(w, code, err.Error())

		return
	}

	answer, _ := json.Marshal(chats)
	w.Write(answer)
}
