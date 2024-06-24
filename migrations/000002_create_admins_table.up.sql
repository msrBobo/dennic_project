CREATE TYPE role_type AS ENUM('admin');

-- /*admins table*/
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
    image_url VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '5 hours'),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX admin_unique_phone_number_deleted_at_null_idx ON admins(phone_number) WHERE deleted_at IS NULL; --unique admin phone number, users with non-soft-deleted accounts cannot register again with the same phone number.
CREATE INDEX deleted_at_idx2 ON admins(deleted_at);

---------------------------------------------------------------------------------------------------------------------------

-- For admins table
INSERT INTO admins (id, admin_order, role, first_name, last_name, birth_date, phone_number, email, password, gender, salary, biography, start_work_year, refresh_token, created_at, image_url)
VALUES
  ('123e4567-e89b-12d3-a456-426614174001', 1, 'admin',    'Mubina', 'Yigitaliyeva', '2007-01-14', '+998972645292', 'mubinayigitaliyeva00@example.com', 'mubina2007', 'female', 1000.00, 'Mubina bio', '2024-01-01', 'refresh_token1', CURRENT_TIMESTAMP, 'https://example.com/admin1.jpg'),
  ('123e4567-e89b-12d3-a456-426614174002', 1, 'admin',    'Bobomurod', 'Artigaliyev', '2000-08-30', '+9998994767316', 'msrbobo@example.com', 'BoBo1234', 'male', 1000.00, 'Bobomurod bio', '2024-01-01', 'refresh_token1', CURRENT_TIMESTAMP, 'https://example.com/admin1.jpg'),
  ('123e4567-e89b-12d3-a456-426614174003', 10, 'admin', 'Toraxon', 'Joraxonov', '2003-05-05', '+998950230605', 'torakhonoffical@gmail.com', 'Toraxon', 'male', 10000.00, 'Toraxon bio', '2024-01-01', 'refresh_token10', CURRENT_TIMESTAMP, 'https://example.com/superadmin10.jpg');;