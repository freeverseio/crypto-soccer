pragma solidity >=0.5.12 <=0.6.3;

import "./Proxy.sol";
import "./Market.sol";

/**
 * @title Auctions operated in cryptocurrency, without anyone's permission (other than sellers and buyers)
 */
 
contract MarketCrypto {

    uint256 internal _auctionDuration = 24 hours; 
    Proxy private _proxy;
    Market private _market;

    mapping (uint256 => uint256) internal _putForSalePlayerIdToAuctionData;

    function setProxyAddress(address addr) external {
        _proxy = Proxy(addr);
    }

    function setMarketAddress(address addr) external {
        _market = Market(addr);
    }
    
    function setActionDuration(uint256 newDuration) external {
        _auctionDuration = newDuration;
    }

    function putPlayerForSale(uint256 playerId, uint256 price) external {
        uint256 currentTeamId  = _market.getCurrentTeamIdFromPlayerId(playerId);
        address currentOwner   = _market.getOwnerTeam(currentTeamId);
        bool OK = (
            (price < 2**205) &&
            // check player is not already frozen
            (!_market.isPlayerFrozen(playerId)) &&  
            // check that the team it belongs to is not already frozen
            !_market.isTeamFrozen(currentTeamId) &&
            // check asset has not been put for sale before
            _putForSalePlayerIdToAuctionData[playerId] == 0 &&
            // check asset is owned by legit address
            (currentOwner != address(0)) && 
            // check asset is owned by sender of this TX
            (currentOwner == msg.sender)   
        );
        require(OK, "conditions to putPlayerForSale not met");
        _putForSalePlayerIdToAuctionData[playerId] = now + _auctionDuration + ((uint256(price) << 40) >> 8);
    }
    
    function isPlayerPutForSale(uint256 playerId) public view returns (bool) {
        return 
    }
    
    // // Main PLAYER auction functions: freeze & complete
    // function freezePlayer(
    //     bytes32 sellerHiddenPrice,
    //     uint256 validUntil,
    //     uint256 playerId,
    //     bytes32[3] memory sig,
    //     uint8 sigV
    // ) public {
    //     require(areFreezePlayerRequirementsOK(sellerHiddenPrice, validUntil, playerId, sig, sigV), "FreePlayer requirements not met");
    //     // // Freeze player
    //     _playerIdToAuctionData[playerId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
    //     emit PlayerFreeze(playerId, _playerIdToAuctionData[playerId], true);
    // }

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
    //     uint256 playerIdWithoutTargetTeam = setTargetTeamId(playerId, 0);
    //     require(!isPlayerWritten(playerIdWithoutTargetTeam), "promo player already in the universe");
    //     uint256 targetTeamId = getTargetTeamId(playerId);
    //     // require that team does not have any constraint from friendlies
    //     (bool isConstrained, uint8 nRemain) = getMaxAllowedAcquisitions(targetTeamId);
    //     require(!(isConstrained && (nRemain == 0)), "trying to accept a promo player, but team is busy in constrained friendlies");
    //     // testing about the target team is already done in _assets.transferPlayer
    //     require(teamExists(targetTeamId), "cannot offer a promo player to a non-existent team");
    //     require(!isBotTeam(targetTeamId), "cannot offer a promo player to a bot team");
                
    //     require(getOwnerTeam(targetTeamId) == 
    //                 recoverAddr(sigBuy[IDX_MSG], sigVBuy, sigBuy[IDX_r], sigBuy[IDX_s]), "Buyer does not own targetTeamId");
         
    //     require(_academyAddr == 
    //                 recoverAddr(sigSel[IDX_MSG], sigVSel, sigSel[IDX_r], sigSel[IDX_s]), "Seller does not own academy");
         
    //     bytes32 signedMsg = prefixed(buildPromoPlayerTxMsg(playerId, validUntil));
    //     require(sigBuy[IDX_MSG] == signedMsg, "buyer msg does not match");
    //     require(sigSel[IDX_MSG] == signedMsg, "seller msg does not match");
    //     transferPlayer(playerIdWithoutTargetTeam, targetTeamId);
    //     decreaseMaxAllowedAcquisitions(targetTeamId);
    // }

    // function completePlayerAuction(
    //     bytes32 sellerHiddenPrice,
    //     uint256 validUntil,
    //     uint256 playerId,
    //     bytes32 buyerHiddenPrice,
    //     uint256 buyerTeamId,
    //     bytes32[3] memory sig,
    //     uint8 sigV,
    //     bool isOffer2StartAuction
    //  ) public {
    //     require(areCompletePlayerAuctionRequirementsOK(
    //         sellerHiddenPrice,
    //         validUntil,
    //         playerId,
    //         buyerHiddenPrice,
    //         buyerTeamId,
    //         sig,
    //         sigV,
    //         isOffer2StartAuction)
    //         , "requirements to complete auction are not met"    
    //     );
    //     transferPlayer(playerId, buyerTeamId);
    //     _playerIdToAuctionData[playerId] = 1;
    //     decreaseMaxAllowedAcquisitions(buyerTeamId);
    //     emit PlayerFreeze(playerId, 1, false);
    // }
    
    // // Main TEAM auction functions: freeze & complete
    // function freezeTeam(
    //     bytes32 sellerHiddenPrice,
    //     uint256 validUntil,
    //     uint256 teamId,
    //     bytes32[3] memory sig,
    //     uint8 sigV
    // ) public {
    //     require(areFreezeTeamRequirementsOK(sellerHiddenPrice, validUntil, teamId, sig, sigV), "FreePlayer requirements not met");
    //     // // Freeze player
    //     _teamIdToAuctionData[teamId] = validUntil + ((uint256(sellerHiddenPrice) << 40) >> 8);
    //     emit TeamFreeze(teamId, _teamIdToAuctionData[teamId], true);
    // }

    // function completeTeamAuction(
    //     bytes32 sellerHiddenPrice,
    //     uint256 validUntil,
    //     uint256 teamId,
    //     bytes32 buyerHiddenPrice,
    //     bytes32[3] memory sig,
    //     uint8 sigV,
    //     bool isOffer2StartAuction
    //  ) public {
    //     (bool ok, address buyerAddress) = areCompleteTeamAuctionRequirementsOK(
    //         sellerHiddenPrice,
    //         validUntil,
    //         teamId,
    //         buyerHiddenPrice,
    //         sig,
    //         sigV,
    //         isOffer2StartAuction
    //     );
    //     require(ok, "requirements to complete auction are not met");
    //     transferTeam(teamId, buyerAddress);
    //     _teamIdToAuctionData[teamId] = 1;
    //     emit TeamFreeze(teamId, 1, false);
    // }

    // function transferPlayer(uint256 playerId, uint256 teamIdTarget) public  {
    //     // warning: check of ownership of players and teams should be done before calling this function
    //     // TODO: checking if they are bots should be done outside this function
    //     require(getIsSpecial(playerId) || playerExists(playerId), "player does not exist");
    //     require(teamExists(teamIdTarget), "unexistent target team");
    //     uint256 state = getPlayerState(playerId);
    //     uint256 newState = state;
    //     uint256 teamIdOrigin = getCurrentTeamId(state);
    //     require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
    //     require(!isBotTeam(teamIdOrigin) && !isBotTeam(teamIdTarget), "cannot transfer player when at least one team is a bot");
    //     uint8 shirtTarget = getFreeShirt(teamIdTarget);
    //     require(shirtTarget != PLAYERS_PER_TEAM_MAX, "target team for transfer is already full");
        
    //     _playerIdToState[playerId] = 
    //         setLastSaleBlock(
    //             setCurrentShirtNum(
    //                 setCurrentTeamId(
    //                     newState, teamIdTarget
    //                 ), shirtTarget
    //             ), block.number
    //         );

    //     teamIdToPlayerIds[teamIdTarget][shirtTarget] = playerId;
    //     if (teamIdOrigin != ACADEMY_TEAM) {
    //         uint256 shirtOrigin = getCurrentShirtNum(state);
    //         teamIdToPlayerIds[teamIdOrigin][shirtOrigin] = FREE_PLAYER_ID;
    //     }
    //     emit PlayerStateChange(playerId, newState);
    // }
    
        
    // function transferTeam(uint256 teamId, address addr) public {
    //     // requiring that team is not bot already ensures that tz and countryIdxInTz exist 
    //     require(!isBotTeam(teamId), "cannot transfer a bot team");
    //     require(addr != NULL_ADDR, "cannot transfer to a null address");
    //     require(teamIdToOwner[teamId] != addr, "buyer and seller are the same addr");
    //     teamIdToOwner[teamId] = addr;
    //     emit TeamTransfer(teamId, addr);
    // }
}
