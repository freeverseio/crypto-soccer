pragma solidity >=0.5.12 <0.6.2;

import "./AssetsLib.sol";
import "./EncodingIDs.sol";
import "./EncodingState.sol";
import "./EncodingSkillsGetters.sol";
import "./EncodingSkillsSetters.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract Market is Storage, AssetsLib, EncodingSkillsSetters, EncodingState {
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
    uint256 constant public ACADEMY_TEAM = 1;
    uint8 constant public MAX_ACQUISITON_CONSTAINTS  = 7;
    
    mapping (uint256 => uint256) private _playerIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToRemainingAcqs;

    event TeamTransfer(uint256 teamId, address to);
    event PlayerStateChange(uint256 playerId, uint256 state);

    address public academyAddr;

    function setAcademyAddr(address addr) public {
        academyAddr = addr;
        emit TeamTransfer(ACADEMY_TEAM, addr);        
    }
    
    function isAcademyPlayer(uint256 playerId) public view returns(bool) {
        return (getIsSpecial(playerId) && _playerIdToState[playerId] == 0);
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

    // function transferPromoPlayer(
    //     uint256 playerId,
    //     uint256 validUntil,
    //     bytes32[3] memory sigSel,
    //     bytes32[3] memory sigBuy,
    //     uint8 sigVSel,
    //     uint8 sigVBuy
    //  ) public {
    //     require(validUntil > now, "validUntil is in the past");
    //     require(validUntil < now + MAX_VALID_UNTIL, "validUntil is too large");
    //      uint256 playerIdWithoutTargetTeam = setTargetTeamId(playerId, 0);
    //     require(!_assets.isPlayerWritten(playerIdWithoutTargetTeam), "promo player already in the universe");
    //     uint256 targetTeamId = getTargetTeamId(playerId);
    //     // require that team does not have any constraint from friendlies
    //     (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
    //     require(!(isConstrained && (nRemain == 0)), "trying to accept a promo player, but team is busy in constrained friendlies");
    //     // testing about the target team is already done in _assets.transferPlayer
    //     require(_assets.teamExists(targetTeamId), "cannot offer a promo player to a non-existent team");
    //     require(!_assets.isBotTeam(targetTeamId), "cannot offer a promo player to a bot team");
                
    //     require(_assets.getOwnerTeam(targetTeamId) == 
    //                 recoverAddr(sigBuy[IDX_MSG], sigVBuy, sigBuy[IDX_r], sigBuy[IDX_s]), "Buyer does not own targetTeamId");
         
    //     require(academyAddr == 
    //                 recoverAddr(sigSel[IDX_MSG], sigVSel, sigSel[IDX_r], sigSel[IDX_s]), "Seller does not own academy");
         
    //     bytes32 signedMsg = prefixed(buildPromoPlayerTxMsg(playerId, validUntil));
    //     require(sigBuy[IDX_MSG] == signedMsg, "buyer msg does not match");
    //     require(sigSel[IDX_MSG] == signedMsg, "seller msg does not match");
    //     _assets.transferPlayer(playerIdWithoutTargetTeam, targetTeamId);
    //     decreaseMaxAllowedAcquisitions(targetTeamId);
    // }

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
        transferPlayer(playerId, buyerTeamId);
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
        transferTeam(teamId, buyerAddress);
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
                // check that auction time is less that the required 34 bit (17179869183 = 2^34 - 1)
                (validUntil < now + MAX_VALID_UNTIL);
        if (!ok) return false;
        if (teamId == ACADEMY_TEAM) return true;
        
        // check that the team itself does not have players already for sale:   
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds = getPlayerIdsInTeam(teamId);
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
        address prevOwner = isAcademyPlayer(playerId) ? academyAddr : getOwnerPlayer(playerId);
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
                (getOwnerTeam(buyerTeamId) != address(0)) && 
                // check buyer and seller refer to the exact same auction
                ((uint256(sellerHiddenPrice) % 2**(256-34)) == (_playerIdToAuctionData[playerId] >> 34)) &&
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
    




    
    function transferPlayer(uint256 playerId, uint256 teamIdTarget) public  {
        // warning: check of ownership of players and teams should be done before calling this function
        // TODO: checking if they are bots should be done outside this function
        require(getIsSpecial(playerId) || playerExists(playerId), "player does not exist");
        require(teamExists(teamIdTarget), "unexistent target team");
        uint256 state = getPlayerState(playerId);
        uint256 newState = state;
        uint256 teamIdOrigin = getCurrentTeamId(state);
        require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
        require(!isBotTeam(teamIdOrigin) && !isBotTeam(teamIdTarget), "cannot transfer player when at least one team is a bot");
        uint8 shirtTarget = getFreeShirt(teamIdTarget);
        require(shirtTarget != PLAYERS_PER_TEAM_MAX, "target team for transfer is already full");
        
        newState = setCurrentTeamId(newState, teamIdTarget);
        newState = setCurrentShirtNum(newState, shirtTarget);
        newState = setLastSaleBlock(newState, block.number);
        _playerIdToState[playerId] = newState;

        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamIdTarget);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtTarget] = playerId;
        if (teamIdOrigin != ACADEMY_TEAM) {
            uint256 shirtOrigin = getCurrentShirtNum(state);
            (timeZone, countryIdxInTZ, teamIdxInCountry) = decodeTZCountryAndVal(teamIdOrigin);
            _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtOrigin] = FREE_PLAYER_ID;
        }
        emit PlayerStateChange(playerId, newState);
    }
    
    
    function getOwnerTeam(uint256 teamId) public view returns(address) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    

    // returns NULL_ADDR if team is bot
    function getOwnerPlayer(uint256 playerId) public view returns(address) {
        require(playerExists(playerId), "unexistent player");
        return getOwnerTeam(getCurrentTeamIdFromPlayerId(playerId));
    }
    
    function playerExists(uint256 playerId) public view returns (bool) {
        if (playerId == 0) return false;
        if (getIsSpecial(playerId)) return (_playerIdToState[playerId] != 0);
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        return _wasPlayerCreatedInCountry(timeZone, countryIdxInTZ, playerIdxInCountry);
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
    
    function transferTeamInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) private {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a bot team");
        require(addr != NULL_ADDR, "cannot transfer to a null address");
        require(_timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner != addr, "buyer and seller are the same addr");
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = addr;
    }

    function transferTeam(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        transferTeamInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
        emit TeamTransfer(teamId, addr);
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


    function teamExists(uint256 teamId) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return _teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }
    
    function _wasPlayerCreatedInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) private view returns(bool) {
        return (playerIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ) * PLAYERS_PER_TEAM_INIT);
    }


    function isBotTeam(uint256 teamId) public view returns(bool) {
        if (teamId == ACADEMY_TEAM) return false;
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function getFreeShirt(uint256 teamId) public view returns(uint8) {
        for (uint8 shirtNum = PLAYERS_PER_TEAM_MAX-1; shirtNum >= 0; shirtNum--) {
            if (isFreeShirt(teamId, shirtNum)) {
                return shirtNum;
            }
        }
        return PLAYERS_PER_TEAM_MAX;
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
            return FREE_PLAYER_ID;
        } else {
            return encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry * PLAYERS_PER_TEAM_INIT + shirtNum);
        }
    }
       
    function isFreeShirt(uint256 teamId, uint8 shirtNum) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry),"cannot query about the shirt of a Bot Team");
        uint256 writtenId = _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
        if (shirtNum > PLAYERS_PER_TEAM_INIT - 1) {
            return (writtenId == 0 || writtenId == FREE_PLAYER_ID);
        } else {
            return writtenId == FREE_PLAYER_ID;
        }
    }


}
