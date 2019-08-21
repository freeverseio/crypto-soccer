pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";

contract FreezableAssets is Assets {
    uint8 constant internal SELL_MSG = 0;
    uint8 constant internal SELL_r   = 1;
    uint8 constant internal SELL_s   = 2;
    uint8 constant internal BUY_MSG  = 3;
    uint8 constant internal BUY_r    = 4;
    uint8 constant internal BUY_s    = 5;
    uint8 constant internal SELLER   = 0;
    uint8 constant internal BUYER    = 1;

    mapping (uint256 => bool) private isPlayerFrozen;

    function freezePlayer(
        bytes32 privHash,
        uint256 validUntil,
        uint256 playerId,
        uint8 typeOfTX,
        uint256 teamId,
        bytes32[6] memory sigs,
        uint8[2] memory vs
    ) public {
        // check that the purpose of this transaction is of type 1 (sell - agree to buy)
        require(typeOfTX == 1, "typeOfTX not valid");

        // check validUntil has not expired
        require(now < validUntil, "these TXs had a valid time that expired already");

        // check player is not already frozen
        require(isPlayerFrozen[playerId] == false, "player already frozen");

        // check assets are owned by someone
        require(getPlayerOwner(playerId) != address(0), "player not owned by anyone");
        require(getTeamOwner(teamId) != address(0), "team not owned by anyone");

        // check signatures are valid by requiring that they own the asset:
        require(getPlayerOwner(playerId) == recoverAddr(sigs[SELL_MSG], vs[SELLER], sigs[SELL_r], sigs[SELL_s]),
            "seller is not owner of player, or seller signature is not valid");
        require(getTeamOwner(teamId) == recoverAddr(sigs[BUY_MSG], vs[BUYER], sigs[BUY_r], sigs[BUY_s]),
            "buyer is not owner of team, or buyer signature not valid");

        // check that they signed what they input data says they signed:
        // ...for the seller:
        bytes32 sellerTxHash = prefixed(buildSellerTxMsg(privHash, validUntil, playerId, typeOfTX));
        require(sellerTxHash == sigs[SELL_MSG], "seller signed a message that does not match the provided pre-hash data");
        // ...for the buyer:
        bytes32 buyerTxHash = prefixed(buildBuyerTxMsg(sellerTxHash, teamId));
        require(buyerTxHash == sigs[BUY_MSG], "buyer signed a message that does not match the provided pre-hash data");

        // // Freeze player
        isPlayerFrozen[playerId] = true;
    }

    function buildSellerTxMsg(bytes32 privHash, uint256 validUntil, uint256 playerId, uint8 typeOfTX) public pure returns (bytes32) {
        return keccak256(abi.encode(privHash, validUntil, playerId, typeOfTX));
    }

    function buildBuyerTxMsg(bytes32 sellerMsg, uint256 teamId) public pure returns (bytes32) {
        return keccak256(abi.encode(sellerMsg, teamId));
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
    function prefixed(bytes32 hash) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", hash));
    }

    function getTeamOwner(uint256 teamId) public view returns (address) {
        string memory name = getTeamName(teamId);
        return getTeamOwner(name);
    }

    function getPlayerOwner(uint256 playerId) public view returns (address) {
        uint256 state = getPlayerState(playerId);
        uint256 teamId = _playerState.getCurrentTeamId(state);
        return getTeamOwner(teamId);
    }

    function transferPlayer(uint256 playerId, uint256 teamIdTarget) public  {
        _transferPlayer(playerId, teamIdTarget);
        isPlayerFrozen[playerId] = false;
    }

    function isFrozen(uint256 playerId) external view returns (bool) {
        require(_playerExists(playerId), "unexistent player");
        return isPlayerFrozen[playerId];
    }
}