pragma solidity >=0.4.21 <0.6.0;

import "./Leagues.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract Market {
    uint8 constant internal SELL_MSG = 0;
    uint8 constant internal SELL_r   = 1;
    uint8 constant internal SELL_s   = 2;
    uint8 constant internal BUY_MSG  = 3;
    uint8 constant internal BUY_r    = 4;
    uint8 constant internal BUY_s    = 5;
    uint8 constant internal SELLER   = 0;
    uint8 constant internal BUYER    = 1;
    uint8 constant internal PUT_FOR_SALE  = 1;
    uint8 constant internal MAKE_AN_OFFER = 2;

    Assets private _assets;

    mapping (uint256 => uint256) private playerIdToTargetTeam;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }

    function freezePlayer(
        bytes32 privHash,
        uint256 validUntil,
        uint256 playerId,
        uint8 typeOfTX,
        uint256 buyerTeamId,
        bytes32[6] memory sigs,
        uint8[2] memory vs
    ) public {
        // check that the purpose of this transaction is of type 1 (sell - agree to buy)
        require(typeOfTX == PUT_FOR_SALE || typeOfTX == MAKE_AN_OFFER, "typeOfTX not valid");

        // check validUntil has not expired
        require(now < validUntil, "these TXs had a valid time that expired already");

        // check player is not already frozen
        require(!isFrozen(playerId), "player already frozen");

        // check assets are owned by someone
        require(_assets.getOwnerPlayer(playerId) != address(0), "player not owned by anyone");
        require(_assets.getOwnerTeam(buyerTeamId) != address(0), "team not owned by anyone");

        // check signatures are valid by requiring that they own the asset:
        require(_assets.getOwnerPlayer(playerId) == recoverAddr(sigs[SELL_MSG], vs[SELLER], sigs[SELL_r], sigs[SELL_s]),
            "seller is not owner of player, or seller signature is not valid");
        require(_assets.getOwnerTeam(buyerTeamId) == recoverAddr(sigs[BUY_MSG], vs[BUYER], sigs[BUY_r], sigs[BUY_s]),
            "buyer is not owner of team, or buyer signature not valid");

        // check that they signed what they input data says they signed:
        // ...for the seller and the buyer:
        bytes32 sellerTxHash;
        bytes32 buyerTxHash;
        if (typeOfTX == PUT_FOR_SALE) {
            sellerTxHash = prefixed(buildPutForSaleTxMsg(privHash, validUntil, playerId, typeOfTX));
            buyerTxHash = prefixed(buildAgreeToBuyTxMsg(sellerTxHash, buyerTeamId));
        } else {
            buyerTxHash = prefixed(buildOfferToBuyTxMsg(privHash, validUntil, playerId, buyerTeamId, typeOfTX));
            sellerTxHash = buyerTxHash;
        }
        require(sellerTxHash == sigs[SELL_MSG], "seller signed a message that does not match the provided pre-hash data");
        require(buyerTxHash == sigs[BUY_MSG], "buyer signed a message that does not match the provided pre-hash data");

        // // Freeze player
        playerIdToTargetTeam[playerId] = buyerTeamId;
    }

    function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) public pure returns (bytes32) {
        return keccak256(abi.encode(currencyId, price, rnd));
    }

    function buildPutForSaleTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint8 typeOfTX) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, typeOfTX));
    }

    function buildOfferToBuyTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint256 buyerTeamId, uint8 typeOfTX) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, buyerTeamId, typeOfTX));
    }

    function buildAgreeToBuyTxMsg(bytes32 sellerMsg, uint256 buyerTeamId) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerMsg, buyerTeamId));
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
        return playerIdToTargetTeam[playerId] != 0;
    }

    function cancelFreeze(uint256 playerId) public {
        require(isFrozen(playerId), "player not frozen, nothing to cancel");
        delete(playerIdToTargetTeam[playerId]);
    }

    function completeFreeze(uint256 playerId) public {
        require(isFrozen(playerId), "player not frozen, nothing to cancel");
        _assets.transferPlayer(playerId, playerIdToTargetTeam[playerId]);
    }

}