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
     * @dev Returns whether the given spender can transfer a given token ID
     * team contract can transfer all token ID
     * @param spender address of the spender to query
     * @param tokenId uint256 ID of the token to be transferred
     * @return bool whether the msg.sender is approved for the given token ID,
     *    is an operator of the owner, or is the owner of the token
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
