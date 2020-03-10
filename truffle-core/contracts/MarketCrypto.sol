pragma solidity >=0.5.12 <=0.6.3;

import "./Proxy.sol";
import "./Market.sol";

/**
 * @title Auctions operated in cryptocurrency, without anyone's permission (other than sellers and buyers)
 */
 
contract MarketCrypto {

    uint32 internal _auctionDuration = 24 hours; 
    uint256 internal _minimumBidIncrement = 0.5 ether; // bid for at least this amount of XDAI, or increase previous by this amount
    uint256 constant private MAX_128_BIT = (2**128)-1; 
    address constant private NULL_ADDR = address(0x0);
    
    Market private _market;

    mapping (uint256 => uint256) internal _playerIdToAuctionId;
    uint256 _nAuctions;
    // an auction exists if id < nAcutions, or if startingPrice > 0
    // an auction had at least one bid if highestBid > 0
    mapping (uint256 => uint256) private _startingPrice;
    mapping (uint256 => uint256) private  _highestBid;
    mapping (uint256 => uint256) private  _validUntil;
    mapping (uint256 => address) private  _seller;
    mapping (uint256 => address) private  _highestBidder;
    mapping (uint256 => uint256) private  _teamIdHighestBidder;
    mapping (uint256 => bool) private  _assetWentToNewOwner;
    mapping (uint256 => mapping(address => uint256)) private _balance;
    
    function setMarketAddress(address proxyAddr) external {
        _market = Market(proxyAddr);
    }
    
    function setActionDuration(uint32 newDuration) external {
        _auctionDuration = newDuration;
    }

    function setMinimumBid(uint128 newMinimum) external {
        _minimumBidIncrement = newMinimum;
    }

    // an asset is ready to be put for sale again if it either has never been assigned to an auction or, in case it has:
    //  - if the auction never received a bid
    //  - if the auction received bids, and the new owner finally executed the "receive asset" function.
    // Note that the difference with isPlayerFrozenCrypto is that isAuctionSettled is a bit more restrictive, and triggered
    // as soon as asset is putForSale. In contrast, isPlayerFrozenCrypto on triggers when first bid arrives.
    function isAuctionSettled(uint256 auctionId) public view returns(bool) {
        if (auctionId == 0) return true;
        if (_assetWentToNewOwner[auctionId]) return true;
        // check if it was put for sale but experied without bids
        return (_highestBid[auctionId] == 0) && (_validUntil[auctionId] < now);
    } 

    // this will check also for local crypto freeze of team when feature is added
    function isTeamFrozenInAnyMarket(uint256 teamId) public view returns (bool) {
        return _market.isTeamFrozen(teamId);
    }

    function putPlayerForSale(uint256 playerId, uint256 startingPrice) external {
        // TODO: cheaper if we return the 4 needed data in just 1 call
        uint256 currentTeamId  = _market.getCurrentTeamIdFromPlayerId(playerId);
        address currentOwner   = _market.getOwnerTeam(currentTeamId);
        uint256 prevAuctionId  = _playerIdToAuctionId[playerId];
        bool OK = (
            // check asset is not already for sale, or if data exists, that the auction has been fully settled
            isAuctionSettled(prevAuctionId) &&
            // check player is not already frozen
            !_market.isPlayerFrozenInAnyMarket(playerId) &&  
            // check that the team it belongs to is not already frozen
            !_market.isTeamFrozen(currentTeamId) &&
            // check asset is owned by legit address
            (currentOwner != NULL_ADDR) && 
            // check asset is owned by sender of this TX
            (currentOwner == msg.sender)   
        );
        require(OK, "conditions to putPlayerForSale not met");
        _nAuctions++;
        uint256 auctionId = _nAuctions;
        _playerIdToAuctionId[playerId] = auctionId;
        _startingPrice[auctionId] = startingPrice;
        _validUntil[auctionId] = now + _auctionDuration;
        _seller[auctionId] = currentOwner;
    }

    
    // TODO: encode isForeignMarket so that we cannot complete the auction ourselves
    function bidForPlayer(uint256 playerId, uint256 bidderTeamId) external payable {
        // TODO: save gas by calling 1 once and returning all data in 1 call
        require(msg.sender != NULL_ADDR, "sender cannot be the null address");
        require(msg.sender != _seller[playerId], "seller is not allowed to bid for its own assets");
        uint256 auctionId = _playerIdToAuctionId[playerId];
        require(auctionId != 0, "player has not been put for sale yet");
        require(now < _validUntil[auctionId], "too late to bid, auction time has expired");
        require(msg.sender == _market.getOwnerTeam(bidderTeamId), "only the owner of the team can bid for a player");

        uint256 bidAmount = _balance[auctionId][msg.sender] + msg.value;

        // if this is the first bid, freeze the asset from every market
        if (!_market.isPlayerFrozenInAnyMarket(playerId)) {
            require (bidAmount >= _startingPrice[auctionId], "bid did not increment the previous bid above the minimum allowed");
            uint256 currentTeamId  = _market.getCurrentTeamIdFromPlayerId(playerId);
            require (!isTeamFrozenInAnyMarket(currentTeamId), "the team that this player belongs to is already frozen. Cannot sale players if team is for sale");
            _market.setIsPlayerFrozenCrypto(playerId, true);
        } else {
            require (bidAmount >= _highestBid[auctionId] + _minimumBidIncrement, "bid did not increment the previous bid above the minimum allowed");
        }

        _balance[auctionId][msg.sender] = bidAmount;
        _highestBid[auctionId] = bidAmount;
        _highestBidder[auctionId] = msg.sender;
        _teamIdHighestBidder[auctionId] = bidderTeamId;
    }

    // Note that this function can be executed by anyone. Of course, only the highest bidder is expected to 
    // be the interest party in executing it, but we allow anyone to operate it, since this was the intention of the auction.
    function withdraw(uint256 auctionId) external {
        require(now > _validUntil[auctionId], "highest bid goes to seller, so highest bidder cannot withdraw");
        require(msg.sender != _highestBidder[auctionId], "highest bid deposit goes to seller, so highest bidder cannot withdraw");
        uint256 amount;
        if (msg.sender == _seller[auctionId]) {
            amount = _highestBid[auctionId];
            _highestBid[auctionId] = 0;
        } else {
            amount = _balance[auctionId][msg.sender];
            _balance[auctionId][msg.sender] = 0;
        }
        require(amount > 0, "nothing to withdraw for this msg sender");
        require(msg.sender.send(amount), "failure when withdrawing legit funds");
    }

    // Note that this function can be executed by anyone. Of course, only the highest bidder is expected to 
    // be the interest party in executing it, but we allow anyone to operate it, since this was the intention of the auction.
    function executePlayerTransfer(uint256 playerId) external {
        uint256 auctionId = _playerIdToAuctionId[playerId];
        require(!_assetWentToNewOwner[auctionId], "the player in this auction was already transferred");
        _assetWentToNewOwner[auctionId] = true;
        _market.setIsPlayerFrozenCrypto(playerId, false);
        _market.transferPlayer(playerId, _teamIdHighestBidder[auctionId]);
        _market.setIsPlayerFrozenCrypto(playerId, false);
    }
}
