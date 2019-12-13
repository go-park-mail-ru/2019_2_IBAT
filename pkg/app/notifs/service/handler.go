package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"

	. "2019_2_IBAT/pkg/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h Service) HandleNotifications(w http.ResponseWriter, r *http.Request) {
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
		Conn: ws,
		Ch:   make(chan uuid.UUID, 1),
	}

	h.ConnectsPool.ConsMu.Lock()
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
	h.ConnectsPool.ConsMu.Unlock()

	go conn.ReadPump()
	go conn.WritePump()

	go sendNewMsgNotifications(node)
}

func sendNewMsgNotifications(clientConn *ConnectsPerUser) {
	for {
		id := <-clientConn.Ch
		fmt.Printf("id %s got from channel for user", id.String())
		clientConn.Mu.Lock()

		for _, ch := range clientConn.Channels {
			ch <- id
			fmt.Printf("id %s sent to user\n", id.String())
		}
		clientConn.Mu.Unlock()
	}
}

func (h Service) Notifications() {
	for {
		notif := <-h.NotifChan
		fmt.Println("Notification accepted")
		fmt.Println(notif)

		ctx := context.Background()
		idsMsg, err := h.RecomService.GetUsersForTags(ctx,
			&recomsproto.IDsMessage{IDs: UuidsToStrings(notif.TagIDs)},
		)

		ids := StringsToUuids(idsMsg.IDs)

		fmt.Println("Users ids intrested in new vacancy")
		fmt.Println(ids)

		if err != nil {
			log.Printf("Notifications %s", err)
		}
		fmt.Println("connects.ConsMu.Lock()")

		h.ConnectsPool.ConsMu.Lock()
		for _, id := range ids {
			fmt.Println("Notification ready  to be sent to user")
			fmt.Println(h.ConnectsPool.Connects[id])
			if _, ok := h.ConnectsPool.Connects[id]; ok {
				h.ConnectsPool.Connects[id].Ch <- notif.VacancyId
				fmt.Printf("Notification was sent to user %s\n", id.String())
			} else {
				fmt.Printf("Notification can not be sent to user %s\n", id.String())
				delete(h.ConnectsPool.Connects, id)
				fmt.Printf("Connections to user were removed\nd")
			}
		}
		h.ConnectsPool.ConsMu.Unlock()
		fmt.Println("connects.ConsMu.Unlock()")

		fmt.Println(ids)
	}
}
