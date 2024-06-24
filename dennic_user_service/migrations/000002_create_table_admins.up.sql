CREATE TYPE role_type AS ENUM('superadmin', 'admin');

/*admins table*/
CREATE TABLE IF NOT EXISTS admins (
    id UUID NOT NULL PRIMARY KEY,
    admin_order SERIAL,
    role role_type NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE NOT NULL,
    phone_number VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    gender gender_type NOT NULL,
    salary FLOAT NOT NULL DEFAULT 0.0, 
    biography TEXT NOT NULL,
    start_work_year DATE NOT NULL,
    end_work_year DATE,
    work_years INTEGER DEFAULT 0,
    refresh_token TEXT NOT NULL,
    image_url   VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX admin_unique_phone_number_deleted_at_null_idx ON admins(phone_number) WHERE deleted_at IS NULL; --unique admin phone number, users with non-soft-deleted accounts cannot register again with the same phone number.
CREATE INDEX deleted_at_idx2 ON admins(deleted_at);
