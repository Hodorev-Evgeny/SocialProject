CREATE SCHEMA trackerapp;


CREATE TABLE trackerapp.users (
    id UUID PRIMARY KEY,
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

CREATE TABLE trackerapp.transactions (
    id UUID PRIMARY KEY,
    sum INTEGER NOT NULL,
    type_transaction VARCHAR(125) NOT NULL CHECK (
        type_transaction = 'Income' OR type_transaction = 'Expenditure'
        ),
    date TIMESTAMP NOT NULL,
    category VARCHAR(125) NOT NULL,
    comments VARCHAR(1000),
    user_id UUID NOT NULL,
    time_create TIMESTAMPTZ NOT NULL,
    time_changes TIMESTAMPTZ,

    FOREIGN KEY (user_id) REFERENCES trackerapp.users(id) ON DELETE CASCADE
);