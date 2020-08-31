create unique index player_already_on_auction on auctions (player_id) WHERE state not in ('ended', 'failed', 'cancelled', 'withadrable_by_seller', 'withadrable_by_buyer','validation');

CREATE FUNCTION update_auction_trigger() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.state = 'cancelled' AND OLD.state != 'started' THEN
        RAISE EXCEPTION 'error: cancel auction in state %', OLD.state;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_auction
    BEFORE UPDATE
    ON auctions
    FOR EACH ROW
    WHEN (OLD.state IS DISTINCT FROM NEW.state)
    EXECUTE PROCEDURE update_auction_trigger();

