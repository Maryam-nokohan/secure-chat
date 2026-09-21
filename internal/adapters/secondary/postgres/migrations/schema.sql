-- ============================================================
-- Secure Chat schema
--
-- This file is executed on every startup (see migrations.go) and must
-- be safe to re-run against both a brand-new database and an existing
-- one. Rules that keep it that way:
--   1. Objects are created in dependency order: extensions -> tables
--      (parents before the tables that reference them) -> indexes ->
--      triggers.
--   2. Every CREATE uses IF NOT EXISTS.
--   3. Tables are created with their full, final column set instead of
--      being patched together with ALTER TABLE across many runs, so
--      there's nothing that depends on a column being added "later in
--      the file" by a previous version of this script.
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ------------------------------------------------------------
-- users
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id                   UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username             VARCHAR(50)  UNIQUE NOT NULL,
    passhash             VARCHAR(255),
    public_key           TEXT         NOT NULL DEFAULT '',
    bio                  TEXT         DEFAULT '',

    email                VARCHAR(255) DEFAULT '',
    provider             VARCHAR(20)  DEFAULT '',
    provider_id          VARCHAR(255) DEFAULT '',

    role                 VARCHAR(20)  NOT NULL DEFAULT 'user',
    deleted_at           TIMESTAMPTZ,

    wrapped_private_key  TEXT         NOT NULL DEFAULT '',
    private_key_iv       VARCHAR(64)  NOT NULL DEFAULT '',
    private_key_salt     VARCHAR(64)  NOT NULL DEFAULT '',

    public_id            VARCHAR(12),
    avatar_path          TEXT         NOT NULL DEFAULT '',

    created_at           TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

-- passhash is nullable because OAuth-only accounts never set one.
ALTER TABLE users ALTER COLUMN passhash DROP NOT NULL;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ------------------------------------------------------------
-- rooms (created before contacts, which references it)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS rooms (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    creator_id  UUID         NOT NULL REFERENCES users(id),
    invite_code VARCHAR(32)  UNIQUE NOT NULL DEFAULT '',
    is_direct   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------
-- contacts (references users + rooms, both already exist above)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS contacts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    room_id       UUID REFERENCES rooms(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT no_self_contact CHECK (requester_id <> addressee_id)
);

-- ------------------------------------------------------------
-- room_users
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS room_users (
    room_id   UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);

-- ------------------------------------------------------------
-- messages
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS messages (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id         UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    sender_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_username VARCHAR(50) NOT NULL,
    content         TEXT        NOT NULL,
    nonce           VARCHAR(64) NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------
-- room_reads
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS room_reads (
    room_id      UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);

-- ------------------------------------------------------------
-- message_keys (per-recipient wrapped AES key for a message)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS message_keys (
    message_id    UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    recipient_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    encrypted_key TEXT NOT NULL,
    PRIMARY KEY (message_id, recipient_id)
);

-- ------------------------------------------------------------
-- attachments
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS attachments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id      UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    uploader_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    storage_key  TEXT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes   BIGINT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- Data normalization (safe no-ops once already normalized)
-- ============================================================
UPDATE users SET email = lower(trim(email)) WHERE email <> '';

-- ============================================================
-- Indexes & constraints (all referenced columns exist by this point)
-- ============================================================

-- users
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_public_id
    ON users (public_id) WHERE public_id IS NOT NULL AND public_id <> '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique
    ON users (LOWER(TRIM(email)))
    WHERE email IS NOT NULL AND TRIM(email) <> '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_identity
    ON users (provider, provider_id) WHERE provider_id <> '';

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_role       ON users(role);

-- contacts
CREATE UNIQUE INDEX IF NOT EXISTS idx_contacts_pair ON contacts (
    LEAST(requester_id, addressee_id), GREATEST(requester_id, addressee_id)
);
CREATE INDEX IF NOT EXISTS idx_contacts_addressee ON contacts(addressee_id, status);
CREATE INDEX IF NOT EXISTS idx_contacts_requester ON contacts(requester_id, status);

-- attachments / messages / message_keys
CREATE INDEX IF NOT EXISTS idx_attachments_room       ON attachments(room_id);
CREATE INDEX IF NOT EXISTS idx_message_keys_recipient ON message_keys(recipient_id);
CREATE INDEX IF NOT EXISTS idx_messages_room_id        ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at     ON messages(created_at);