DROP TABLE IF EXISTS offers_histories;
DROP TABLE IF EXISTS offers;

CREATE TABLE offers (
    auction_id TEXT NOT NULL,
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
    buyer_team_id TEXT NOT NULL,
    PRIMARY KEY(auction_id)
);
CREATE INDEX idx_offers_player_id ON offers (player_id);
CREATE INDEX idx_offers_auction_id ON offers (auction_id);


CREATE TABLE offers_histories (
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    auction_id TEXT NOT NULL REFERENCES offers(auction_id),
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
    buyer_team_id TEXT NOT NULL
);