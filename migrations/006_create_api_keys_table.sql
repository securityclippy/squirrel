-- +goose Up
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    key_prefix VARCHAR(20) NOT NULL,
    permissions JSONB NOT NULL DEFAULT '["read"]'::jsonb,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for efficient lookups
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_key_prefix ON api_keys(key_prefix);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);

-- Create a default user for existing reminders
INSERT INTO users (auth0_id, email, name, email_verified)
VALUES ('migration-user', 'migration@example.com', 'Migration User', false)
ON CONFLICT (auth0_id) DO NOTHING;

-- Add new user_id column as nullable first
ALTER TABLE reminders ADD COLUMN user_id_new INTEGER REFERENCES users(id) ON DELETE CASCADE;

-- Update existing reminders to reference the migration user
UPDATE reminders SET user_id_new = (
    SELECT id FROM users WHERE auth0_id = 'migration-user' LIMIT 1
);

-- Make the column NOT NULL and drop the old column
ALTER TABLE reminders ALTER COLUMN user_id_new SET NOT NULL;
ALTER TABLE reminders DROP COLUMN user_id;
ALTER TABLE reminders RENAME COLUMN user_id_new TO user_id;

-- Create index for efficient user-based reminder queries
CREATE INDEX idx_reminders_user_id ON reminders(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_reminders_user_id;
ALTER TABLE reminders 
DROP COLUMN user_id,
ADD COLUMN user_id VARCHAR(255) NOT NULL DEFAULT 'user123';

DROP INDEX IF EXISTS idx_api_keys_expires_at;
DROP INDEX IF EXISTS idx_api_keys_is_active;
DROP INDEX IF EXISTS idx_api_keys_key_prefix;
DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP TABLE IF EXISTS api_keys;