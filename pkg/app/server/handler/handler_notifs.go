package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/app/auth"
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 60 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	maxMessageSize = 512
// )

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// type Connect struct {
// 	Conn *websocket.Conn
// 	Ch   chan uuid.UUID
// }

// type ConnectsPerUser struct {
// 	Channels []chan uuid.UUID
// 	Mu       *sync.Mutex
// 	Ch       chan uuid.UUID
// }

// type WsConnects struct {
// 	ConsMu   *sync.Mutex
// 	Connects map[uuid.UUID]*ConnectsPerUser
// }

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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Printf("Cookie: %s\n", cookie.Value)

	fmt.Println(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	h.ConnectsPool.ConsMu.Lock()

	conn := Connect{
		Conn: ws,
		Ch:   make(chan uuid.UUID, 1),
	}

	//TODO can be deleted after check
	node, ok := h.ConnectsPool.Connects[authInfo.ID]

	if !ok {
		node = &ConnectsPerUser{
			Channels: make([]chan uuid.UUID, 0),
			Mu:       &sync.Mutex{},
			Ch:       make(chan uuid.UUID, 5),
		}
		node.Channels = append(node.Channels, conn.Ch)
		h.ConnectsPool.Connects[authInfo.ID] = node
		fmt.Printf("Connection pool for user %s was created\n", authInfo.ID)
	} else {
		node.Mu.Lock()
		node.Channels = append(node.Channels, conn.Ch) //careful with mu
		h.ConnectsPool.Connects[authInfo.ID] = node
		node.Mu.Unlock()
		fmt.Printf("Connection pool for user %s was updated\n", authInfo.ID)
	}
	h.ConnectsPool.ConsMu.Unlock() //careful

	go conn.ReadPump()
	go conn.WritePump()

	go sendNewMsgNotifications(node)
	// fmt.Println(h.WsConnects)
}

func sendNewMsgNotifications(clientConn *ConnectsPerUser) {
	for {
		// select {
		id := <-clientConn.Ch
		fmt.Printf("id %s got from channel for user", id.String())
		clientConn.Mu.Lock()

		for _, ch := range clientConn.Channels {
			ch <- id
			fmt.Printf("id %s sent to user\n", id.String())
		}
		clientConn.Mu.Unlock()
		// }_
	}
	//clear gorutine
}

// func (c *Connect) readPump() {
// 	// defer func() {
// 	// 	c.hub.unregister <- c
// 	// 	c.conn.Close()
// 	// }()
// 	c.Conn.SetReadLimit(maxMessageSize)
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		if _, _, err := c.Conn.NextReader(); err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			// c.Close()
// 			break
// 		}
// 	}
// }

// // for {
// // 	_, message, err := c.conn.ReadMessage()
// // 	if err != nil {
// // 		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// // 			log.Printf("error: %v", err)
// // 		}
// // 		break
// // 	}
// // 	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
// // 	c.hub.broadcast <- message
// // }

// func (c *Connect) writePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()
// 	for {
// 		select {
// 		case id, ok := <-c.Ch:
// 			if ok {
// 				w, _ := c.Conn.NextWriter(websocket.TextMessage)
// 				// if err != nil {
// 				// 	ticker.Stop()
// 				// 	break
// 				// }
// 				idJSON, _ := json.Marshal(Id{Id: id.String()})
// 				w.Write(idJSON)
// 				w.Close()
// 				fmt.Printf("id %s was sent user", id.String())
// 			} else {
// 				fmt.Println("Channel closed!")
// 				//close channel
// 				return
// 				// break
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
