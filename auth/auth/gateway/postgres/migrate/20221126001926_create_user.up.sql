CREATE TABLE auth_user(
    id SERIAL PRIMARY KEY,

    username TEXT NOT NULL UNIQUE,
    secret   TEXT NOT NULL
);
