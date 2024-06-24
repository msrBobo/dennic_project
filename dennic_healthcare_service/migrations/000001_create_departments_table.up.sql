
CREATE TABLE IF NOT EXISTS departments
(
    id          UUID PRIMARY KEY,
    department_order   SERIAL,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(200) NOT NULL,
    image_url   VARCHAR(200) NOT NULL,
    floor_number INT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP
    );

CREATE INDEX name_idx ON departments(lower(name));
