pragma solidity ^0.4.24;

import "./CryptoPlayersMintable.sol";
import "./TeamsContractRole.sol";
import "./CryptoTeamsPlayers.sol";

/**
 * @title CryptoPlayersTeam
 * @dev CryptPlayers team logic
 */
contract CryptoPlayersTeam is CryptoPlayersMintable, TeamsContractRole {
    // Mapping from player ID to team ID
    mapping(uint256 => uint256) _playerTeam;
    
    /**
     * @return team ID of existing player
     */
    function getTeam(uint256 playerId) public view returns (uint256) {
        require(_exists(playerId));
        return _playerTeam[playerId];
    }

    /**
     * @dev set the team of existing player
     */ 
    function setTeam(uint256 playerId, uint256 teamId) public onlyTeamsContract {
        require(_exists(playerId));
        _playerTeam[playerId] = teamId;
    }

    /**
     * @dev Transfers the ownership of a given player to another address
     * The team of the player is reset
     * Requires the msg sender to be the owner, approved, or operator
     * @param from current owner of the player
     * @param to address to receive the ownership of the given player ID
     * @param playerId uint256 ID of the player to be transferred
    */
    function transferFrom(address from, address to, uint256 playerId) public {
        super.transferFrom(from, to, playerId);
        _playerTeam[playerId] = 0;
    }

    /**
     * @dev Returns whether the given spender can transfer a given token ID
     * @param spender address of the spender to query
     * @param tokenId uint256 ID of the token to be transferred
     * @return bool whether the msg.sender is approved for the given token ID,
     *    is an operator of the owner, or is the owner of the token
     */
    function _isApprovedOrOwner(address spender, uint256 tokenId) internal view returns (bool) {
        return isTeamsContract(spender) || super._isApprovedOrOwner(spender, tokenId);
    }
}
