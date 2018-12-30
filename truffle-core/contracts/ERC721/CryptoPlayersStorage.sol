pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

/**
 * @title CryptoPlayersStorage ERC721 Token Standard with its props
 */
contract CryptoPlayersStorage is ERC721, ERC721Enumerable {
    struct Props {
        string name;
        uint256 state;
        uint256 teamId;
    }

    mapping(uint256 => Props) private _playerProps;

    /**
     * @dev Transfers the ownership of a given token ID to another address
     * and reset the owner Team
     * Usage of this method is discouraged, use `safeTransferFrom` whenever possible
     * Requires the msg sender to be the owner, approved, or operator
     * @param from current owner of the token
     * @param to address to receive the ownership of the given token ID
     * @param tokenId uint256 ID of the token to be transferred
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
        _playerProps[playerId].teamId = teamId;
    }

    /**
     * @dev returns team id of existing player
     */
    function getTeam(uint256 playerId) public view returns (uint256) {
        require(_exists(playerId));
        return _playerProps[playerId].teamId;
    }

    /**
     * @dev sets state of existing player
     */
    function _setState(uint256 playerId, uint256 state) internal {
        require(_exists(playerId));
        _playerProps[playerId].state = state;
    }

    /**
     * @dev returns state of existing player
     */
    function getState(uint256 playerId) public view returns(uint) {
        require(_exists(playerId));
        return _playerProps[playerId].state;
    }

    /**
     * @dev returns name of exiting player
     */
    function getName(uint256 playerId) external view returns(string) {
        require(_exists(playerId));
        return _playerProps[playerId].name;
    }

    /**
     * @dev sets name of existing player
     */
    function _setName(uint256 playerId, string name) internal {
        require(_exists(playerId));
        _playerProps[playerId].name = name;
    }
}
