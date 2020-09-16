CREATE FUNCTION is_higher_offer() RETURNS TRIGGER AS $$
DECLARE
    numHigherOffers integer;
BEGIN
    SELECT COUNT(auction_id)
    INTO numHigherOffers
    FROM offers_v2
    WHERE state='started' AND player_id=NEW.player_id AND price >= NEW.price;

    IF numHigherOffers > 0 THEN
        RAISE EXCEPTION 'error: Price not the highest';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_offers_v2
    BEFORE INSERT
    ON offers_v2
    FOR EACH ROW
    EXECUTE PROCEDURE is_higher_offer();