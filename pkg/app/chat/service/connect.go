package service

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	. "2019_2_IBAT/pkg/app/chat/models"
	. "2019_2_IBAT/pkg/pkg/models"
)

type Connect struct {
	Conn      *websocket.Conn
	Ch        chan InChatMessage
	UserId    uuid.UUID
	ConnIndex int
}

// func (s Service) ReadPump(c *Connect) {
// 	defer func() {
// 		fmt.Println("ReadPump CONNECTION WAS CLOSED")
// 		c.Conn.Close()
// 	}()

// 	c.Conn.SetReadLimit(maxMessageSize)
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		if _, _, err := c.Conn.NextReader(); err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 	}
// }

// func (c *Client) readPump() {
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

func (s Service) ReadPump(c *Connect, authInfo AuthStorageValue, stopCh chan bool, mu *sync.Mutex) {
	defer func() {
		mu.Lock()
		select {
		case <-stopCh:
			close(stopCh)
			return
		default:
			stopCh <- true
			s.RemoveConnect(c)
		}
		mu.Unlock()

		fmt.Println("ReadPump CONNECTION WAS CLOSED")
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg InChatMessage
		err = json.Unmarshal(bytes, &msg)
		if err == nil {
			msg.OwnerInfo = authInfo
			fmt.Printf("ReadPump msg %s was read\n", msg.Text)

			s.MainChan <- msg
			fmt.Printf("ReadPump msg %s was send to main channel\n", msg.Text)
		} else {
			fmt.Printf("ReadPump invalid message\n")
		}
	}
	fmt.Printf("ReadPump CYCLE FOR USER %s WAS STOPPED\n", authInfo.ID)
	return //useless
}

func (s Service) WritePump(c *Connect, stopCh chan bool, mu *sync.Mutex) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		mu.Lock()
		select {
		case <-stopCh:
			close(stopCh)
			return
		default:
			stopCh <- true
			s.RemoveConnect(c)
		}
		mu.Unlock()

		fmt.Println("WritePump CONNECTION WAS CLOSED")
	}()

	for {
		select {
		case msg := <-c.Ch:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				ticker.Stop()
				// break
				return
			}

			finalMsg := OutChatMessage{
				ChatID:    msg.ChatID,
				Timestamp: msg.Timestamp,
				Text:      msg.Text,
			}

			messageJSON, _ := json.Marshal(finalMsg)
			w.Write(messageJSON)

			if err := w.Close(); err != nil {
				return
			}
			fmt.Printf("WritePump msg %s was sent to user\n", finalMsg)
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
				// break
			}
		}
	}
	// fmt.Println("WritePump CYCLE FOR USER %s WAS STOPPED")
}

// func (c *Client) readPump() {
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

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// func (c *Client) writePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			// Add queued chat messages to the current websocket message.
// 			n := len(c.send)
// 			for i := 0; i < n; i++ {
// 				w.Write(newline)
// 				w.Write(<-c.send)
// 			}

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

func (s Service) RemoveConnect(c *Connect) {
	s.ConnectsPool.ConsMu.Lock()
	c.Conn.Close()
	close(c.Ch)

	node := s.ConnectsPool.Connects[c.UserId]

	if len(node.Connects) > 1 {
		node.Connects[len(node.Connects)-1].ConnIndex = c.ConnIndex
		node.Connects[c.ConnIndex] = node.Connects[len(node.Connects)-1] // Copy last element to index i.
		node.Connects = node.Connects[:len(node.Connects)-1]             // Truncate slice.

		s.ConnectsPool.Connects[c.UserId] = node
		fmt.Printf("RemoveConnect: connections node len became %d\n", len(node.Connects))
	} else {
		delete(s.ConnectsPool.Connects, c.UserId)
		fmt.Println("RemoveConnect: connections node was removed")
	}

	s.ConnectsPool.ConsMu.Unlock()
}
