package repository

import (
	. "2019_2_IBAT/pkg/app/chat/models"
	. "2019_2_IBAT/pkg/pkg/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBStorage struct {
	DbConn *sqlx.DB
}

func (m DBStorage) CreateChat(seekerId uuid.UUID, employerId uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	_, err := m.DbConn.Exec(InsertChat, id, seekerId, employerId)
	if err != nil {
		fmt.Printf("CreateChat: %s\n", err)
		return id, errors.New(InvalidIdMsg)
	}
	return id, nil
}

func (m DBStorage) GetCompanionId(msg InChatMessage) (uuid.UUID, error) {
	var id uuid.UUID
	var err error

	if msg.OwnerInfo.Role == SeekerStr {
		err = m.DbConn.QueryRowx("SELECT employer_id FROM chats WHERE chat_id = $1 AND "+
			"seeker_id = $2;", msg.ChatID, msg.OwnerInfo.ID).Scan(&id)
	} else if msg.OwnerInfo.Role == EmployerStr {
		err = m.DbConn.QueryRowx("SELECT seeker_id FROM chats WHERE chat_id = $1 AND "+
			"employer_id = $2;", msg.ChatID, msg.OwnerInfo.ID).Scan(&id)
	}

	if err != nil {
		fmt.Printf("GetCompanionId: %s\n", err)
		return id, errors.New(InvalidIdMsg)
	}

	return id, nil
}

func (m DBStorage) GetChats(authInfo AuthStorageValue) ([]Chat, error) {
	var rows *sqlx.Rows
	var err error

	if authInfo.Role == EmployerStr {
		rows, err = m.DbConn.Queryx(SelectChatsForEmpl, authInfo.ID)
	} else { //if authInfo.Role == SeekerStr {
		rows, err = m.DbConn.Queryx(SelectChatsForSeek, authInfo.ID)
	}

	if err != nil {
		fmt.Printf("GetChats: %s\n", err)
		return []Chat{}, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	var messages []Chat
	for rows.Next() {
		var chat Chat

		err = rows.StructScan(&chat)

		if err != nil {
			fmt.Printf("GetChats: %s\n", err)
			return messages, errors.New(InternalErrorMsg)
		}
		messages = append(messages, chat)
	}

	return messages, nil
}

func (m DBStorage) CreateMessage(msg InChatMessage) error {
	_, err := m.DbConn.Exec(InsertMessage, msg.ChatID, msg.OwnerInfo.ID, msg.Text)
	if err != nil {
		fmt.Printf("CreateMessage: %s\n", err)
		return errors.New(BadRequestMsg)
	}
	return nil
}

func (m DBStorage) GetChatHistory(authInfo AuthStorageValue, chatId uuid.UUID) ([]OutChatMessage, error) {
	var rows *sqlx.Rows
	var err error

	if authInfo.Role == EmployerStr {
		rows, err = m.DbConn.Queryx(SelectChatHistoryForEmpl, chatId, authInfo.ID)
	} else {
		rows, err = m.DbConn.Queryx(SelectChatHistoryForSeek, chatId, authInfo.ID)
	}

	if err != nil {
		fmt.Printf("GetChatHistory: %s\n", err)
		return []OutChatMessage{}, errors.New(BadRequestMsg)
	}
	defer rows.Close()

	var messages []OutChatMessage
	for rows.Next() {
		var message OutChatMessage

		err = rows.StructScan(&message)

		if err != nil {
			fmt.Printf("GetChatHistoryForSeeker: %s\n", err)
			return messages, errors.New(InternalErrorMsg)
		}
		messages = append(messages, message)
	}

	return messages, nil
}
