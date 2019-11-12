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
    uint256 constant public ROSTER_TEAM = 1;
    uint8 constant public MAX_ACQUISITON_CONSTAINTS  = 7;

    Assets private _assets;

    mapping (uint256 => uint256) private _playerIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToRemainingAcqs;

    address public rosterAddr;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }
    function setRosterAddr(address addr) public {rosterAddr = addr;}
    
    function isRosterPlayer(uint256 playerId) public view returns(bool) {
        return (_assets.getIsSpecial(playerId) && !_assets.isPlayerWritten(playerId));
    }

    function addAcquisitionConstraint(uint256 teamId, uint32 validUntil, uint8 nRemain) public {
        require(nRemain > 0, "nRemain = 0, which does not make sense for a constraint");
        uint256 remainingAcqs = _teamIdToRemainingAcqs[teamId];
        bool success;
        for (uint8 acq = 0; acq < MAX_ACQUISITON_CONSTAINTS; acq++) {
            if (isAcquisitionFree(remainingAcqs, acq)) {
                _teamIdToRemainingAcqs[teamId] = setAcquisitionConstraint(remainingAcqs, validUntil, nRemain, acq);
                success = true;
                continue;
            }
        }
        require(success, "this team is already signed up in 7 contrained friendly championships");
    }
    
    function decreaseMaxAllowedAcquisitions(uint256 teamId) public {
        uint256 remainingAcqs = _teamIdToRemainingAcqs[teamId];
        if (remainingAcqs == 0) return;
        for (uint8 acq = 0; acq < MAX_ACQUISITON_CONSTAINTS; acq++) {
            if (!isAcquisitionFree(remainingAcqs, acq)) {
                remainingAcqs = decreaseAcquisitionConstraint(remainingAcqs, acq);
            }
        }
        _teamIdToRemainingAcqs[teamId] = remainingAcqs;
    }
    
    function getMaxAllowedAcquisitions(uint256 teamId) public view returns (bool isConstrained, uint8) {
        uint256 remainingAcqs = _teamIdToRemainingAcqs[teamId];
        if (remainingAcqs == 0) return (false, 0);
        uint8 nRemain = 255;
        for (uint8 acq = 0; acq < MAX_ACQUISITON_CONSTAINTS; acq++) {
            if (!isAcquisitionFree(remainingAcqs, acq)) {
                uint8 thisNRemain = getAcquisitionConstraintNRemain(remainingAcqs, acq);
                if (thisNRemain == 0) return (true, 0);
                if (thisNRemain < nRemain) nRemain = thisNRemain;
                
            }
        }
        return (nRemain == 255) ? (false, 0) : (true, nRemain);
    }    
    
    function isAcquisitionFree(uint256 remainingAcqs, uint8 acq) public view returns (bool) {
        uint32 validUntil = getAcquisitionConstraintValidUntil(remainingAcqs, acq);
        return (validUntil == 0) || (validUntil < now);
    }
    
    function getAcquisitionConstraintValidUntil(uint256 remainingAcqs, uint8 acq) public pure returns (uint32) {
        return uint32((remainingAcqs >> 36 * acq) & (2**32-1)) ; 
    }

    function getAcquisitionConstraintNRemain(uint256 remainingAcqs, uint8 acq) public pure returns (uint8) {
        return uint8((remainingAcqs >> 32 + 36 * acq) & 15); 
    }
    
    function setAcquisitionConstraint(uint256 remainingAcqs, uint32 validUntil, uint8 nRemain, uint8 acq) public pure returns (uint256) {
        return  (remainingAcqs & ~(uint256(2**36-1) << (36 * acq)))
                | (uint256(validUntil) << (36 * acq))
                | (uint256(nRemain) << (32 + 36 * acq));
    }

    function decreaseAcquisitionConstraint(uint256 remainingAcqs, uint8 acq) public pure returns (uint256) {
        uint8 nRemain = getAcquisitionConstraintNRemain(remainingAcqs, acq);
        return (remainingAcqs & ~(uint256(15) << (32 + 36 * acq))) | (uint256(nRemain-1) << (32 + 36 * acq));
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
        _playerIdToAuctionData[playerId] = validUntil + uint256(sellerHiddenPrice << 34);
        emit PlayerFreeze(playerId, _playerIdToAuctionData[playerId], true);
    }

    function transferPromoPlayer(
        uint256 playerId,
        uint256 validUntil,
        bytes32[3] memory sigSel,
        bytes32[3] memory sigBuy,
        uint8 sigVSel,
        uint8 sigVBuy
     ) public {
        require(validUntil > now, "validUntil is in the past");
        require(validUntil < now + MAX_VALID_UNTIL, "validUntil is too large");
        uint256 playerIdWithoutTargetTeam = _assets.setTargetTeamId(playerId, 0);
        require(!_assets.isPlayerWritten(playerIdWithoutTargetTeam), "promo player already in the universe");
        uint256 targetTeamId = _assets.getTargetTeamId(playerId);
        // testing about the target team is already done in _assets.transferPlayer
        require(_assets.teamExists(targetTeamId), "cannot offer a promo player to a non-existent team");
        require(!_assets.isBotTeam(targetTeamId), "cannot offer a promo player to a bot team");
                
        require(_assets.getOwnerTeam(targetTeamId) == 
                    recoverAddr(sigBuy[IDX_MSG], sigVBuy, sigBuy[IDX_r], sigBuy[IDX_s]), "Buyer does not own targetTeamId");
         
        require(rosterAddr == 
                    recoverAddr(sigSel[IDX_MSG], sigVSel, sigSel[IDX_r], sigSel[IDX_s]), "Seller does not own roster");
         
        bytes32 signedMsg = prefixed(buildPromoPlayerTxMsg(playerId, validUntil));
        require(sigBuy[IDX_MSG] == signedMsg, "buyer msg does not match");
        require(sigSel[IDX_MSG] == signedMsg, "seller msg does not match");
        _assets.transferPlayer(playerIdWithoutTargetTeam, targetTeamId);
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
        _playerIdToAuctionData[playerId] = 1;
        decreaseMaxAllowedAcquisitions(buyerTeamId);
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
        _teamIdToAuctionData[teamId] = validUntil + uint256(sellerHiddenPrice << 34);
        emit TeamFreeze(teamId, _teamIdToAuctionData[teamId], true);
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
        _teamIdToAuctionData[teamId] = 1;
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
                ((uint256(sellerHiddenPrice) % 2**(256-34)) == (_teamIdToAuctionData[teamId] >> 34)) &&
                // // check player is still frozen
                isTeamFrozen(teamId) &&
                // // check that they signed what they input data says they signed:
                sig[IDX_MSG] == prefixed(buildAgreeToBuyTeamTxMsg(sellerTxHash, buyerHiddenPrice, isOffer2StartAuction));

        if (isOffer2StartAuction) {
            // in this case: validUntil is interpreted as offerValidUntil
            ok = ok && (validUntil > (_teamIdToAuctionData[teamId] & VALID_UNTIL_MASK) - AUCTION_TIME);
        } else {
            ok = ok && (validUntil == (_teamIdToAuctionData[teamId] & VALID_UNTIL_MASK));
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
        address prevOwner = isRosterPlayer(playerId) ? rosterAddr : _assets.getOwnerPlayer(playerId);
        return (
            // check validUntil has not expired
            (now < validUntil) &&
            // check player is not already frozen
            (!isPlayerFrozen(playerId)) &&  
            // check that the team it belongs to not already frozen
            !isTeamFrozen(_assets.getCurrentTeamIdFromPlayerId(playerId)) &&
            // check asset is owned by legit address
            (prevOwner != address(0)) && 
            // check signatures are valid by requiring that they own the asset:
            (prevOwner == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&    
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
        
        (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(buyerTeamId);
        if (isConstrained && nRemain == 0) return false;
        ok =    // check asset is owned by buyer
                (_assets.getOwnerTeam(buyerTeamId) != address(0)) && 
                // check buyer and seller refer to the exact same auction
                ((uint256(sellerHiddenPrice) % 2**(256-34)) == (_playerIdToAuctionData[playerId] >> 34)) &&
                // check signatures are valid by requiring that they own the asset:
                (_assets.getOwnerTeam(buyerTeamId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&
                // check player is still frozen
                isPlayerFrozen(playerId) &&
                // check that they signed what they input data says they signed:
                sig[IDX_MSG] == prefixed(buildAgreeToBuyPlayerTxMsg(sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction));


        if (isOffer2StartAuction) {
            // in this case: validUntil is interpreted as offerValidUntil
            ok = ok && (validUntil > (_playerIdToAuctionData[playerId] & VALID_UNTIL_MASK) - AUCTION_TIME);
        } else {
            ok = ok && (validUntil == (_playerIdToAuctionData[playerId] & VALID_UNTIL_MASK));
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

    function buildPromoPlayerTxMsg(uint256 playerId, uint256 validUntil) public pure returns (bytes32) {
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
        return (_playerIdToAuctionData[playerId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }

    function isTeamFrozen(uint256 teamId) public view returns (bool) {
        if (teamId == _assets.ROSTER_TEAM()) return false;
        require(_assets.teamExists(teamId), "unexistent team");
        return (_teamIdToAuctionData[teamId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }
}