package repository

const (
	InsertChat = "INSERT INTO chats(chat_id, seeker_id, employer_id)VALUES" +
		"($1, $2, $3);"
	InsertMessage = "INSERT INTO messages(chat_id, owner_id, content)VALUES" +
		"($1, $2, $3);"
	SelectChatsForEmpl = "SELECT * FROM chats WHERE employer_id = $1;"
	SelectChatsForSeek = "SELECT * FROM chats WHERE seeker_id = $1;"

	SelectChatHistoryForEmpl = "SELECT M.chat_id, M.owner_id, M.created_at, M.content " +
		"FROM messages AS M JOIN chats AS C ON (M.chat_id = C.chat_id) WHERE M.chat_id = $1 AND C.employer_id = $2;"
	SelectChatHistoryForSeek = "SELECT M.chat_id, M.owner_id, M.created_at, M.content " +
		"FROM messages AS M JOIN chats AS C ON (M.chat_id = C.chat_id) WHERE M.chat_id = $1 AND C.seeker_id = $2;"
)

// -- INSERT INTO messages(message_id, chat_id, owner_id, content)
// -- VALUES(gen_random_uuid(),gen_random_uuid(),gen_random_uuid(), 'default msg');

// type OutChatMessage struct {
// 	ChatID uuid.UUID `json:"chat_id"                 db:"id"`
// 	// OwnerInfo AuthStorageValue `json:"-"                 db:"-"`
// 	// OwnerId
// 	Timestamp string `json:"timestamp"                 db:"id"`
// 	Text      string `json:"text"                 db:"id"`
// }

// type Chat struct {
// 	ChatID   uuid.UUID `json:"chat_id"         db:"chat_id"`
// 	SeekerID uuid.UUID `json:"seeker_id"       db:"seeker_id"`
// 	Employer uuid.UUID `json:"employer_id"     db:"employer_id"`
// }

// CREATE TABLE chats(
//     chat_id uuid PRIMARY KEY,
//     seeker_id uuid REFERENCES persons(id) ON DELETE CASCADE NOT NULL,
//     employer_id uuid REFERENCES persons(id) ON DELETE CASCADE  NOT NULL,
//     UNIQUE(seeker_id, employer_id)
// );

// CREATE TABLE messages(
//     chat_id uuid REFERENCES chats(chat_id) ON DELETE CASCADE NOT NULL,
//     owner_id uuid REFERENCES persons(id) ON DELETE CASCADE  NOT NULL,
//     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//     content text NOT NULL
// );
