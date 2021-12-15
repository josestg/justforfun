-- up script here...
CREATE TABLE IF NOT EXISTS users
(
    id            UUID         NOT NULL PRIMARY KEY,
    name          VARCHAR(40)  NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    date_created  TIMESTAMP,
    date_updated  TIMESTAMP,

    CONSTRAINT users__unique_email_each_users UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS email_index ON users (email);

---+split+---

-- down script here...
DROP INDEX IF EXISTS users.email_index;
DROP TABLE IF EXISTS users;
