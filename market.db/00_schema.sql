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
CREATE INDEX idx_auctions_player_id ON auctions (player_id);

CREATE TABLE auctions_histories (
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id TEXT NOT NULL REFERENCES auctions(id),
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price BIGINT NOT NULL,
    rnd BIGINT NOT NULL,
    valid_until BIGINT NOT NULL,
    signature TEXT NOT NULL,
    state auction_state NOT NULL,
    state_extra TEXT NOT NULL,
    payment_url TEXT NOT NULL,
    seller TEXT NOT NULL
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

CREATE TABLE bids_histories (
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    auction_id TEXT NOT NULL, 
    extra_price INT NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state bid_state NOT NULL,
    state_extra TEXT NOT NULL,
    payment_id TEXT NOT NULL,
    payment_url TEXT NOT NULL,
    payment_deadline TEXT NOT NULL,
    FOREIGN KEY (auction_id, extra_price) REFERENCES bids (auction_id, extra_price)
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
    PRIMARY KEY(purchase_token)
);

CREATE TABLE playstore_orders_histories(
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_id TEXT NOT NULL,
    package_name TEXT NOT NULL,
    product_id TEXT NOT NULL,
    purchase_token TEXT NOT NULL REFERENCES playstore_orders(purchase_token),
    player_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state playstore_order_state NOT NULL,
    state_extra TEXT NOT NULL
);

CREATE TYPE offer_state AS ENUM ('started', 'failed', 'cancelled', 'ended', 'accepted');
CREATE TABLE offers (
    id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price BIGINT NOT NULL,
    rnd BIGINT NOT NULL,
    valid_until BIGINT NOT NULL,
    signature TEXT NOT NULL,
    state offer_state NOT NULL,
    state_extra TEXT NOT NULL,
    seller TEXT NOT NULL,
    buyer TEXT NOT NULL,
    auction_id TEXT REFERENCES auctions(id),
    buyer_team_id TEXT NOT NULL,
    PRIMARY KEY(id)
);
CREATE INDEX idx_offers_player_id ON offers (player_id);
CREATE INDEX idx_offers_auction_id ON offers (auction_id);

CREATE TABLE offers_histories (
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id TEXT NOT NULL REFERENCES offers(id),
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price BIGINT NOT NULL,
    rnd BIGINT NOT NULL,
    valid_until BIGINT NOT NULL,
    signature TEXT NOT NULL,
    state offer_state NOT NULL,
    state_extra TEXT NOT NULL,
    seller TEXT NOT NULL,
    buyer TEXT NOT NULL,
    auction_id TEXT REFERENCES auctions(id),
    buyer_team_id TEXT NOT NULL
);
