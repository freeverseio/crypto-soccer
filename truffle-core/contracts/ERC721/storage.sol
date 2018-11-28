pragma solidity ^ 0.4.24;

import "../CryptoSoccer.sol";
import "./CryptoTeams.sol";
/*
    Defines all storage structures and mappings
*/

contract Storage is CryptoSoccer, CryptoTeams {
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


    /// @dev Upong deployment of the game, we create the first null player
    /// @dev Choose a silly serialized state (meaningless age, skills, etc)
    /// @dev to differentiate it from 0.
    constructor() public {
        players.push(Player({name: "_", state: uint(-1) }));
    }
    




    function addPlayer(string memory name, uint state) public {
        players.push(Player({name: name, state: state}));
    }

    function getPlayerState(uint playerIdx) public view returns(uint) {
        return players[playerIdx + 1].state;
    }

    function getPlayerName(uint playerIdx) public view returns(string) {
        return players[playerIdx + 1].name;
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return players.length - 1;
    }

    function teamNameByPlayer(bytes32 playerHashName) public view returns(string){
        uint256 teamIdx = playerToTeam[playerHashName];
        return(teams[teamIdx].name);
    }

    function addPlayerToTeam(bytes32 playerHashName, uint256 idx) public {
        require(idx != 0);
        playerToTeam[playerHashName] = idx;
    }

 

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        require(team != 0);
        teams[team].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 team) public returns (uint256) {
        require(team != 0);
        return teams[team].playersIdx;
    }


    function _teamExists(uint256 idx) internal returns (bool){
        return teams[idx].owner != address(0);
    }
}
