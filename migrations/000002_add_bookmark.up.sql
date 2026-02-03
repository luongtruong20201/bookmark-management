CREATE TABLE bookmarks (
    id          VARCHAR(36) UNIQUE,
    description VARCHAR(255),
    url         VARCHAR(2048) NOT NULL,
    code        VARCHAR(10)   NOT NULL,
    user_id     VARCHAR(36)   NOT NULL,

    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT bookmarks_pkey PRIMARY KEY (id),
    CONSTRAINT uni_bookmark_code UNIQUE (code),
    CONSTRAINT fk_bookmarks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
