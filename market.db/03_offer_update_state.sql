CREATE FUNCTION update_offer_trigger() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.state = 'cancelled' AND OLD.state != 'started' THEN
        RAISE EXCEPTION 'error: cancel offer in state %', OLD.state;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_offer
    BEFORE UPDATE
    ON offers_v2
    FOR EACH ROW
    WHEN (OLD.state IS DISTINCT FROM NEW.state)
    EXECUTE PROCEDURE update_offer_trigger();


