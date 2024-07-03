CREATE TABLE users (
    "id" char(20) NOT NULL,
    name varchar(32) NOT NULL,
    password char(60) NOT NULL,
    email varchar(255) NOT NULL,
    avatar varchar(255),
    role smallint NOT NULL,
    two_factor_auth varchar(64),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT USERS_PK_1 PRIMARY KEY ("id")
);

CREATE TABLE sessions (
    "id" char(20) NOT NULL,
    session char(128) NOT NULL,
    ip varchar(64) NOT NULL,
    device varchar(64) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT SESSIONS_PK_1 PRIMARY KEY ("id", session)
);

CREATE INDEX SESSIONS_FK_1 ON sessions ("id");