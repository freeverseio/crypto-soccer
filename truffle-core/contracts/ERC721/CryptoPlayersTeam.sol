pragma solidity ^0.4.24;

import "./CryptoPlayersMintable.sol";
import "./CryptoTeamsPlayers.sol";

/**
 * @title CryptoPlayersTeam
 * @dev CryptPlayers team logic
 */
contract CryptoPlayersTeam is CryptoPlayersMintable {
    CryptoTeamsPlayers private _cryptoTeams;
    mapping(uint256 => uint256) _playerTeam;

    /**
     * @dev Transfers the ownership of a given player ID to another address
     * and reset the owner Team ID
     * Usage of this method is discouraged, use `safeTransferFrom` whenever possible
     * Requires the msg sender to be the owner, approved, or operator
     * @param from current owner of the player
     * @param to address to receive the ownership of the given player ID
     * @param playerId uint256 ID of the player to be transferred
    */
    function transferFrom(address from, address to, uint256 playerId) public {
        super.transferFrom(from, to, playerId);
        _setTeam(playerId, 0);
    }

    /**
     * @dev sets team id for existing player
     */
    function _setTeam(uint256 playerId, uint256 teamId) internal {
        require(_exists(playerId));
        _playerTeam[playerId] = teamId;
    }

    /**
     * @dev returns team id of existing player
     */
    function getTeam(uint256 playerId) public view returns (uint256) {
        require(_exists(playerId));
        return _playerTeam[playerId];
    }
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
        _cryptoTeams = CryptoTeamsPlayers(cryptoTeams);
    }

    function getTeamsContract() external view returns (address) {
        return _cryptoTeams;
    }

    function setTeam(uint256 playerId, uint256 teamId) public onlyTeamsContract {
        _setTeam(playerId, teamId);
    }

    /**
     * @dev Throws if called by any account other that the CyptoTeams.
     */
    modifier onlyTeamsContract() {
        require(msg.sender == address(_cryptoTeams));
        _;
    }
}
