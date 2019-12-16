package service

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/app/auth"

	. "2019_2_IBAT/pkg/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s Service) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := FromContext(r.Context())
	if !ok {
		log.Println("Notifications Handler: unauthorized")
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

	fmt.Println(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	conn := Connect{
		Conn:   ws,
		Ch:     make(chan uuid.UUID, 1),
		UserId: authInfo.ID,
	}

	s.ConnectsPool.ConsMu.Lock()
	node, ok := s.ConnectsPool.Connects[authInfo.ID]

	if !ok {
		node = &ConnectsPerUser{
			Connects: make([]*Connect, 0),
			Mu:       &sync.Mutex{},
			// Ch:       make(chan uuid.UUID, 5),
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
