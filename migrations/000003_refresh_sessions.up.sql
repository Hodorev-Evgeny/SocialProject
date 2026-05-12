CREATE TABLE trackerapp.refresh_sessions (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES trackerapp.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX refresh_sessions_token_hash_idx ON trackerapp.refresh_sessions (token_hash);
CREATE INDEX refresh_sessions_user_id_idx ON trackerapp.refresh_sessions (user_id);
