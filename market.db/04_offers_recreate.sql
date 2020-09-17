ALTER TABLE IF EXISTS offers_histories RENAME TO old_offers_histories;
ALTER TABLE IF EXISTS offers RENAME TO old_offers;
ALTER INDEX IF EXISTS idx_offers_player_id RENAME TO idx_old_offers_player_id;
ALTER INDEX IF EXISTS idx_offers_auction_id RENAME TO idx_old_offers_auction_id;

comment on table old_offers_histories is
  E'@omit';
comment on table old_offers is
  E'@omit';

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

CREATE TRIGGER update_offer
    BEFORE UPDATE
    ON offers
    FOR EACH ROW
    WHEN (OLD.state IS DISTINCT FROM NEW.state)
    EXECUTE PROCEDURE update_offer_trigger();


ALTER TABLE auctions ADD COLUMN offer_valid_until BIGINT NOT NULL;
ALTER TABLE auctions_histories ADD COLUMN offer_valid_until BIGINT NOT NULL;