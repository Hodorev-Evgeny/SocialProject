ALTER TABLE trackerapp.users
    ADD COLUMN role VARCHAR(32) NOT NULL DEFAULT 'passenger'
    CHECK (role IN ('passenger', 'driver', 'admin'));
