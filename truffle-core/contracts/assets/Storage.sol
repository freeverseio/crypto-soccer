pragma solidity >=0.4.21 <0.6.0;

import "../state/PlayerState.sol";

/**
 * teamId == 0 is invalid and represents the null team
 */
contract Storage {
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

    uint8 constant public PLAYERS_PER_TEAM = 11;
    uint8 constant internal BITS_PER_SKILL = 14;
    uint16 constant internal SKILL_MASK = 0x3fff;
    uint8 constant public NUM_SKILLS = 5;

    mapping(uint256 => uint256) private _playerIdToState;

    /// @dev An array containing the Team struct for all teams in existence.
    /// @dev The ID of each team is actually his index in this array.
    Team[] private teams;

    PlayerState internal _playerState;

    constructor(address playerState) public {
        _playerState = PlayerState(playerState);
        teams.push(Team("_", 0, 0, 0, 0));
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

    function getPlayerPosInTeam(uint256 playerId) public view returns (uint256) {
        require(_playerExists(playerId), "unexistent player");
        uint256 teamId = getPlayerTeam(playerId);
        uint256 state = getPlayerState(playerId);
        return _playerState.getCurrentShirtNum(state);
    }

    function countTeams() public view returns (uint256){
        return teams.length - 1;
    }

    function getTeamName(uint256 teamId) public view returns (string memory) {
        require(_teamExists(teamId), "invalid team id");
        return teams[teamId].name;
    }

    /// this function uses the inverse of the following formula
    /// playerId = playersPerTeam * (teamId -1) + 1 + posInTeam;
    function getPlayerTeam(uint256 playerId) public view returns (uint256) {
        require(_playerExists(playerId), "unexistent player");
        uint256 state = getPlayerState(playerId);
        return _playerState.getCurrentTeamId(state);
    }

    function getPlayerState(uint256 playerId) public view returns (uint256) {
        require(_playerExists(playerId), "unexistent player");
        if (_isVirtual(playerId)) {
            uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM;
            uint256 posInTeam = playerId - PLAYERS_PER_TEAM * (teamId - 1) - 1;
            string memory teamName = getTeamName(teamId);
            uint256 seed = uint256(keccak256(abi.encodePacked(teamName, posInTeam)));
            uint16[5] memory skills = _computeSkills(seed);
            uint16 birth = _computeBirth(seed, block.timestamp);
            return _playerState.playerStateCreate(
                skills[0], // defence,
                skills[1], //uint256 speed,
                skills[2], //uint256 pass,
                skills[3], //uint256 shoot,
                skills[4], //uint256 endurance,
                birth, //uint256 monthOfBirthInUnixTime,
                playerId,
                teamId,
                posInTeam, //uint256 currentShirtNum,
                0, //uint256 prevLeagueId,
                0, //uint256 prevTeamPosInLeague,
                0, //uint256 prevShirtNumInLeague,
                0 //uint256 lastSaleBloc
            );
        }
        else
            return _playerIdToState[playerId];
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
        return teamId != 0 && teamId < teams.length;
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

    /// Compute a random age between 16 and 35
    /// @return monthOfBirth in monthUnixTime
    function _computeBirth(uint256 rnd, uint256 currentTime) internal pure returns (uint16) {
        rnd >>= BITS_PER_SKILL*NUM_SKILLS;
        uint16 seed = uint16(rnd & SKILL_MASK);
        /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
        uint16 age = 16 + (seed % 20);

        /// @dev Convert age to monthOfBirthAfterUnixEpoch.
        /// @dev I leave it this way for clarity, for the time being.
        uint years2secs = 365 * 24 * 3600; // TODO: make it a constant
        uint month2secs = 30 * 24 * 3600; // TODO: make it a constant

        return uint16((currentTime - age * years2secs) / month2secs);
    }    

     /**
     * @dev Compute the pseudorandom skills, sum of the skills is 250
     * @param rnd is a random number
     * @return 5 skills
     */
    function _computeSkills(uint256 rnd) internal pure returns (uint16[NUM_SKILLS] memory) {
        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(rnd & SKILL_MASK);
            rnd >>= BITS_PER_SKILL;
        }

        /// The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < 5; i++) {
            skills[i] = skills[i] % 50;
            excess += skills[i];
        }

        /// At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        uint16 delta = (250 - excess) / 5;
        for (uint8 i = 0; i < 5; i++)
            skills[i] = skills[i] + delta;

        uint16 remainder = (250 - excess) % 5;
        for (uint8 i = 0 ; i < remainder ; i++)
            skills[i]++;

        return skills;
    }
}
