CREATE TABLE trackerapp.revoked_access_tokens (
    jti UUID PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX revoked_access_tokens_expires_at_idx ON trackerapp.revoked_access_tokens (expires_at);
