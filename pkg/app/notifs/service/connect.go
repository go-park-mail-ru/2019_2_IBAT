package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/pkg/config"
	. "2019_2_IBAT/pkg/pkg/models"
)

type Connect struct {
	Conn      *websocket.Conn
	Ch        chan uuid.UUID
	UserId    uuid.UUID
	ConnIndex int
}

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

		fmt.Println("Notif ReadPump CONNECTION WAS CLOSED")
	}()

	c.Conn.SetReadLimit(config.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(config.PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
	for {
		if _, _, err := c.Conn.NextReader(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (s Service) WritePump(c *Connect, stopCh chan bool, mu *sync.Mutex) {
	ticker := time.NewTicker(config.PingPeriod)
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

		fmt.Println("Notif WritePump CONNECTION WAS CLOSED")
	}()

	for {
		select {
		case id := <-c.Ch: //ok
			c.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				ticker.Stop()
				break
			}
			idJSON, _ := Id{Id: id.String()}.MarshalJSON()
			w.Write(idJSON)
			w.Close()
			fmt.Printf("id %s was sent user", id.String())
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// func (s Service) WritePump(c *Connect, stopCh chan bool, mu *sync.Mutex) {
// 	ticker := time.NewTicker(config.PingPeriod)
// 	defer func() {
// 		ticker.Stop()

// 		mu.Lock()
// 		select {
// 		case <-stopCh:
// 			close(stopCh)
// 			return
// 		default:
// 			stopCh <- true
// 			s.RemoveConnect(c)
// 		}
// 		mu.Unlock()

// 		fmt.Println("WritePump CONNECTION WAS CLOSED")
// 	}()

// 	for {
// 		select {
// 		case msg := <-c.Ch:
// 			c.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				ticker.Stop()
// 				// break
// 				return
// 			}

// 			// finalMsg := OutChatMessage{
// 			// 	ChatID:    msg.ChatID,
// 			// 	Timestamp: msg.Timestamp,
// 			// 	Text:      msg.Text,
// 			// }

// 			messageJSON, _ := json.Marshal(msg)
// 			w.Write(messageJSON)

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 			fmt.Printf("WritePump msg %s was sent to user\n", msg)
// 		case <-ticker.C:
// 			err := c.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
// 			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 				// break
// 			}
// 		}
// 	}
// 	// fmt.Println("WritePump CYCLE FOR USER %s WAS STOPPED")
// }

func (s Service) RemoveConnect(c *Connect) {
	s.ConnectsPool.ConsMu.Lock()
	c.Conn.Close()
	close(c.Ch)

	node := s.ConnectsPool.Connects[c.UserId]

	if len(node.Connects) > 1 {
		node.Connects[len(node.Connects)-1].ConnIndex = c.ConnIndex
		node.Connects[c.ConnIndex] = node.Connects[len(node.Connects)-1]
		node.Connects = node.Connects[:len(node.Connects)-1]

		s.ConnectsPool.Connects[c.UserId] = node
		fmt.Printf("RemoveConnect: connections node len became %d\n", len(node.Connects))
	} else {
		delete(s.ConnectsPool.Connects, c.UserId)
		fmt.Println("RemoveConnect: connections node was removed")
	}

	s.ConnectsPool.ConsMu.Unlock()
}
