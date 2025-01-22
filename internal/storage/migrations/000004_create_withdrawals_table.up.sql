BEGIN;

CREATE TABLE IF NOT EXISTS withdrawals(
    order_id TEXT PRIMARY KEY NOT NULL,
    user_login TEXT NOT NULL,
    sum double precision NOT NULL,
    processed_at timestamp
);

COMMIT;