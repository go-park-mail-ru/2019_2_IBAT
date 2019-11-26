package handler

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/websocket"

// 	"2019_2_IBAT/pkg/app/auth"
// 	. "2019_2_IBAT/pkg/pkg/interfaces"
// )

// // const (
// // 	// Time allowed to write a message to the peer.
// // 	writeWait = 10 * time.Second

// // 	// Time allowed to read the next pong message from the peer.
// // 	pongWait = 60 * time.Second

// // 	// Send pings to peer with this period. Must be less than pongWait.
// // 	pingPeriod = (pongWait * 9) / 10

// // 	maxMessageSize = 512
// // )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // type Connect struct {
// // 	Conn *websocket.Conn
// // 	Ch   chan uuid.UUID
// // }

// // type ConnectsPerUser struct {
// // 	Channels []chan uuid.UUID
// // 	Mu       *sync.Mutex
// // 	Ch       chan uuid.UUID
// // }

// // type WsConnects struct {
// // 	ConsMu   *sync.Mutex
// // 	Connects map[uuid.UUID]*ConnectsPerUser
// // }

// func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
// 	authInfo, ok := FromContext(r.Context())
// 	// if !ok {
// 	// 	log.Println("Notifications Handler: unauthorized")
// 	// 	w.WriteHeader(http.StatusUnauthorized)
// 	// 	return
// 	// }

// 	cookie, err := r.Cookie(auth.CookieName)
// 	if err != nil {
// 		fmt.Println("Failed to fetch cookie")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	fmt.Printf("Cookie: %s\n", cookie.Value)

// 	fmt.Println(r)

// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	conn := Connect{
// 		Conn: ws,
// 		Ch:   make(chan uuid.UUID, 1),
// 	}

// 	h.ConnectsPool.ConsMu.Lock()
// 	node, ok := h.ConnectsPool.Connects[authInfo.ID]

// 	if !ok {
// 		node = &ConnectsPerUser{
// 			Channels: make([]chan uuid.UUID, 0),
// 			Mu:       &sync.Mutex{},
// 			Ch:       make(chan uuid.UUID, 5),
// 		}
// 		node.Channels = append(node.Channels, conn.Ch)
// 		h.ConnectsPool.Connects[authInfo.ID] = node
// 		fmt.Printf("Connection pool for user %s was created\n", authInfo.ID)
// 	} else {
// 		node.Mu.Lock()
// 		node.Channels = append(node.Channels, conn.Ch) //careful with mu
// 		h.ConnectsPool.Connects[authInfo.ID] = node
// 		node.Mu.Unlock()
// 		fmt.Printf("Connection pool for user %s was updated\n", authInfo.ID)
// 	}
// 	h.ConnectsPool.ConsMu.Unlock() //careful

// 	go conn.ReadPump()
// 	go conn.WritePump()

// 	go sendNewMsgNotifications(node)
// 	// fmt.Println(h.WsConnects)
// }

// func sendNewMsgNotifications(clientConn *ConnectsPerUser) {
// 	for {
// 		// select {
// 		id := <-clientConn.Ch
// 		fmt.Printf("id %s got from channel for user", id.String())
// 		clientConn.Mu.Lock()

// 		for _, ch := range clientConn.Channels {
// 			ch <- id
// 			fmt.Printf("id %s sent to user\n", id.String())
// 		}
// 		clientConn.Mu.Unlock()
// 		// }_
// 	}
// }
