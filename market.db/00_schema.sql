CREATE TYPE auction_state AS ENUM ('started', 'failed', 'cancelled', 'ended', 'asset_frozen', 'paying', 'withadrable_by_seller', 'withadrable_by_buyer');
CREATE TABLE auctions (
    id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price BIGINT NOT NULL,
    rnd BIGINT NOT NULL,
    valid_until BIGINT NOT NULL,
    signature TEXT NOT NULL,
    state auction_state NOT NULL,
    state_extra TEXT NOT NULL DEFAULT '',
    payment_url TEXT NOT NULL DEFAULT '',
    seller TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE bid_states(
    state TEXT NOT NULL PRIMARY KEY
);
INSERT INTO bid_states(state) VALUES ('ACCEPTED');
INSERT INTO bid_states(state) VALUES ('REFUSED');
INSERT INTO bid_states(state) VALUES ('PAYING');
INSERT INTO bid_states(state) VALUES ('PAID');
INSERT INTO bid_states(state) VALUES ('FAILED');

CREATE TABLE bids (
    auction_id TEXT NOT NULL REFERENCES auctions(id),
    extra_price INT NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state TEXT NOT NULL REFERENCES bid_states(state),
    state_extra TEXT NOT NULL DEFAULT '',
    payment_id TEXT NOT NULL DEFAULT '',
    payment_url TEXT NOT NULL DEFAULT '',
    payment_deadline TEXT NOT NULL DEFAULT '0',
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

