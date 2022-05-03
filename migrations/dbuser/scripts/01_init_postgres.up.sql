CREATE TABLE users (
    id uuid PRIMARY KEY,
    username character varying NOT NULL,
    email character varying NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);