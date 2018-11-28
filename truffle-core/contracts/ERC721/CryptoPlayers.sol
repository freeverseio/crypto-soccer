pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";

contract CryptoPlayers is ERC721Full("CryptoSoccerPlayers", "CSP") {
    constructor() public {
        players.push(Player({name: "_", state: uint(-1) }));
    }

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
    struct Player {
        string name;
        uint state;
    }

    /// @dev An array containing the Player struct for all players in existence. 
    /// @dev The ID of each player is actually his index this array.
    Player[] private players;

    /// @dev A mapping from hash(playerName) to a Team struct.
    /// @dev Facilitates checking if a playerName already exists.
    mapping(bytes32 => uint256) private playerToTeam;

    function addPlayer(string memory name, uint state, uint256 teamIdx) public {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        players.push(Player({name: name, state: state}));
        playerToTeam[playerNameHash] = teamIdx;
    }

    function getPlayerState(uint playerIdx) public view returns(uint) {
        return players[playerIdx + 1].state;
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return players.length - 1;
    }

    function getPlayerName(uint playerIdx) public view returns(string) {
        return players[playerIdx + 1].name;
    }

    function getTeamIndexByPlayer(string name) public view returns (uint256){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        return playerToTeam[playerNameHash];
    }

    function playerExists(string name) public view returns (bool){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 teamIdx = playerToTeam[playerNameHash];
        return teamIdx != 0;
    }
}
