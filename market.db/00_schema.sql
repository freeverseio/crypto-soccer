CREATE TABLE auction_states(
    state TEXT NOT NULL PRIMARY KEY
);
INSERT INTO auction_states(state) VALUES ('STARTED');
INSERT INTO auction_states(state) VALUES ('ASSET_FROZEN');
INSERT INTO auction_states(state) VALUES ('PAYING');
INSERT INTO auction_states(state) VALUES ('PAID');
INSERT INTO auction_states(state) VALUES ('NO_BIDS');
INSERT INTO auction_states(state) VALUES ('CANCELLED_BY_SELLER');
INSERT INTO auction_states(state) VALUES ('FREEZE_FAILED');
INSERT INTO auction_states(state) VALUES ('PAYMENT_FAILED');

CREATE TABLE auctions (
    uuid UUID NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price INT NOT NULL,
    rnd INT NOT NULL,
    valid_until TEXT NOT NULL,
    signature TEXT NOT NULL,
    state TEXT NOT NULL REFERENCES auction_states(state),
    paymet_link TEXT NOT NULL DEFAULT '',
    PRIMARY KEY(uuid)
);

CREATE TABLE bid_states(
    state TEXT NOT NULL PRIMARY KEY
);
INSERT INTO bid_states(state) VALUES ('ACCEPTED');
INSERT INTO bid_states(state) VALUES ('REFUSED');
INSERT INTO bid_states(state) VALUES ('PAYING');
INSERT INTO bid_states(state) VALUES ('PAID');
INSERT INTO bid_states(state) VALUES ('FAILED_TO_PAY');

CREATE TABLE bids (
    auction UUID NOT NULL REFERENCES auctions(uuid),
    extra_price NUMERIC(15,2) NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    state TEXT NOT NULL REFERENCES bid_states(state),
    paymet_link TEXT NOT NULL DEFAULT '',
    PRIMARY KEY(auction, extra_price)
);

