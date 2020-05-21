CREATE TYPE auction_state AS ENUM ('started', 'failed', 'cancelled', 'ended', 'asset_frozen', 'paying', 'withadrable_by_seller', 'withadrable_by_buyer','validation');
CREATE TABLE auctions (
    id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price BIGINT NOT NULL,
    rnd BIGINT NOT NULL,
    valid_until BIGINT NOT NULL,
    signature TEXT NOT NULL,
    state auction_state NOT NULL,
    state_extra TEXT NOT NULL,
    payment_url TEXT NOT NULL,
    seller TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TYPE bid_state AS ENUM ('accepted','paying','paid','failed');
CREATE TABLE bids (
    auction_id TEXT NOT NULL REFERENCES auctions(id),
    extra_price INT NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state bid_state NOT NULL,
    state_extra TEXT NOT NULL,
    payment_id TEXT NOT NULL,
    payment_url TEXT NOT NULL,
    payment_deadline TEXT NOT NULL,
    PRIMARY KEY(auction_id, extra_price)
);

CREATE TABLE shop_items (
    uuid UUID NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    type INT NOT NULL,
    quantity INT NOT NULL,
    price INT NOT NULL,
    PRIMARY KEY(uuid)
);

CREATE TYPE playstore_order_state AS ENUM ('open','acknowledged', 'complete', 'refunding', 'refunded', 'failed');
CREATE TABLE playstore_orders(
    order_id TEXT NOT NULL,
    package_name TEXT NOT NULL,
    product_id TEXT NOT NULL,
    purchase_token TEXT NOT NULL,
    player_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state playstore_order_state NOT NULL,
    state_extra TEXT NOT NULL,
    PRIMARY KEY(order_id)
)

CREATE TABLE playstore_orders_histories(
    order_id TEXT NOT NULL REFERENCES playstore_orders(order_id),
    package_name TEXT NOT NULL,
    product_id TEXT NOT NULL,
    purchase_token TEXT NOT NULL,
    player_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state playstore_order_state NOT NULL,
    state_extra TEXT NOT NULL,
)

