CREATE TYPE auction_pass_playstore_order_state AS ENUM ('open','acknowledged', 'complete', 'refunding', 'refunded', 'failed');
CREATE TABLE auction_pass_playstore_orders(
    order_id TEXT NOT NULL,
    package_name TEXT NOT NULL,
    product_id TEXT NOT NULL,
    purchase_token TEXT NOT NULL,
    team_id TEXT NOT NULL,
    owner TEXT NOT NULL,
    signature TEXT NOT NULL,
    state playstore_order_state NOT NULL,
    state_extra TEXT NOT NULL,
    PRIMARY KEY(purchase_token)
);

CREATE TABLE auction_pass_playstore_orders_histories(
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_id TEXT NOT NULL,
    package_name TEXT NOT NULL,
    product_id TEXT NOT NULL,
    purchase_token TEXT NOT NULL REFERENCES auction_pass_playstore_orders(purchase_token),
    team_id TEXT NOT NULL,
    owner TEXT NOT NULL,
    signature TEXT NOT NULL,
    state playstore_order_state NOT NULL,
    state_extra TEXT NOT NULL
);

CREATE TABLE auction_pass(
    owner TEXT NOT NULL,
    PRIMARY KEY(owner)
);

