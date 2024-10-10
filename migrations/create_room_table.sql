CREATE TABLE rooms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    discription VARCHAR(50) NOT NULL,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    created_by VARCHAR(50)
)