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

    Assets private _assets;

    mapping (uint256 => uint256) private playerIdToAuctionEnd;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }

    function freezePlayer(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        uint8 typeOfTX,
        bytes32[3] memory sig,
        uint8 sigV
    ) public {
        // check that the purpose of this transaction is of type 1 (sell - agree to buy)
        require(typeOfTX == PUT_FOR_SALE || typeOfTX == MAKE_AN_OFFER, "typeOfTX not valid");

        // check validUntil has not expired
        require(now < validUntil, "these TXs had a valid time that expired already");

        // check player is not already frozen
        require(!isFrozen(playerId), "player already frozen");

        // check asset is owned by someone
        require(_assets.getOwnerPlayer(playerId) != address(0), "player not owned by anyone");

        // check signatures are valid by requiring that they own the asset:
        require(_assets.getOwnerPlayer(playerId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s]),
            "seller is not owner of player, or seller signature is not valid");

        // check that they signed what they input data says they signed:
        // ...for the seller and the buyer:
        bytes32 sellerTxHash;
        sellerTxHash = prefixed(buildPutForSaleTxMsg(sellerHiddenPrice, validUntil, playerId, typeOfTX));
        require(sellerTxHash == sig[IDX_MSG], "seller signed a message that does not match the provided pre-hash data");

        // // Freeze player
        playerIdToAuctionEnd[playerId] = validUntil;
        emit PlayerFreeze(playerId, true);
    }

     function completeAuction(
        bytes32 sellerHiddenPrice,
        uint256 validUntil,
        uint256 playerId,
        uint8 typeOfTX,
        bytes32 buyerHiddenPrice,
        uint256 buyerTeamId,
        bytes32[3] memory sig,
        uint8 sigV
     ) public {
        // check asset is owned by buyer
        require(_assets.getOwnerTeam(buyerTeamId) != address(0), "team not owned by anyone");
        // check signatures are valid by requiring that they own the asset:
        require(_assets.getOwnerTeam(buyerTeamId) == recoverAddr(sig[IDX_MSG], sigV, sig[IDX_r], sig[IDX_s]),
            "buyer is not owner of team, or buyer signature not valid");
    
        // make sure that the playerId is the same that was used by the seller to sign
        bytes32 sellerTxHash = prefixed(buildPutForSaleTxMsg(sellerHiddenPrice, validUntil, playerId, typeOfTX));
        // check that they signed what they input data says they signed:
        // ...for the seller and the buyer:
        bytes32 buyerTxHash = prefixed(buildAgreeToBuyTxMsg(sellerTxHash, buyerHiddenPrice, buyerTeamId));
        require(buyerTxHash == sig[IDX_MSG], "buyer signed a message that does not match the provided pre-hash data");
        require(isFrozen(playerId), "player not frozen, nothing to cancel");
        _assets.transferPlayer(playerId, buyerTeamId);
        playerIdToAuctionEnd[playerId] = 1;
        emit PlayerFreeze(playerId, false);
    }
    // this function is not used in the contract. It's only for external helps
    function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) external pure returns (bytes32) {
        return keccak256(abi.encode(currencyId, price, rnd));
    }

    function buildPutForSaleTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint8 typeOfTX) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, typeOfTX));
    }

    function buildOfferToBuyTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint256 buyerTeamId, uint8 typeOfTX) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, buyerTeamId, typeOfTX));
    }

    function buildAgreeToBuyTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerTxHash, buyerHiddenPrice, buyerTeamId));
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
        return playerIdToAuctionEnd[playerId] > now;
    }

}