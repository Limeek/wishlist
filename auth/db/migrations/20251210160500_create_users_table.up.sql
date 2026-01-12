CREATE TABLE users
(
    id         uuid PRIMARY KEY,
    email      text,
    password   text,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

CREATE TABLE refresh_tokens
(
    user_id    uuid PRIMARY KEY references users(id),
    token      text,
    expires_at timestamptz
);