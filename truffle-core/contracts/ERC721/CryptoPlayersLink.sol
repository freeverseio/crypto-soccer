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

    function getTeamsContract() external view returns (address) {
        return _cryptoTeams;
    }

    function setTeam(uint256 playerId, uint256 teamId) public {
        require(msg.sender == address(_cryptoTeams));
        _setTeam(playerId, teamId);
    }
}
