CREATE SCHEMA social;


CREATE TABLE social.users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(15) CHECK (
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
        ),
    password VARCHAR NOT NULL,
    time_add TIMESTAMPTZ NOT NULL
);

CREATE TABLE social.drive (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    time_to TIMESTAMPTZ NOT NULL,
    time_from TIMESTAMPTZ,
    status VARCHAR(128) CHECK (
        status IN ('Active', 'Finished')
        ),

    FOREIGN KEY (user_id) REFERENCES social.users(id) ON DELETE CASCADE
);