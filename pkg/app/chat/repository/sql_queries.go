package repository

const (
	InsertChat = "INSERT INTO chats(chat_id, seeker_id, employer_id)VALUES" +
		"($1, $2, $3);"
	InsertMessage = "INSERT INTO messages(chat_id, owner_id, content)VALUES" +
		"($1, $2, $3);"
	SelectChatsForEmpl = "SELECT C.chat_id, C.seeker_id, P.first_name, P.second_name FROM chats AS C " +
		"JOIN persons AS P ON (C.seeker_id = P.id) WHERE C.employer_id = $1;"
	SelectChatsForSeek = "SELECT C.chat_id, C.employer_id, COMP.company_name FROM chats AS C " +
		"JOIN persons AS P ON (C.employer_id = P.id) JOIN companies AS COMP ON (P.id = COMP.own_id) " +
		"WHERE C.seeker_id = $1;"

	SelectChatHistoryForEmpl = "SELECT M.chat_id, M.owner_id, M.created_at, M.content " +
		"FROM messages AS M JOIN chats AS C ON (M.chat_id = C.chat_id) WHERE M.chat_id = $1 AND C.employer_id = $2;"
	SelectChatHistoryForSeek = "SELECT M.chat_id, M.owner_id, M.created_at, M.content FROM messages AS M " +
		"JOIN chats AS C ON (M.chat_id = C.chat_id) WHERE M.chat_id = $1 AND C.seeker_id = $2;"
	SelectCompanyName = "SELECT COMP.company_name FROM chats AS C " +
		"JOIN companies AS COMP ON COMP.own_id = C.employer_id " +
		"WHERE C.chat_id = $1"
	SelectSeekerName = "SELECT P.first_name, P.second_name FROM chats AS C " +
		"JOIN persons AS P ON P.id = C.seeker_id " +
		"WHERE C.chat_id = $1"
)

// type OutChatMessage struct {
// 	ChatID uuid.UUID `json:"chat_id"              db:"chat_id"`
// 	// OwnerInfo AuthStorageValue `json:"-"                 db:"-"`
// 	OwnerId   uuid.UUID `json:"owner_id"          db:"owner_id"`
// 	OwnerName uuid.UUID `json:"owner_name"        db:"owner_name"`
// 	Timestamp string    `json:"created_at"        db:"created_at"`
// 	Text      string    `json:"content"           db:"content"`
// 	IsYours   bool      `json:"is_yours"          db:"content"`
// }
