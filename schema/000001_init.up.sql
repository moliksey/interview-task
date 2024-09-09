CREATE TABLE "users"
(
    id            SERIAL PRIMARY KEY,
    email         varchar(255) NOT NULL,
    refresh_token varchar(255)
)
