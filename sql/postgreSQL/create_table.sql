-------------------------------------------------------------
CREATE TABLE wallet_tab (
    user_id INTEGER PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0,
    frozen_balance INTEGER NOT NULL DEFAULT 0,
    ext_info VARCHAR(50) NOT NULL,
    created_time BIGINT NOT NULL ,
    updated_time BIGINT NOT NULL
);
-------------------------------------------------------------
CREATE TABLE transaction_tab (
    transaction_id BIGINT PRIMARY KEY,
    order_type INTEGER NOT NULL, -- deposit:0, withdraw:1, transfer2
    transaction_type INTEGER NOT NULL, -- money_in:0, money_out:1
    amount INTEGER NOT NULL,
    status INTEGER DEFAULT 0, -- pending:0, complete:1, failed:2, canceled 3;
    user_id INTEGER NOT NULL,
    oppo_user_id INTEGER NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL,
    update_time BIGINT NOT NULL,
    last_process_time BIGINT NOT NULL
);

CREATE INDEX idx_user_id_process_time ON transaction_tab(user_id, last_process_time);
-------------------------------------------------------------