package service

import (
	"2019_2_IBAT/pkg/app/chat"
	. "2019_2_IBAT/pkg/app/chat/models"
	. "2019_2_IBAT/pkg/pkg/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Service struct {
	MainChan     chan InChatMessage
	ConnectsPool WsConnects
	// AuthService  session.ServiceClient
	// RecomService recomsproto.ServiceClient
	Storage chat.Repository
}

func (s Service) CreateChat(authInfo AuthStorageValue, compId uuid.UUID) (uuid.UUID, error) {
	if authInfo.Role == SeekerStr {
		return s.Storage.CreateChat(authInfo.ID, compId)

	} else if authInfo.Role == EmployerStr {
		return s.Storage.CreateChat(compId, authInfo.ID)
	}
	return uuid.UUID{}, errors.New(ForbiddenMsg)
}

func (s Service) GetChats(authInfo AuthStorageValue) ([]Chat, error) {
	return s.Storage.GetChats(authInfo)
}

func (s Service) ProcessMessage() {
	for {
		msg := <-s.MainChan
		fmt.Printf("ProcessMessage msg %s was read from main channel\n", msg.Text)
		id, err := s.Storage.GetCompanionId(msg)
		fmt.Printf("ProcessMessage companion id %s was accepted\n", id.String())

		fmt.Println(id)
		fmt.Println(err)
		s.ConnectsPool.ConsMu.Lock()
		if cons, ok := s.ConnectsPool.Connects[id]; ok {
			s.ConnectsPool.Connects[id].Mu.Lock()
			// cons := s.ConnectsPool.Connects[id]
			for _, con := range cons.Connects {
				fmt.Printf("ProcessMessage msg %s was sent to output channel\n", msg.Text)
				con.Ch <- msg
			}
			s.ConnectsPool.Connects[id].Mu.Unlock()
		}
		s.ConnectsPool.ConsMu.Unlock()

		// if not ok set unread
		go s.Storage.CreateMessage(msg)
	}
}

func (s Service) GetChatHistory(authInfo AuthStorageValue, chatId uuid.UUID) ([]OutChatMessage, error) {
	if authInfo.Role == EmployerStr {
		return s.Storage.GetChatHistoryForEmployer(authInfo, chatId)

	} else {
		return s.Storage.GetChatHistoryForSeeker(authInfo, chatId)
	}
}
