pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract Market {
    event PlayerFreeze(uint256 playerId, uint256 auctionData, bool frozen);
    event TeamFreeze(uint256 teamId, uint256 auctionData, bool frozen);

    uint8 constant internal IDX_MSG = 0;
    uint8 constant internal IDX_r   = 1;
    uint8 constant internal IDX_s   = 2;
    uint8 constant internal PUT_FOR_SALE  = 1;
    uint8 constant internal MAKE_AN_OFFER = 2;
    // POST_AUCTION_TIME: is how long does the buyer have to pay in fiat, after auction is finished.
    //  ...it includes time to ask for a 2nd-best bidder, or 3rd-best.
    uint256 constant public POST_AUCTION_TIME   = 6 hours; 
    uint256 constant public AUCTION_TIME        = 24 hours; 
    uint256 constant public MAX_VALID_UNTIL     = 30 hours; // the sum of the previous two
    uint256 constant private VALID_UNTIL_MASK   = 0x3FFFFFFFF; // 2^34-1 (34 bit)
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0

    Assets private _assets;

    mapping (uint256 => uint256) private playerIdToAuctionData;
    mapping (uint256 => uint256) private teamIdToAuctionData;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }
    
    // Main PLAYER auction functions: freeze & complete
    function freezePlayer(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32[3] memory sig,
        uint8 sigV
    ) public {
        require(areFreezePlayerRequirementsOK(sellerHiddenPrice, validUntil, playerId, sig, sigV), "FreePlayer requirements not met");
        // // Freeze player
        playerIdToAuctionData[playerId] = validUntil + uint256(sellerHiddenPrice << 34);
        emit PlayerFreeze(playerId, playerIdToAuctionData[playerId], true);
    }

    function freezePromoPlayer(
        uint256 playerId, 
        uint256 validUntil
    ) public {
        require(msg.sender == _assets.rosterAddr() , "Only the Roster can create promo players");
        require(!isPlayerFrozen(playerId));
        require(_assets.teamExists(_assets.getTargetTeamId(playerId)), "cannot offer a promo player to a non-existent team");
        require(!_assets.isBotTeam(_assets.getTargetTeamId(playerId)), "cannot offer a promo player to a bot team");
        require(validUntil > now, "validUntil is in the past");
        require(validUntil < now + MAX_VALID_UNTIL, "validUntil is too large");
        playerIdToAuctionData[playerId] = validUntil;
        emit PlayerFreeze(playerId, playerIdToAuctionData[playerId], true);
    }

    function completePromoPlayerTransfer(
        uint256 playerId,
        uint256 validUntil,
        bytes32[3] memory sig,
        uint8 sigV
     ) public {
        uint256 buyerTeamId = _assets.getTargetTeamId(playerId);
        require(isPlayerFrozen(playerId), "promo player not frozen, cannot complete transfer");
        require(_assets.getOwnerTeam(buyerTeamId) == 
                    recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s]), "Buyer is not own targetTeamId");
        require(now < validUntil);
        require(validUntil == playerIdToAuctionData[playerId], "provided validUntil does not match freeze validUntil");
        require(sig[IDX_MSG] == prefixed(buildAgreeToBuyPromoPlayerTxMsg(playerId, validUntil)), "buyer msg does not match");
        _assets.transferPlayer(playerId, buyerTeamId);
        playerIdToAuctionData[playerId] = 1;
        emit PlayerFreeze(playerId, 1, false);
    }

    function completePlayerAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32 buyerHiddenPrice,
        uint256 buyerTeamId,
        bytes32[3] memory sig,
        uint8 sigV,
        bool isOffer2StartAuction
     ) public {
        require(areCompletePlayerAuctionRequirementsOK(
            sellerHiddenPrice,
            validUntil,
            playerId,
            buyerHiddenPrice,
            buyerTeamId,
            sig,
            sigV,
            isOffer2StartAuction)
            , "requirements to complete auction are not met"    
        );
        _assets.transferPlayer(playerId, buyerTeamId);
        playerIdToAuctionData[playerId] = 1;
        emit PlayerFreeze(playerId, 1, false);
    }
    
    // Main TEAM auction functions: freeze & complete
    function freezeTeam(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 teamId,
        bytes32[3] memory sig,
        uint8 sigV
    ) public {
        require(areFreezeTeamRequirementsOK(sellerHiddenPrice, validUntil, teamId, sig, sigV), "FreePlayer requirements not met");
        // // Freeze player
        teamIdToAuctionData[teamId] = validUntil + uint256(sellerHiddenPrice << 34);
        emit TeamFreeze(teamId, teamIdToAuctionData[teamId], true);
    }

    function completeTeamAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 teamId,
        bytes32 buyerHiddenPrice,
        bytes32[3] memory sig,
        uint8 sigV,
        bool isOffer2StartAuction
     ) public {
        (bool ok, address buyerAddress) = areCompleteTeamAuctionRequirementsOK(
            sellerHiddenPrice,
            validUntil,
            teamId,
            buyerHiddenPrice,
            sig,
            sigV,
            isOffer2StartAuction
        );
        require(ok, "requirements to complete auction are not met");
        _assets.transferTeam(teamId, buyerAddress);
        teamIdToAuctionData[teamId] = 1;
        emit TeamFreeze(teamId, 1, false);
    }

    function areFreezeTeamRequirementsOK(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 teamId,
        bytes32[3] memory sig,
        uint8 sigV
    ) 
        public 
        view 
        returns (bool ok)
    {
        address teamOwner = _assets.getOwnerTeam(teamId);
        ok =    // check validUntil has not expired
                (now < validUntil) &&
                // check player is not already frozen
                (!isTeamFrozen(teamId)) &&  
                // check asset is owned by legit address
                (teamOwner != address(0)) && 
                // check signatures are valid by requiring that they own the asset:
                (teamOwner == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&    
                // check that they signed what they input data says they signed:
                (sig[IDX_MSG] == prefixed(buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, teamId))) &&
                // check that auction time is less that the required 34 bit (17179869183 = 2^34 - 1)
                (validUntil < now + MAX_VALID_UNTIL);
        if (!ok) return false;
        if (teamId == _assets.ROSTER_TEAM()) return true;
        
        // check that the team itself does not have players already for sale:   
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds = _assets.getPlayerIdsInTeam(teamId);
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if ((playerIds[p] != FREE_PLAYER_ID) && isPlayerFrozen(playerIds[p])) return false;
        }
    }

    function areCompleteTeamAuctionRequirementsOK(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 teamId,
        bytes32 buyerHiddenPrice,
        bytes32[3] memory sig,
        uint8 sigV,
        bool isOffer2StartAuction
     ) 
        public
        view
        returns(bool ok, address buyerAddress) 
    {
        // the next line will verify that the teamId is the same that was used by the seller to sign
        bytes32 sellerTxHash = prefixed(buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, teamId));
        buyerAddress = recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s]);
        ok =    // check buyerAddress is legit and signature is valid
                (buyerAddress != address(0)) && 
                // check buyer and seller refer to the exact same auction
                ((uint256(sellerHiddenPrice) % 2**(256-34)) == (teamIdToAuctionData[teamId] >> 34)) &&
                // // check player is still frozen
                isTeamFrozen(teamId) &&
                // // check that they signed what they input data says they signed:
                sig[IDX_MSG] == prefixed(buildAgreeToBuyTeamTxMsg(sellerTxHash, buyerHiddenPrice, isOffer2StartAuction));

        if (isOffer2StartAuction) {
            // in this case: validUntil is interpreted as offerValidUntil
            ok = ok && (validUntil > (teamIdToAuctionData[teamId] & VALID_UNTIL_MASK) - AUCTION_TIME);
        } else {
            ok = ok && (validUntil == (teamIdToAuctionData[teamId] & VALID_UNTIL_MASK));
        } 
    }

    function areFreezePlayerRequirementsOK(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32[3] memory sig,
        uint8 sigV
    ) 
        public 
        view 
        returns (bool)
    {
        return (
            // check validUntil has not expired
            (now < validUntil) &&
            // check player is not already frozen
            (!isPlayerFrozen(playerId)) &&  
            // check that the team it belongs to not already frozen
            !isTeamFrozen(_assets.getCurrentTeamIdFromPlayerId(playerId)) &&
            // check asset is owned by legit address
            (_assets.getOwnerPlayer(playerId) != address(0)) && 
            // check signatures are valid by requiring that they own the asset:
            (_assets.getOwnerPlayer(playerId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&    
            // check that they signed what they input data says they signed:
            (sig[IDX_MSG] == prefixed(buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerId))) &&
            // check that auction time is less that the required 34 bit (17179869183 = 2^34 - 1)
            (validUntil < now + MAX_VALID_UNTIL)
        );
    }


    function areCompletePlayerAuctionRequirementsOK(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32 buyerHiddenPrice,
        uint256 buyerTeamId,
        bytes32[3] memory sig,
        uint8 sigV,
        bool isOffer2StartAuction
     ) 
        public
        view
        returns(bool ok) 
    {
        // the next line will verify that the playerId is the same that was used by the seller to sign
        bytes32 sellerTxHash = prefixed(buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerId));

        ok =    // check asset is owned by buyer
                (_assets.getOwnerTeam(buyerTeamId) != address(0)) && 
                // check buyer and seller refer to the exact same auction
                ((uint256(sellerHiddenPrice) % 2**(256-34)) == (playerIdToAuctionData[playerId] >> 34)) &&
                // check signatures are valid by requiring that they own the asset:
                (_assets.getOwnerTeam(buyerTeamId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&
                // check player is still frozen
                isPlayerFrozen(playerId) &&
                // check that they signed what they input data says they signed:
                sig[IDX_MSG] == prefixed(buildAgreeToBuyPlayerTxMsg(sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction));


        if (isOffer2StartAuction) {
            // in this case: validUntil is interpreted as offerValidUntil
            ok = ok && (validUntil > (playerIdToAuctionData[playerId] & VALID_UNTIL_MASK) - AUCTION_TIME);
        } else {
            ok = ok && (validUntil == (playerIdToAuctionData[playerId] & VALID_UNTIL_MASK));
        } 
    }


    
    
    // this function is not used in the contract. It's only for external helps
    function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) external pure returns (bytes32) {
        return keccak256(abi.encode(currencyId, price, rnd));
    }

    function hashBidHiddenPrice(uint256 extraPrice, uint256 rnd) external pure returns (bytes32) {
        return keccak256(abi.encode(extraPrice, rnd));
    }

    function buildPutAssetForSaleTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 assetId) public pure returns (bytes32) {
        return keccak256(abi.encode(hiddenPrice, validUntil, assetId));
    }

    function buildOfferToBuyTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 playerId, uint256 buyerTeamId) public pure returns (bytes32) {
        return keccak256(abi.encode(hiddenPrice, validUntil, playerId, buyerTeamId));
    }

    function buildAgreeToBuyPlayerTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bool isOffer2StartAuction) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction));
    }

    function buildAgreeToBuyPromoPlayerTxMsg(uint256 playerId, uint256 validUntil) public pure returns (bytes32) {
        return keccak256(abi.encode(playerId, validUntil));
    }

    function buildAgreeToBuyTeamTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, bool isOffer2StartAuction) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerTxHash, buyerHiddenPrice, isOffer2StartAuction));
    }

    // FUNCTIONS FOR SIGNATURE MANAGEMENT
    // retrieves the addr that signed a message
    function recoverAddr(bytes32 msgHash, uint8 v, bytes32 r, bytes32 s) internal pure returns (address) {
        return ecrecover(msgHash, v, r, s);
    }

    // (currently not used) checks if the signature of a message says that it was signed by the provided address
    function isSigned(address _addr, bytes32 msgHash, uint8 v, bytes32 r, bytes32 s) public pure returns (bool) {
        return ecrecover(msgHash, v, r, s) == _addr;
    }

    // Builds a prefixed hash to mimic the behavior of eth_sign.
    function prefixed(bytes32 hash) public pure returns (bytes32) {
        return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", hash));
    }

    function getBlockchainNowTime() external view returns (uint) {
        return now;
    }

    function isPlayerFrozen(uint256 playerId) public view returns (bool) {
        require(_assets.getIsSpecial(playerId) || _assets.playerExists(playerId), "player does not exist");
        return (playerIdToAuctionData[playerId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }

    function isTeamFrozen(uint256 teamId) public view returns (bool) {
        if (teamId == _assets.ROSTER_TEAM()) return false;
        require(_assets.teamExists(teamId), "unexistent team");
        return (teamIdToAuctionData[teamId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }
}