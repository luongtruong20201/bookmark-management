CREATE TABLE users (
    id            VARCHAR(36),
    display_name  VARCHAR(255)  NOT NULL,
    username      VARCHAR(255)  NOT NULL,
    password      VARCHAR(2048) NOT NULL,
    email         VARCHAR(2048) NOT NULL,

    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT uni_username UNIQUE (username),
    CONSTRAINT uni_email    UNIQUE (email)
);
