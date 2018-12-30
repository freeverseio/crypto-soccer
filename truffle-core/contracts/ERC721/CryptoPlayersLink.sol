pragma solidity ^0.4.24;

import "./CryptoPlayersMintable.sol";
import "./CryptoTeamsLink.sol";

/**
 * @title CryptoPlayersLink
 * @dev CryptPlayers team link logic
 */
contract CryptoPlayersLink is CryptoPlayersMintable {
    CryptoTeamsLink private _cryptoTeams;

    /**
     * @dev Tells whether an operator is approved by a given owner
     * It approves the team contract
     * @param owner owner address which you want to query the approval of
     * @param operator operator address which you want to query the approval of
     * @return bool whether the given operator is approved by the given owner
     */
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
