create unique index player_already_on_auction on auctions (player_id) WHERE state not in ('ended', 'failed', 'cancelled', 'withadrable_by_seller', 'withadrable_by_buyer','validation');
