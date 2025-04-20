CREATE TABLE if not exists tasks
(
    id           TEXT PRIMARY KEY not null,
    status       varchar(255)     NOT NULL,
    payload      text             not null,
    result       text not null    DEFAULT '',
    "created_at" timestamptz             NOT NULL DEFAULT (now())
);

