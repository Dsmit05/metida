-- Creation users table
CREATE TABLE users
(
    id         serial PRIMARY KEY,
    name       text,
    password   text NOT NULL,
    email      text NOT NULL,
    role       text NOT NULL DEFAULT 'User',
    is_deleted bool default false,
    verifay bool NOT NULL DEFAULT false
);

create unique index emails_index
    on users (email);

CREATE TABLE sessions
(
    id            serial PRIMARY KEY,
    user_email    text,
    refresh_token text,
    access_token  text,                                            -- give token for role services
    user_agent    text,
    ip            varchar(20),
    expires_in    bigint                   NOT NULL,
    created_at    timestamp with time zone NOT NULL DEFAULT now(), -- UTC
    FOREIGN KEY (user_email) REFERENCES users (email) ON DELETE SET NULL
);

create unique index refresh_token_index
    on sessions (refresh_token);


CREATE TABLE content
(
    id          serial PRIMARY KEY,
    user_email  text,
    name        text,
    description text,
    FOREIGN KEY (user_email) REFERENCES users (email) ON DELETE SET NULL
);

create unique index content_name_index
    on content (user_email, name);

CREATE TABLE blog
(
    id          serial PRIMARY KEY,
    name        text,
    description text
);

create unique index blog_name_index
    on blog (name);
