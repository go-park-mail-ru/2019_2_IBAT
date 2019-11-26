package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/app/auth"
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

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

	//TODO can be deleted after check
	node, ok := h.ConnectsPool.Connects[authInfo.ID]
	if !ok {
		node = &ConnectsPerUser{
			Conns: []*websocket.Conn{ws},
			Mu:    &sync.Mutex{},
			Ch:    make(chan uuid.UUID, 5),
		}
		h.ConnectsPool.Connects[authInfo.ID] = node
		fmt.Printf("Connection pool for user %s was created\n", authInfo.ID)
	} else {
		node.Mu.Lock()
		node.Conns = append(node.Conns, ws) //careful with mu
		// h.WsConnects[authInfo.ID.String()] = node
		h.ConnectsPool.Connects[authInfo.ID] = node
		node.Mu.Unlock()
		fmt.Printf("Connection pool for user %s was updated\n", authInfo.ID)
	}
	h.ConnectsPool.ConsMu.Unlock() //careful

	go sendNewMsgNotifications(node)
	// fmt.Println(h.WsConnects)
}

func sendNewMsgNotifications(clientConn *ConnectsPerUser) {
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
	for {
		select {
		case id, _ := <-clientConn.Ch:
			// if !ok {
			// 	continue
			// }

			fmt.Printf("id %s got from channel for user", id.String())
			clientConn.Mu.Lock()
			for i, client := range clientConn.Conns {
				// select {
				// case client <- message:
				// default:
				// 	close(client.send)
				// 	delete(h.clients, client)
				// }
				// client.
				if client != nil {
					w, _ := client.NextWriter(websocket.TextMessage)
					// if err != nil {
					// 	ticker.Stop()
					// 	break
					// }

					idJSON, _ := json.Marshal(Id{Id: id.String()})
					w.Write(idJSON)
					w.Close()
					fmt.Printf("id %s was sent user", id.String())
				} else {
					// ticker.Stop()
					fmt.Println("connection disconnected")
					clientConn.Conns[i] = clientConn.Conns[len(clientConn.Conns)-1]
					clientConn.Conns[len(clientConn.Conns)-1] = nil
					clientConn.Conns = clientConn.Conns[:len(clientConn.Conns)-1]
				}
			}
			clientConn.Mu.Unlock()
		}
	}
}

// func readPump() {
// 	defer func() {
// 		c.hub.unregister <- c
// 		c.conn.Close()
// 	}()
// 	c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		_, message, err := c.conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
// 		c.hub.broadcast <- message
// 	}
// }
