pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

contract CryptoPlayersBase is ERC721, ERC721Enumerable {
    /// @dev The main Player struct.
    /// @dev name is a string, unique for every Player
    /// @dev state is a uint256 that serializes age, skills, role.
    /// @dev Each skill is sent as a uint16 for serialization, occupying 20 bits of the state
    /// @dev The order of elements serialized are:
    ///         0-monthOfBirthAfterUnixEpoch; if this goes up to 9999, then the game will run for at least 800 more years.
    ///         1-defense
    ///         2-speed
    ///         3-pass
    ///         4-shoot (for a goalkeeper, this is interpreted as ability to block a shoot)
    ///         5-endurance
    ///         6-role
    struct Props {
        string name;
        uint256 state;
        uint256 teamId;
    }

    /// @dev An array containing the Player struct for all players in existence. 
    /// @dev The ID of each player is actually his index this array.
    mapping(uint256 => Props) private _playerProps;

    /// @dev A mapping from hash(playerName) to a Team struct.
    /// @dev Facilitates checking if a playerName already exists.
    mapping(bytes32 => uint256) private _nameHashPlayer;

    function _addPlayer(string memory name, uint state, uint256 teamIdx, address owner) internal {
        require(teamIdx != 0);
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 playerId = totalSupply() + 1;
        _mint(owner, playerId);
        _playerProps[playerId].name = name;
        _playerProps[playerId].state = state;
        _playerProps[playerId].teamId = teamIdx;
        _nameHashPlayer[playerNameHash] = playerId;
    }

    function _setState(uint256 playerId, uint256 state) internal {
        require(_exists(playerId));
        _playerProps[playerId].state = state;
    }

    function _getState(uint playerId) internal view returns(uint) {
        require(_exists(playerId));
        return _playerProps[playerId].state;
    }

    function _getName(uint playerId) internal view returns(string) {
        require(_exists(playerId));
        return _playerProps[playerId].name;
    }

    function _getTeamIndexByPlayer(string name) internal view returns (uint256){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 id = _nameHashPlayer[playerNameHash];
        require(id != 0);
        return _playerProps[id].teamId;
    }

    function _playerExists(string name) internal view returns (bool){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 teamIdx = _nameHashPlayer[playerNameHash];
        return teamIdx != 0;
    }
}
