pragma solidity >=0.4.21 <0.6.0;

import "../state/PlayerState.sol";

/**
 * teamId == 0 is invalid and represents the null team
 */
contract Storage {
    uint8 constant public PLAYERS_PER_TEAM = 11;

    mapping(uint256 => uint256) private _playerIdToState;

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
        uint256 currentLeagueId;
        uint8 posInCurrentLeague;
        uint256 prevLeagueId;
        uint8 posInPrevLeague;
    }

    /// @dev An array containing the Team struct for all teams in existence.
    /// @dev The ID of each team is actually his index in this array.
    Team[] private teams;

    PlayerState internal _playerState;

    constructor(address playerState) public {
        _playerState = PlayerState(playerState);
        teams.push(Team("_", 0, 0, 0, 0));
    }

    function _setPlayerState(uint256 state) internal {
        uint256 playerId = _playerState.getPlayerId(state);
        require(_playerExists(playerId), "unexistent player");
        uint256 teamId = _playerState.getCurrentTeamId(state);
        require(_teamExists(teamId), "unexistent team");
        uint256 shirtNumber = _playerState.getCurrentShirtNum(state);
        require(shirtNumber < PLAYERS_PER_TEAM, "invalid shirt number");
        shirtNumber = _playerState.getPrevShirtNumInLeague(state);
        require(shirtNumber < PLAYERS_PER_TEAM, "invalid shirt number");
        uint256 saleBlock = _playerState.getLastSaleBlock(state);
        require(saleBlock != 0 && saleBlock <= block.number, "invalid sale block");
        _playerIdToState[playerId] = state;
    }

    function getPlayerState(uint256 playerId) public view returns (uint256) {
        require(_playerExists(playerId), "unexistent player");
        if (_isVirtual(playerId)) {
            // uint256[5] skills = 
            return 0;
        }
        else
            return _playerIdToState[playerId];
    }

    function countTeams() public view returns (uint256){
        return teams.length - 1;
    }

    function getTeamName(uint256 teamId) public view returns (string memory) {
        require(_teamExists(teamId), "invalid team id");
        return teams[teamId].name;
    }

    function getTeamCurrentHistory(uint256 teamId) external view returns (
        uint256 currentLeagueId,
        uint8 posInCurrentLeague,
        uint256 prevLeagueId,
        uint8 posInPrevLeague
        )
    {
        require(_teamExists(teamId), "invalid team id");
        return (
            teams[teamId].currentLeagueId,
            teams[teamId].posInCurrentLeague,
            teams[teamId].prevLeagueId,
            teams[teamId].posInPrevLeague);
    }

    function _updateTeamCurrentHistory(
        uint256 teamId,
        uint256 currentLeagueId,
        uint8 posInCurrentLeague
    )
    internal
    {
        require(_teamExists(teamId), "invalid team id");
        teams[teamId].prevLeagueId = teams[teamId].currentLeagueId;
        teams[teamId].posInPrevLeague = teams[teamId].posInCurrentLeague;
        teams[teamId].currentLeagueId = currentLeagueId;
        teams[teamId].posInCurrentLeague = posInCurrentLeague;
    }


    function _addTeam(string memory name) internal returns (uint256) {
        teams.push(Team(name, 0, 0, 0, 0));
        return teams.length - 1;
    }

    function _teamExists(uint256 teamId) internal view returns (bool) {
        return teamId != 0 && teamId <= teams.length;
    }

    function _playerExists(uint256 playerId) internal view returns (bool) {
        if (playerId == 0) return false;
        if (_playerIdToState[playerId] != 0) return true;
        uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM;
        return teamId <= countTeams();
    }

    function _isVirtual(uint256 playerId) internal view returns (bool) {
        require(_playerExists(playerId), "unexistent player");
        return _playerIdToState[playerId] == 0;
    }
}
