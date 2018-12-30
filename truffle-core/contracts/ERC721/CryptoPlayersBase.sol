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

    mapping(uint256 => Props) private _playerProps;

    function transferFrom(address from, address to, uint256 playerId) public {
        super.transferFrom(from, to, playerId);
        _setTeam(playerId, 0);
    }

    function _setTeam(uint256 playerId, uint256 teamId) internal {
        require(_exists(playerId));
        _playerProps[playerId].teamId = teamId;
    }

    function getTeam(uint256 playerId) public view returns (uint256) {
        require(_exists(playerId));
        return _playerProps[playerId].teamId;
    }

    function _setState(uint256 playerId, uint256 state) internal {
        require(_exists(playerId));
        _playerProps[playerId].state = state;
    }

    function getState(uint256 playerId) public view returns(uint) {
        require(_exists(playerId));
        return _playerProps[playerId].state;
    }

    function getName(uint256 playerId) external view returns(string) {
        require(_exists(playerId));
        return _playerProps[playerId].name;
    }

    function _setName(uint256 playerId, string name) internal {
        require(_exists(playerId));
        _playerProps[playerId].name = name;
    }
}
