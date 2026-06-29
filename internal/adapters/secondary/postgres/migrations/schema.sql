DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='rooms' AND column_name='invite_code'
    ) THEN
        ALTER TABLE rooms ADD COLUMN invite_code VARCHAR(32) UNIQUE NOT NULL DEFAULT '';
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='users' AND column_name='public_key'
    ) THEN
        ALTER TABLE users ADD COLUMN public_key TEXT NOT NULL DEFAULT '';
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='text'
    ) THEN
        ALTER TABLE messages RENAME COLUMN text TO encrypted_content;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='encrypted_content'
    ) THEN
        ALTER TABLE messages ADD COLUMN encrypted_content TEXT NOT NULL DEFAULT '';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='encrypted_aes_key'
    ) THEN
        ALTER TABLE messages ADD COLUMN encrypted_aes_key TEXT NOT NULL DEFAULT '';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='iv'
    ) THEN
        ALTER TABLE messages ADD COLUMN iv TEXT NOT NULL DEFAULT '';
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='receiver_id'
    ) THEN
        ALTER TABLE messages RENAME COLUMN receiver_id TO recipient_id;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='messages' AND column_name='recipient_id'
    ) THEN
        ALTER TABLE messages ADD COLUMN recipient_id UUID REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50)  UNIQUE NOT NULL,
    passhash    VARCHAR(255) NOT NULL,
    public_key  TEXT         NOT NULL,     
    bio         TEXT         DEFAULT '',
    created_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS rooms (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) UNIQUE NOT NULL,
    creator_id UUID         NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS room_users (
    room_id   UUID REFERENCES rooms(id)  ON DELETE CASCADE,
    user_id   UUID REFERENCES users(id)  ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id                UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id           UUID  NOT NULL REFERENCES rooms(id)  ON DELETE CASCADE,
    sender_id         UUID  NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    recipient_id      UUID  NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    encrypted_content TEXT  NOT NULL,
    encrypted_aes_key TEXT  NOT NULL,
    iv                TEXT  NOT NULL,
    created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_room_id   ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
