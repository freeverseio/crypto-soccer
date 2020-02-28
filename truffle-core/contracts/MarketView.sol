pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsLib.sol";
import "./EncodingState.sol";
import "./EncodingSkillsSetters.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract MarketView is AssetsLib, EncodingSkillsSetters, EncodingState {

    
    
    function isAcademyPlayer(uint256 playerId) public view returns(bool) {
        return (getIsSpecial(playerId) && _playerIdToState[playerId] == 0);
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
        address teamOwner = getOwnerTeam(teamId);
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
                // check that auction time is less that the required 32 bit (2^32 - 1)
                (validUntil < now + MAX_VALID_UNTIL);
        if (!ok) return false;
        if (teamId == ACADEMY_TEAM) return true;
        
        // check that the team itself does not have players already for sale:   
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds = getPlayerIdsInTeam(teamId);
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (!isFreeShirt(playerIds[p], p) && isPlayerFrozen(playerIds[p])) return false;
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
                ((uint256(sellerHiddenPrice) & KILL_LEFTMOST_40BIT_MASK) == (_teamIdToAuctionData[teamId] >> 32)) &&
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
                (getOwnerTeam(buyerTeamId) != address(0)) && 
                // check buyer and seller refer to the exact same auction
                ((uint256(sellerHiddenPrice) & KILL_LEFTMOST_40BIT_MASK) == (_playerIdToAuctionData[playerId] >> 32)) &&
                // check signatures are valid by requiring that they own the asset:
                (getOwnerTeam(buyerTeamId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&
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
        address prevOwner = isAcademyPlayer(playerId) ? _academyAddr : getOwnerTeam(getCurrentTeamIdFromPlayerId(playerId));
        return (
            // check validUntil has not expired
            (now < validUntil) &&
            // check player is not already frozen
            (!isPlayerFrozen(playerId)) &&  
            // check that the team it belongs to not already frozen
            !isTeamFrozen(getCurrentTeamIdFromPlayerId(playerId)) &&
            // check asset is owned by legit address
            (prevOwner != address(0)) && 
            // check signatures are valid by requiring that they own the asset:
            (prevOwner == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&    
            // check that they signed what they input data says they signed:
            (sig[IDX_MSG] == prefixed(buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerId))) &&
            // check that auction time is less that the required 32 bit
            (validUntil < now + MAX_VALID_UNTIL)
        );
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
        require(getIsSpecial(playerId) || playerExists(playerId), "player does not exist");
        return (_playerIdToAuctionData[playerId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }

    function isTeamFrozen(uint256 teamId) public view returns (bool) {
        if (teamId == ACADEMY_TEAM) return false;
        require(teamExists(teamId), "unexistent team");
        return (_teamIdToAuctionData[teamId] & VALID_UNTIL_MASK) + POST_AUCTION_TIME > now;
    }
    
    function getOwnerTeam(uint256 teamId) public view returns(address) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function getPlayerStateAtBirth(uint256 playerId) public pure returns (uint256) {
        if (getIsSpecial(playerId)) return encodePlayerState(playerId, ACADEMY_TEAM, 0, 0, 0);
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        uint256 currentTeamId = encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry);
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        return encodePlayerState(playerId, currentTeamId, shirtNum, 0, 0);
    }

    function getPlayerState(uint256 playerId) public view returns (uint256) {
        return (isPlayerWritten(playerId) ? _playerIdToState[playerId] : getPlayerStateAtBirth(playerId));
    }
    
    // TODO: we really don't need this function. Only for external use. Consider removal
    function getPlayerIdsInTeam(uint256 teamId) public view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        if (isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry)) {
            for (uint8 shirtNum = 0 ; shirtNum < PLAYERS_PER_TEAM_MAX ; shirtNum++){
                playerIds[shirtNum] = getDefaultPlayerIdForTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum);
            }
        } else {
            for (uint8 shirtNum = 0 ; shirtNum < PLAYERS_PER_TEAM_MAX ; shirtNum++){
                uint256 writtenId = _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
                if (writtenId == 0) {
                    playerIds[shirtNum] = getDefaultPlayerIdForTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum);
                } else {
                    playerIds[shirtNum] = writtenId;
                }
            }
        }
    }
    
    function getCurrentTeamIdFromPlayerId(uint256 playerId) public view returns(uint256) {
        return getCurrentTeamId(getPlayerState(playerId));
    }

    function isBotTeam(uint256 teamId) public view returns(bool) {
        if (teamId == ACADEMY_TEAM) return false;
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function getFreeShirt(uint256 teamId) public view returns(uint8) {
        // already assumes that there was a previous check that this team is not a bot
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        for (uint8 shirtNum = PLAYERS_PER_TEAM_MAX-1; shirtNum >= 0; shirtNum--) {
            uint256 writtenId = _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
            if (isFreeShirt(writtenId, shirtNum)) { return shirtNum; }
        }
        return PLAYERS_PER_TEAM_MAX;
    }
    
    function isFreeShirt(uint256 playerId, uint8 shirtNum) public pure returns(bool) {
        if (shirtNum < PLAYERS_PER_TEAM_INIT) {
            return (playerId == FREE_PLAYER_ID);
        } else {
            return (playerId == 0 || playerId == FREE_PLAYER_ID);
        }
    }

    function isPlayerWritten(uint256 playerId) public view returns (bool) { return (_playerIdToState[playerId] != 0); }       

    function getDefaultPlayerIdForTeamInCountry(
        uint8 timeZone,
        uint256 countryIdxInTZ,
        uint256 teamIdxInCountry,
        uint8 shirtNum
    )
        public
        pure
        returns(uint256)
    {
        if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
            return 0;
        } else {
            return encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry * PLAYERS_PER_TEAM_INIT + shirtNum);
        }
    }
       
    function getOwnerPlayer(uint256 playerId) public view returns(address) {
        require(playerExists(playerId), "unexistent player");
        return getOwnerTeam(getCurrentTeamIdFromPlayerId(playerId));
    }
}
