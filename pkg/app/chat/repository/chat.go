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

func (m DBStorage) GetCompanionIdAndName(msg InChatMessage) (uuid.UUID, string, error) {
	var id uuid.UUID
	var err error
	var name string

	if msg.OwnerInfo.Role == SeekerStr {
		err = m.DbConn.QueryRowx("SELECT C.employer_id, COMP.company_name FROM chats AS C "+
			"JOIN companies AS COMP ON (COMP.own_id = C.employer_id) WHERE C.chat_id = $1 AND "+
			"C.seeker_id = $2;", msg.ChatID, msg.OwnerInfo.ID).Scan(&id, &name)
	} else if msg.OwnerInfo.Role == EmployerStr {
		var firstName, secondName string
		err = m.DbConn.QueryRowx("SELECT C.seeker_id, P.first_name, P.second_name FROM chats AS C "+
			"JOIN persons AS P ON (P.id = C.seeker_id) WHERE C.chat_id = $1 AND "+
			"C.employer_id = $2;", msg.ChatID, msg.OwnerInfo.ID).Scan(&id, &firstName, &secondName)
		name = firstName + " " + secondName
	}

	if err != nil {
		fmt.Printf("GetCompanionId: %s\n", err)
		return id, name, errors.New(InvalidIdMsg)
	}

	return id, name, nil
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

		if authInfo.Role == EmployerStr {
			var firstName, secondName string
			err = rows.Scan(&chat.ChatID, &chat.CompanionID, &firstName, &secondName)
			chat.CompanionName = firstName + " " + secondName
		} else {
			err = rows.Scan(&chat.ChatID, &chat.CompanionID, &chat.CompanionName)
		}

		if err != nil {
			fmt.Printf("GetChats: %s\n", err)
			return messages, errors.New(InternalErrorMsg)
		}
		messages = append(messages, chat)
	}

	return messages, nil
}

func (m DBStorage) CreateMessage(msg InChatMessage) error {
	_, err := m.DbConn.Exec(InsertMessage, msg.ChatID, msg.OwnerInfo.ID, msg.Text, msg.Timestamp)
	if err != nil {
		fmt.Printf("CreateMessage: %s\n", err)
		return errors.New(BadRequestMsg)
	}
	return nil
}

func (m DBStorage) GetChatHistoryForEmployer(authInfo AuthStorageValue, chatId uuid.UUID) ([]OutChatMessage, error) {
	var firstName, secondName string
	err := m.DbConn.QueryRowx(SelectSeekerName, chatId).Scan(&firstName, &secondName)
	if err != nil {
		fmt.Printf("GetChatHistoryForEmployer: %s\n", err)
		return []OutChatMessage{}, errors.New(BadRequestMsg)
	}

	rows, err := m.DbConn.Queryx(SelectChatHistoryForEmpl, chatId, authInfo.ID)

	if err != nil {
		fmt.Printf("GetChatHistoryForEmployer: %s\n", err)
		return []OutChatMessage{}, errors.New(BadRequestMsg)
	}
	defer rows.Close()

	var messages []OutChatMessage
	for rows.Next() {
		var msg OutChatMessage

		err = rows.Scan(&msg.ChatID, &msg.OwnerId, &msg.Timestamp, &msg.Text)

		if err != nil {
			fmt.Printf("GetChatHistoryForEmployer: %s\n", err)
			return messages, errors.New(InternalErrorMsg)
		}

		if msg.OwnerId != authInfo.ID {
			msg.OwnerName = firstName + " " + secondName
			msg.IsNotYours = true
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (m DBStorage) GetChatHistoryForSeeker(authInfo AuthStorageValue, chatId uuid.UUID) ([]OutChatMessage, error) {
	var companyName string

	err := m.DbConn.QueryRowx(SelectCompanyName, chatId).Scan(&companyName)
	if err != nil {
		fmt.Printf("GetChatHistoryForSeeker: %s\n", err)
		return []OutChatMessage{}, errors.New(BadRequestMsg)
	}

	rows, err := m.DbConn.Queryx(SelectChatHistoryForSeek, chatId, authInfo.ID)

	if err != nil {
		fmt.Printf("GetChatHistoryForSeeker: %s\n", err)
		return []OutChatMessage{}, errors.New(BadRequestMsg)
	}
	defer rows.Close()

	var messages []OutChatMessage
	for rows.Next() {
		var msg OutChatMessage

		err = rows.Scan(&msg.ChatID, &msg.OwnerId, &msg.Timestamp, &msg.Text)
		if err != nil {
			fmt.Printf("GetChatHistoryForSeeker: %s\n", err)
			return messages, errors.New(InternalErrorMsg)
		}

		if msg.OwnerId != authInfo.ID {
			msg.OwnerName = companyName
			msg.IsNotYours = true
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
