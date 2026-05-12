CREATE TABLE trackerapp.register_verifications (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES trackerapp.users(id) ON DELETE CASCADE,
    purpose VARCHAR(32) NOT NULL,
    code_hash VARCHAR NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX register_verifications_user_id_idx ON trackerapp.register_verifications (user_id);
