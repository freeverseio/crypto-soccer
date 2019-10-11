pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";

/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract Market {
    event PlayerFreeze(uint256 playerId, bool frozen);

    uint8 constant internal IDX_MSG = 0;
    uint8 constant internal IDX_r   = 1;
    uint8 constant internal IDX_s   = 2;
    uint8 constant internal PUT_FOR_SALE  = 1;
    uint8 constant internal MAKE_AN_OFFER = 2;
    // POST_AUCTION_TIME: is how long does the buyer have to pay in fiat, after auction is finished.
    //  ...it includes time to ask for a 2nd-best bidder, or 3rd-best.
    uint256 constant public POST_AUCTION_TIME   = 6 hours; 
    uint256 constant public AUCTION_TIME        = 24 hours; 

    Assets private _assets;

    mapping (uint256 => uint256) private playerIdToAuctionEnd;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
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
            (!isFrozen(playerId))) &&  
            // check asset is owned by legit address
            (_assets.getOwnerPlayer(playerId) != address(0)) && 
            // check signatures are valid by requiring that they own the asset:
            (_assets.getOwnerPlayer(playerId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&    
            // check that they signed what they input data says they signed:
            (sig[IDX_MSG] == prefixed(buildPutForSaleTxMsg(sellerHiddenPrice, validUntil, playerId))
        );
    }

    function freezePlayer(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32[3] memory sig,
        uint8 sigV
    ) public {
        require(areFreezePlayerRequirementsOK(sellerHiddenPrice, validUntil, playerId, sig, sigV), "FreePlayer requirements not met");
        // // Freeze player
        playerIdToAuctionEnd[playerId] = validUntil;
        emit PlayerFreeze(playerId, true);
    }


    function areCompleteAuctionRequirementsOK(
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
        bytes32 sellerTxHash = prefixed(buildPutForSaleTxMsg(sellerHiddenPrice, validUntil, playerId));
        ok =    // check asset is owned by buyer
                (_assets.getOwnerTeam(buyerTeamId) != address(0)) && 
                // check signatures are valid by requiring that they own the asset:
                (_assets.getOwnerTeam(buyerTeamId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s])) &&
                // check player is still frozen
                isFrozen(playerId) &&
                // check that they signed what they input data says they signed:
                sig[IDX_MSG] == prefixed(buildAgreeToBuyTxMsg(sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction));

        if (isOffer2StartAuction) {
            // in this case: validUntil is interpreted as offerValidUntil
            ok = ok && (validUntil > playerIdToAuctionEnd[playerId] - AUCTION_TIME);
        } else {
            ok = ok && (validUntil == playerIdToAuctionEnd[playerId]);
        } 
    }

    function completeAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        bytes32 buyerHiddenPrice,
        uint256 buyerTeamId,
        bytes32[3] memory sig,
        uint8 sigV,
        bool isOffer2StartAuction
     ) public {
        require(areCompleteAuctionRequirementsOK(
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
        playerIdToAuctionEnd[playerId] = 1;
        emit PlayerFreeze(playerId, false);
    }
    
    
    // this function is not used in the contract. It's only for external helps
    function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) external pure returns (bytes32) {
        return keccak256(abi.encode(currencyId, price, rnd));
    }

    function buildPutForSaleTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId));
    }

    function buildOfferToBuyTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint256 buyerTeamId) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, buyerTeamId));
    }

    function buildAgreeToBuyTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bool isOffer2StartAuction) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction));
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

    function isFrozen(uint256 playerId) public view returns (bool) {
        require(_assets.playerExists(playerId), "unexistent player");
        return playerIdToAuctionEnd[playerId] + POST_AUCTION_TIME > now;
    }

}