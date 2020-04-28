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
    event PlayerFreezeCrypto(uint256 playerId, bool frozen);
    event TeamFreeze(uint256 teamId, uint256 auctionData, bool frozen);
    event PlayerStateChange(uint256 playerId, uint256 state);

    function setIsPlayerFrozenCrypto(uint256 playerId, bool isFrozen) public {
        _playerIdToIsFrozenCrypto[playerId] = isFrozen;
        emit PlayerFreezeCrypto(playerId, isFrozen);
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
    
    // Main PLAYER auction functions: freeze & complete
    function freezePlayer(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32[2] memory sig,
        uint8 sigV
    ) public {
        require(areFreezePlayerRequirementsOK(sellerHiddenPrice, validUntil, playerId, sig, sigV), "FreezePlayer requirements not met");
        // // Freeze player
        _playerIdToAuctionData[playerId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
        emit PlayerFreeze(playerId, _playerIdToAuctionData[playerId], true);
    }

    function transferBuyNowPlayer(
        uint256 playerId,
        uint256 targetTeamId
     ) public {
        // isAcademy checks that player isSpecial, and not written.
        require(getCurrentTeamIdFromPlayerId(playerId) == ACADEMY_TEAM, "only Academy players can be sold via buy-now");
        require(getTargetTeamId(playerId) == 0, "cannot have buy-now players with non-null targetTeamId");

        // note that wasTeamCreatedVirtually(targetTeamId) &  !isBotTeam(targetTeamId) => already part of transferPlayer
        (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
        require(!(isConstrained && (nRemain == 0)), "trying to accept a promo player, but team is busy in constrained friendlies");
        transferPlayer(playerId, targetTeamId);
        decreaseMaxAllowedAcquisitions(targetTeamId);
    }
    
    function transferPromoPlayer(
        uint256 playerId,
        uint256 validUntil,
        bytes32[2] memory sigSel,
        bytes32[2] memory sigBuy,
        uint8 sigVSel,
        uint8 sigVBuy
     ) public {
        require(validUntil > now, "validUntil is in the past");
        require(validUntil < now + MAX_VALID_UNTIL, "validUntil is too large");
        require(getCurrentTeamIdFromPlayerId(playerId) == ACADEMY_TEAM, "only Academy Players can be offered as promo players");
        uint256 playerIdWithoutTargetTeam = setTargetTeamId(playerId, 0);
        require(_playerIdToState[playerIdWithoutTargetTeam] == 0, "promo player already in the universe");
        uint256 targetTeamId = getTargetTeamId(playerId);
        // require that team does not have any constraint from friendlies
        (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
        require(!(isConstrained && (nRemain == 0)), "trying to accept a promo player, but team is busy in constrained friendlies");
        // testing require(wasTeamCreatedVirtually(targetTeamId) and  require(!isBotTeam(targetTeamId) is already done in _assets.transferPlayer:
                
        bytes32 signedMsg = prefixed(buildPromoPlayerTxMsg(playerId, validUntil));
        require(getOwnerTeam(targetTeamId) == 
                    recoverAddr(signedMsg, sigVBuy, sigBuy[IDX_r], sigBuy[IDX_s]), "Buyer does not own targetTeamId");
        require(_academyAddr == 
                    recoverAddr(signedMsg, sigVSel, sigSel[IDX_r], sigSel[IDX_s]), "Seller does not own academy");
         
        transferPlayer(playerIdWithoutTargetTeam, targetTeamId);
        decreaseMaxAllowedAcquisitions(targetTeamId);
    }

    function completePlayerAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32 buyerHiddenPrice,
        uint256 buyerTeamId,
        bytes32[2] memory sig,
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
        bytes32[2] memory sig,
        uint8 sigV
    ) public {
        require(areFreezeTeamRequirementsOK(sellerHiddenPrice, validUntil, teamId, sig, sigV), "FreezePlayer requirements not met");
        // // Freeze player
        _teamIdToAuctionData[teamId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
        emit TeamFreeze(teamId, _teamIdToAuctionData[teamId], true);
    }

    function completeTeamAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 teamId,
        bytes32 buyerHiddenPrice,
        bytes32[2] memory sig,
        uint8 sigV,
        address buyerAddress,
        bool isOffer2StartAuction
     ) public {
        bool ok = areCompleteTeamAuctionRequirementsOK(
            sellerHiddenPrice,
            validUntil,
            teamId,
            buyerHiddenPrice,
            sig,
            sigV,
            buyerAddress,
            isOffer2StartAuction
        );
        require(ok, "requirements to complete auction are not met");
        transferTeam(teamId, buyerAddress);
        _teamIdToAuctionData[teamId] = 1;
        emit TeamFreeze(teamId, 1, false);
    }

    function transferPlayer(uint256 playerId, uint256 teamIdTarget) public  {
        // warning: check of ownership of players and teams should be done before calling this function
        // warning: checking if they are bots should be moved outside this function

        // part related to origin team:
        uint256 state = getPlayerState(playerId);
        uint256 teamIdOrigin = getCurrentTeamIdFromPlayerState(state);
        require(teamIdOrigin != NULL_TEAMID, "the player does not belong to any team");
    
        if (teamIdOrigin != ACADEMY_TEAM) {
            uint256 shirtOrigin = getCurrentShirtNum(state);
            teamIdToPlayerIds[teamIdOrigin][shirtOrigin] = FREE_PLAYER_ID;
        }
                
        // part related to both teams:
        require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
        require(!isBotTeam(teamIdOrigin) && !isBotTeam(teamIdTarget), "cannot transfer player when at least one team is a bot");

        // part related to target team:
        require(wasTeamCreatedVirtually(teamIdTarget), "unexistent target team");
        uint8 shirtTarget = getFreeShirt(teamIdTarget);
        if (shirtTarget < PLAYERS_PER_TEAM_MAX) {
            state = setLastSaleBlock(
                        setCurrentShirtNum(
                            setCurrentTeamId(
                                state, teamIdTarget
                            ), shirtTarget
                        ), block.number
                    );
            teamIdToPlayerIds[teamIdTarget][shirtTarget] = playerId;
        } else {
            _playerInTransitToTeam[playerId] = teamIdTarget;
            _nPlayersInTransitInTeam[teamIdTarget] += 1;
            state = setLastSaleBlock(
                        setCurrentTeamId(
                            state, IN_TRANSIT_TEAM
                        ), block.number
                    );
        }
        _playerIdToState[playerId] = state;

        emit PlayerStateChange(playerId, state);
    }
    
    function completePlayerTransit(uint256 playerId) public  {
        uint256 teamIdTarget = _playerInTransitToTeam[playerId];
        require(teamIdTarget != 0, "player not in transit");
        uint8 shirtTarget = getFreeShirt(teamIdTarget);
        require(shirtTarget < PLAYERS_PER_TEAM_MAX, "cannot complete player transit because targetTeam is still full");
        uint256 state = getPlayerState(playerId);
        state = setCurrentShirtNum(
                    setCurrentTeamId(
                        state, teamIdTarget
                    ), shirtTarget
                );
        _playerIdToState[playerId] = state;
        teamIdToPlayerIds[teamIdTarget][shirtTarget] = playerId;
        _nPlayersInTransitInTeam[teamIdTarget] -= 1;
        delete _playerInTransitToTeam[playerId];
        emit PlayerStateChange(playerId, state);
    }
        
    function transferTeam(uint256 teamId, address addr) public {
        // requiring that team is not bot already ensures that tz and countryIdxInTz exist 
        require(!isBotTeam(teamId), "cannot transfer a bot team");
        require(addr != NULL_ADDR, "cannot transfer to a null address");
        require(teamIdToOwner[teamId] != addr, "buyer and seller are the same addr");
        teamIdToOwner[teamId] = addr;
        emit TeamTransfer(teamId, addr);
    }
    
    // function dismissPlayer(uint256 playerId) public {
    //     require(playerExists(playerId), "player does not exist");
    //     require(!isPlayerFrozenInAnyMarket(playerId),"cannot dismiss a player that is frozen");

    //     uint256 state = getPlayerState(playerId);
    //     uint256 teamIdOrigin = getCurrentTeamIdFromPlayerState(state);
    //     require(teamIdOrigin != ACADEMY_TEAM, "cannot dimiss a player from the Academy team");
    //     uint256 shirtOrigin = getCurrentShirtNum(state);
    //     teamIdToPlayerIds[teamIdOrigin][shirtOrigin] = FREE_PLAYER_ID;
    //     require(_nPlayersInTransitInTeam[teamIdTarget] != 0, "cannot dimiss a player unless team is full");

    //     require(!isBotTeam(teamIdOrigin), "cannot transfer player when at least one team is a bot");

    //     _playerIdToState[playerId] = state;

    //     emit PlayerStateChange(playerId, state);
    // }
    
    
    
}
