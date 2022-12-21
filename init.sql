DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS bot_user;

CREATE TABLE bot_user
(
    id           INT8 NOT NULL PRIMARY KEY,
    user_name    VARCHAR(32),
    chat_id      INT8 NOT NULL,
    lang_code    VARCHAR(4),
    is_active    BOOL,
    created_at   TIMESTAMP DEFAULT now(),
    last_seen_at TIMESTAMP
);

CREATE TABLE location
(
    id        BIGSERIAL NOT NULL PRIMARY KEY,
    user_id   INT8      NOT NULL REFERENCES bot_user (id) ON DELETE CASCADE,
    latitude  FLOAT8,
    longitude FLOAT8,
    time_zone VARCHAR(30)
);

CREATE TABLE notification
(
    id                     BIGSERIAL NOT NULL PRIMARY KEY,
    user_id                INT8      NOT NULL REFERENCES bot_user (id) ON DELETE CASCADE,
    location_id            INT8      NOT NULL REFERENCES location (id) ON DELETE CASCADE,
    is_active              BOOL,
    notify_at              TIME,
    last_time_processed_at TIMESTAMP,
    UNIQUE (location_id)
);