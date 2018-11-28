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


    // Team[] private teams;a
    uint256 private teamsCount = 1;

    /// @dev A mapping from hash(playerName) to a Team struct.
    /// @dev Facilitates checking if a playerName already exists.
    mapping(bytes32 => uint256) private playerToTeam;

    /// @dev A mapping from team hash(name) to the owner's address.
    /// @dev Facilitates checking if a teamName already exists.
    mapping(bytes32 => uint256) private teamToOwnerAddr;
    

    /// @dev Upong deployment of the game, we create the first null player
    /// @dev Choose a silly serialized state (meaningless age, skills, etc)
    /// @dev to differentiate it from 0.
    constructor() public {
        players.push(Player({name: "_", state: uint(-1) }));
    }
    
    function getTeamOwner(bytes32 teamHashName) public view returns(address){
        uint256 teamIdx = teamToOwnerAddr[teamHashName];
        return teams[teamIdx].owner;
    }

    function teamOwnerOf(uint256 _tokenId) external view returns (address){
        require(_tokenId != 0);
        address owner = teams[_tokenId].owner;
        require(owner != address(0));
        return owner;
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

    function getNCreatedTeams() public view returns(uint) {
        return teamsCount - 1;
    }

    function getTeamName(uint idx) public view returns(string) { 
        require(idx != 0);
        require(_teamExists(idx));
        return teams[idx].name;
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        require(team != 0);
        teams[team].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 team) public returns (uint256) {
        require(team != 0);
        return teams[team].playersIdx;
    }

    function addTeam(string memory name, address owner) public {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        require(getTeamOwner(nameHash) == 0);

        teams[teamsCount] = Team({name: name, playersIdx: 0, owner: owner});
        teamToOwnerAddr[nameHash] = teamsCount;
        teamsCount++;
    }

    function _teamExists(uint256 idx) internal returns (bool){
        return teams[idx].owner != address(0);
    }
}
