package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	. "2019_2_IBAT/pkg/pkg/models"
)

type NotifStruct struct {
	VacancyId uuid.UUID
	TagIDs    []uuid.UUID
}

type Connect struct {
	Conn *websocket.Conn
	Ch   chan uuid.UUID
}

type ConnectsPerUser struct {
	Channels []chan uuid.UUID
	Mu       *sync.Mutex
	Ch       chan uuid.UUID
}

type WsConnects struct {
	ConsMu   *sync.Mutex
	Connects map[uuid.UUID]*ConnectsPerUser
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

func (c *Connect) ReadPump() {
	defer func() {
		fmt.Println("ReadPump CONNECTION WAS CLOSED")
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		if _, _, err := c.Conn.NextReader(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Connect) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		fmt.Println("WritePump CONNECTION WAS CLOSED")
	}()

	for {
		select {
		case id := <-c.Ch: //ok
			// if ok {
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				ticker.Stop()
				break
			}
			idJSON, _ := Id{Id: id.String()}.MarshalJSON()
			w.Write(idJSON)
			w.Close()
			fmt.Printf("id %s was sent user", id.String())
			// } else {
			// 	fmt.Println("Channel closed!")
			// 	//close channel
			// 	return
			// 	// break
			// }
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
