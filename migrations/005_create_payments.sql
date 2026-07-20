CREATE TABLE IF NOT EXISTS payments (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id               INTEGER NOT NULL,
    phone_number          TEXT NOT NULL,
    amount                INTEGER NOT NULL,
    account_reference     TEXT NOT NULL,
    checkout_request_id   TEXT,
    merchant_request_id   TEXT,
    mpesa_receipt_number  TEXT,
    status                TEXT NOT NULL DEFAULT 'pending',
    result_code           TEXT,
    result_desc           TEXT,
    created_at            TEXT NOT NULL,
    updated_at            TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);OA
