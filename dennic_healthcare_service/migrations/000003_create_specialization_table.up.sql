CREATE TABLE IF NOT EXISTS specializations
(
    id            UUID PRIMARY KEY,
    specializations_order SERIAL,
    name          VARCHAR(100) NOT NULL,
    description   VARCHAR(200) NOT NULL,
    department_id UUID REFERENCES departments(id) ON DELETE CASCADE NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP,
    deleted_at    TIMESTAMP
);

CREATE INDEX specialization_name_idx ON specializations(name);
