-- Patch rooms table
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='rooms' AND column_name='invite_code'
    ) THEN
        ALTER TABLE rooms ADD COLUMN invite_code VARCHAR(32) UNIQUE NOT NULL DEFAULT '';
    END IF;
END $$;

-- Patch users table
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='users' AND column_name='public_key'
    ) THEN
        ALTER TABLE users ADD COLUMN public_key TEXT NOT NULL DEFAULT '';
    END IF;
END $$;


CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- user table
CREATE TABLE IF NOT EXISTS users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    passhash   VARCHAR(255) NOT NULL,
    public_key TEXT        NOT NULL DEFAULT '',
    bio        TEXT        DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    creator_id  UUID         NOT NULL REFERENCES users(id),
    invite_code VARCHAR(32)  UNIQUE NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);
-- room users table
CREATE TABLE IF NOT EXISTS room_users (
    room_id   UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);
-- messages table
CREATE TABLE IF NOT EXISTS messages (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id         UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    sender_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_username VARCHAR(50) NOT NULL,
    content         TEXT        NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- room reads table
CREATE TABLE IF NOT EXISTS room_reads (
    room_id      UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);
-- patch messages table 
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='nonce'
    ) THEN
        ALTER TABLE messages ADD COLUMN nonce VARCHAR(64) NOT NULL DEFAULT '';
    END IF;
END $$;

-- create message_keys table
CREATE TABLE IF NOT EXISTS message_keys (
    message_id    UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    recipient_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    encrypted_key TEXT NOT NULL,
    PRIMARY KEY (message_id, recipient_id)
);
CREATE INDEX IF NOT EXISTS idx_message_keys_recipient ON message_keys(recipient_id);
CREATE INDEX IF NOT EXISTS idx_messages_room_id    ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
