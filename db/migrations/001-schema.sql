-- +migrate Up
CREATE TABLE IF NOT EXISTS calls (
    id INTEGER PRIMARY KEY,
    title TEXT UNIQUE NOT NULL,
    deadline TEXT NOT NULL,
    requirements TEXT NOT NULL,
    apply_module TEXT NOT NULL
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS bot_users (
    chat_id INTEGER PRIMARY KEY,
    first_interaction TEXT NOT NULL,
    is_active INTEGER DEFAULT 1 NOT NULL
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS users_calls (
    user_chat_id INTEGER NOT NULL,
    call_id INTEGER NOT NULL,
    FOREIGN KEY (user_chat_id) REFERENCES bot_users (chat_id),
    FOREIGN KEY (call_id) REFERENCES calls (id),
    PRIMARY KEY (user_chat_id, call_id)
);

-- +migrate Down
DROP TABLE IF EXISTS calls;
-- +migrate Down
DROP TABLE IF EXISTS bot_users;
-- +migrate Down
DROP TABLE IF EXISTS users_calls;

