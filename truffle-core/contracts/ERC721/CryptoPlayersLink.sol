pragma solidity ^0.4.24;

import "./CryptoPlayersBase.sol";

contract CryptoPlayersLink is CryptoPlayersBase {
    function transferFrom(address from, address to, uint256 tokenId) public {
        super.transferFrom(from, to, tokenId);
        _setTeam(tokenId, 0);
    }
}