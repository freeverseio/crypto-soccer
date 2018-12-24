pragma solidity ^0.4.24;

import "./CryptoPlayersBase.sol";
import "./CryptoTeamsLink.sol";

contract CryptoPlayersLink is CryptoPlayersBase {
    CryptoTeamsLink private _cryptoTeams;

    function _isApprovedOrOwner(address spender, uint256 tokenId) internal view returns (bool) {
        if (spender == address(_cryptoTeams))
            return true;
        return super._isApprovedOrOwner(spender, tokenId);
    }

    function setTeamsContract(address cryptoTeams) public {
        _cryptoTeams = CryptoTeamsLink(cryptoTeams);
    }

    function transferFrom(address from, address to, uint256 playerId) public {
        super.transferFrom(from, to, playerId);
        setTeam(playerId, 0);
    }
}
