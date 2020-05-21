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
    event PlayerRetired(uint256 playerId, uint256 teamId);
    event ProposedNewMaxSumSkillsBuyNowPlayer(uint256 newSumSkills, uint256 newLapseTime);
    event UpdatedNewMaxSumSkillsBuyNowPlayer(uint256 newSumSkills, uint256 newLapseTime);

    modifier onlyCryptoMarket() {
        require(msg.sender == _cryptoMktAddr, "Only CryptoMarket is authorized.");
        _;
    }

    function setCryptoMarketAddress(address addr) external onlyCOO {
        _cryptoMktAddr = addr;
    }
    
    function setIsBuyNowAllowedByOwner(uint256 teamId, bool isAllowed) external {
        require(msg.sender == getOwnerTeam(teamId), "only owner of team can change isBuyNowAlloed");
        _teamIdToIsBuyNowForbidden[teamId] = !isAllowed;
    }
    
    function setIsPlayerFrozenCrypto(uint256 playerId, bool isFrozen) public onlyCryptoMarket {
        _playerIdToIsFrozenCrypto[playerId] = isFrozen;
        emit PlayerFreezeCrypto(playerId, isFrozen);
    }

    function proposeNewMaxSumSkillsBuyNowPlayer(uint256 newSumSkills, uint256 newLapseTime) public onlyCOO{
        _maxSumSkillsBuyNowPlayerProposed = newSumSkills;
        _maxSumSkillsBuyNowPlayerMinLapseProposed = newLapseTime;
        _maxSumSkillsBuyNowPlayerLastUpdate = now;
        emit ProposedNewMaxSumSkillsBuyNowPlayer(newSumSkills, newLapseTime);
    }

    // maxSumSkills can always be lowered, regardless of lapse period 
    function lowerNewMaxSumSkillsBuyNowPlayer(uint256 newMaxSum) public onlyCOO {
        require (newMaxSum < _maxSumSkillsBuyNowPlayer, "newMaxSum is not lower than previous");
        _maxSumSkillsBuyNowPlayer = newMaxSum;
        emit UpdatedNewMaxSumSkillsBuyNowPlayer(newMaxSum, _maxSumSkillsBuyNowPlayerMinLapse);
    }
    
    // maxSumSkills can only grow if enough time has passed 
    function updateNewMaxSumSkillsBuyNowPlayer() public onlyCOO {
        require (now >= (_maxSumSkillsBuyNowPlayerLastUpdate + _maxSumSkillsBuyNowPlayerMinLapse),
            "not enough time passed to update new maxSumSkills"
        );
        _maxSumSkillsBuyNowPlayer = _maxSumSkillsBuyNowPlayerProposed;
        _maxSumSkillsBuyNowPlayerMinLapse = _maxSumSkillsBuyNowPlayerMinLapseProposed;
        emit UpdatedNewMaxSumSkillsBuyNowPlayer(_maxSumSkillsBuyNowPlayer, _maxSumSkillsBuyNowPlayerMinLapse);
    }
    
    // TODO: require signature from team owner
    function addAcquisitionConstraint(uint256 teamId, uint32 validUntil, uint8 nRemain) public onlyCOO {
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
    
    function decreaseMaxAllowedAcquisitions(uint256 teamId) private {
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
    ) public onlyMarket {
        require(areFreezePlayerRequirementsOK(sellerHiddenPrice, validUntil, playerId, sig, sigV), "FreezePlayer requirements not met");
        // // Freeze player
        _playerIdToAuctionData[playerId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
        emit PlayerFreeze(playerId, _playerIdToAuctionData[playerId], true);
    }

    function transferBuyNowPlayer(
        uint256 playerId,
        uint256 targetTeamId
     ) public onlyMarket {
        // isAcademy checks that player isSpecial, and not written.
        require(getCurrentTeamIdFromPlayerId(playerId) == ACADEMY_TEAM, "only Academy players can be sold via buy-now");
        require(getSumOfSkills(playerId) < _maxSumSkillsBuyNowPlayer, "buy now player has sum of skills larger than allowed");
        require(!isBotTeam(targetTeamId), "cannot transfer to bot teams");
        require(!_teamIdToIsBuyNowForbidden[targetTeamId], "user has explicitly forbidden buyNow");
        require(targetTeamId != ACADEMY_TEAM, "targetTeam of buyNow player cannot be Academy Team");

        // note that wasTeamCreatedVirtually(targetTeamId) &  !isBotTeam(targetTeamId) => already part of transferPlayer
        (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
        require(!(isConstrained && (nRemain == 0)), "trying to accept a buyNow player, but team is busy in constrained friendlies");
        transferPlayer(playerId, targetTeamId);
        decreaseMaxAllowedAcquisitions(targetTeamId);
    }
    
    function transferPlayerFromCryptoMkt(
        uint256 playerId,
        uint256 targetTeamId
     ) external onlyCryptoMarket {
        transferPlayer(playerId, targetTeamId);
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
     ) public onlyMarket {
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
    ) public onlyMarket {
        require(areFreezeTeamRequirementsOK(sellerHiddenPrice, validUntil, teamId, sig, sigV), "FreezeTeam requirements not met");
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
     ) public onlyMarket {
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
    
    function dismissPlayer(
        uint256 validUntil,
        uint256 playerId,
        bytes32 sigR,
        bytes32 sigS,
        uint8 sigV,
        bool returnToAcademy
    ) public onlyMarket {
        uint256 state = getPlayerState(playerId);
        uint256 teamIdOrigin = getCurrentTeamIdFromPlayerState(state);
        address owner = getOwnerTeam(teamIdOrigin);
        bytes32 msgHash = prefixed(keccak256(abi.encode(validUntil, playerId, returnToAcademy)));
        require (
            // check validUntil has not expired
            (now < validUntil) &&
            // check player is not already frozen
            (!isPlayerFrozenInAnyMarket(playerId)) &&  
            // check that the team it belongs to not already frozen
            !isTeamFrozen(getCurrentTeamIdFromPlayerId(playerId)) &&
            // check asset is owned by legit address
            (owner != address(0)) && 
            // check signatures are valid by requiring that they own the asset:
            (owner == recoverAddr(msgHash, sigV, sigR, sigS)) &&    
            // check that auction time is less that the required 32 bit
            (validUntil < now + MAX_VALID_UNTIL),
            "conditions to dismiss player are not met"
        );  
        if (returnToAcademy) { 
            transferPlayer(playerId, ACADEMY_TEAM); 
        } else {
            uint256 shirtOrigin = getCurrentShirtNum(state);
            teamIdToPlayerIds[teamIdOrigin][shirtOrigin] = FREE_PLAYER_ID;
            emit PlayerRetired(playerId, teamIdOrigin);
        }
    }

    function transferPlayer(uint256 playerId, uint256 teamIdTarget) private  {
        // warning: check of ownership of players and teams should be done before calling this function
        // so in this function, both teams are asumed to exist, be different, and belong to the rightful (nonBot) owners

        // part related to origin team:
        uint256 state = getPlayerState(playerId);
        uint256 teamIdOrigin = getCurrentTeamIdFromPlayerState(state);
    
        if (teamIdOrigin != ACADEMY_TEAM) {
            uint256 shirtOrigin = getCurrentShirtNum(state);
            teamIdToPlayerIds[teamIdOrigin][shirtOrigin] = FREE_PLAYER_ID;
        }

        // part related to target team:
        // - determine new state of player
        // - if not Academy, write playerId in target team's shirt
        if (teamIdTarget == ACADEMY_TEAM) {
            state = setCurrentTeamId(state, ACADEMY_TEAM);
        } else {
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
        
    function transferTeam(uint256 teamId, address addr) private {
        // requiring that team is not bot already ensures that tz and countryIdxInTz exist 
        require(!isBotTeam(teamId), "cannot transfer a bot team");
        require(addr != NULL_ADDR, "cannot transfer to a null address");
        require(teamIdToOwner[teamId] != addr, "buyer and seller are the same addr");
        teamIdToOwner[teamId] = addr;
        emit TeamTransfer(teamId, addr);
    }
    
}
