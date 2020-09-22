ALTER TABLE bids DROP CONSTRAINT bids_auction_id_fkey;
ALTER TABLE bids add CONSTRAINT bids_auction_id_fkey FOREIGN KEY (auction_id) REFERENCES auctions(id) DEFERRABLE;