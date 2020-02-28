pragma solidity >=0.5.12 <=0.6.3;

import "./MarketView.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 * @dev The serialized structs appearing here are "AcquisitonConstraints" and "AuctionData"
 * @dev     Both use validUntil (in seconds) which uses 32b, hence allowing 2**32/(3600*24*365) = 136 years after 1970
 * @dev     AuctionData encodes, (8b of zeroes, 216b for sellerHiddenPrice, 32b for validUntil), 
 * @dev         where sellerHiddenPrice has the leftmost 40 bit killed, 
 * @dev         => validUntil + (uint256(sellerHiddenPrice) << 40)) >> 8;
 * @dev     AcquisitonConstraints: serializes the number of trades left (4b), and until when, for the 6 possible constraints
 * @dev         => (n5, validUntil5, n4, validUntil4,... n0, validUntil0), 
 * @dev         => so it leaves the leftmost 256 - 6 * 36 = 40b free
 */
 
contract Market is MarketView {
    event PlayerFreeze(uint256 playerId, uint256 auctionData, bool frozen);
    event TeamFreeze(uint256 teamId, uint256 auctionData, bool frozen);
    event PlayerStateChange(uint256 playerId, uint256 state);

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
        _playerIdToAuctionData[playerId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
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
        uint256 playerIdWithoutTargetTeam = setTargetTeamId(playerId, 0);
        require(!isPlayerWritten(playerIdWithoutTargetTeam), "promo player already in the universe");
        uint256 targetTeamId = getTargetTeamId(playerId);
        // require that team does not have any constraint from friendlies
        (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
        require(!(isConstrained && (nRemain == 0)), "trying to accept a promo player, but team is busy in constrained friendlies");
        // testing about the target team is already done in _assets.transferPlayer
        require(teamExists(targetTeamId), "cannot offer a promo player to a non-existent team");
        require(!isBotTeam(targetTeamId), "cannot offer a promo player to a bot team");
                
        require(getOwnerTeam(targetTeamId) == 
                    recoverAddr(sigBuy[IDX_MSG], sigVBuy, sigBuy[IDX_r], sigBuy[IDX_s]), "Buyer does not own targetTeamId");
         
        require(_academyAddr == 
                    recoverAddr(sigSel[IDX_MSG], sigVSel, sigSel[IDX_r], sigSel[IDX_s]), "Seller does not own academy");
         
        bytes32 signedMsg = prefixed(buildPromoPlayerTxMsg(playerId, validUntil));
        require(sigBuy[IDX_MSG] == signedMsg, "buyer msg does not match");
        require(sigSel[IDX_MSG] == signedMsg, "seller msg does not match");
        transferPlayer(playerIdWithoutTargetTeam, targetTeamId);
        decreaseMaxAllowedAcquisitions(targetTeamId);
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
        _teamIdToAuctionData[teamId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
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
        
        _playerIdToState[playerId] = 
            setLastSaleBlock(
                setCurrentShirtNum(
                    setCurrentTeamId(
                        newState, teamIdTarget
                    ), shirtTarget
                ), block.number
            );

        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamIdTarget);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToPlayerIds[teamIdxInCountry][shirtTarget] = playerId;
        if (teamIdOrigin != ACADEMY_TEAM) {
            uint256 shirtOrigin = getCurrentShirtNum(state);
            (timeZone, countryIdxInTZ, teamIdxInCountry) = decodeTZCountryAndVal(teamIdOrigin);
            _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToPlayerIds[teamIdxInCountry][shirtOrigin] = FREE_PLAYER_ID;
        }
        emit PlayerStateChange(playerId, newState);
    }
    
    function transferTeamInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) private {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a bot team");
        require(addr != NULL_ADDR, "cannot transfer to a null address");
        require(_timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToOwner[teamIdxInCountry] != addr, "buyer and seller are the same addr");
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToOwner[teamIdxInCountry] = addr;
    }

    function transferTeam(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        transferTeamInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
        emit TeamTransfer(teamId, addr);
    }
}
