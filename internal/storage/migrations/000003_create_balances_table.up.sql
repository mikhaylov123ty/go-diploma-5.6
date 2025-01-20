BEGIN;

CREATE TABLE IF NOT EXISTS balances(
    user_login TEXT PRIMARY KEY NOT NULL,
    current double precision NOT NULL,
    withdrawn double precision NOT NULL
);

COMMIT;