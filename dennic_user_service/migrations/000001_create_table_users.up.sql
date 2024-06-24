CREATE TYPE gender_type AS ENUM('male', 'female');

/*users table*/
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    user_order SERIAL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE NOT NULL,
    phone_number VARCHAR(100) NOT NULL,
    password VARCHAR(50) NOT NULL,
    gender gender_type NOT NULL,
    refresh_token TEXT NOT NULL,
    image_url   VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX unique_phone_number_deleted_at_null_idx ON users(phone_number) WHERE deleted_at IS NULL; --unique phone number, users with non-soft-deleted accounts cannot register again with the same phone number.
CREATE INDEX deleted_at_idx ON users(deleted_at);