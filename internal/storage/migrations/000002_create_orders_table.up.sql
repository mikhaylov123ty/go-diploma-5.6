BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    number TEXT PRIMARY KEY NOT NULL,
    user_login TEXT NOT NULL,
    status TEXT NOT NULL,
    accrual double precision,
    uploaded_at timestamp

);

COMMIT;