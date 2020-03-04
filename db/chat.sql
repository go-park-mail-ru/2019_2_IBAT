DROP TABLE messages;
DROP TABLE chats;


CREATE TABLE chats(
    chat_id uuid PRIMARY KEY,
    seeker_id uuid REFERENCES persons(id) ON DELETE CASCADE NOT NULL,
    employer_id uuid REFERENCES persons(id) ON DELETE CASCADE  NOT NULL,
    UNIQUE(seeker_id, employer_id)
);

CREATE TABLE messages(
    chat_id uuid REFERENCES chats(chat_id) ON DELETE CASCADE NOT NULL,
    owner_id uuid REFERENCES persons(id) ON DELETE CASCADE  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    content text NOT NULL
);


CREATE INDEX chats_seeker_idx ON chats(seeker_id);
CREATE INDEX chats_employer_idx ON chats(employer_id);

CREATE INDEX messages_idx ON messages(chat_id, created_at);


-- INSERT INTO messages(message_id, chat_id, owner_id, content)
-- VALUES(gen_random_uuid(),gen_random_uuid(),gen_random_uuid(), 'default msg');