pragma solidity >=0.4.21 <0.6.0;

import "../state/PlayerState.sol";
 
/// teamId == 0 is invalid and represents the null team
/// TODO: fix the playerPos <=> playerShirt doubt
contract Assets {
    event TeamCreated (uint256 id);
    event TeamTransfer(uint256 teamId, address to);

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        uint256[PLAYERS_PER_TEAM_MAX] playerIds;
        address owner; // timestamp as seconds since unix epoch
    }

    struct Country {
        uint256 nDivisions;
        uint8 nDivisionsToAddNextRound;
        mapping (uint256 => uint256) divisonIdxToRound;
        mapping (uint256 => Team) teamIdxInCountryToTeam;
    }

    struct TimeZone {
        Country[] countries;
        uint8 nCountriesToAdd;
        uint256[2] orgMapHash;
        uint256[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        uint256 scoresRoot;
        uint8 updateCycleIdx;
        uint256 lastUpdateBlockNum;
        uint256 actionsRoot;
        uint256 blockHash;
        uint256 lastMarketClosureBlockNum;
    }    
    
    uint8 constant public PLAYERS_PER_TEAM_INIT = 18;
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint256 constant public FREE_PLAYER_ID  = uint256(-1);
    uint8 constant internal BITS_PER_SKILL = 14;
    uint16 constant internal SKILL_MASK = 0x3fff;
    uint8 constant public N_SKILLS = 5;
    uint8 constant public LEAGUES_PER_DIV = 16;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    address constant public FREEVERSE = address(1);

    mapping(uint256 => uint256) private _playerIdToState;

    PlayerState internal _playerStateLib;
    TimeZone[] internal _timeZones;
    mapping (uint256 => uint256) internal _playerIdxToPlayerState;

    constructor(address playerState) public {
        _playerStateLib = PlayerState(playerState);
        // the first timeZone is a dummy one, without any country. Forbidden to use timeZone[0].
        _timeZones.length++;
        // It then creates the remaining 24 timezones.
        for (uint8 tz = 0; tz < 24; tz++) {
            createTimeZone();
        }
    }

    function createTimeZone() private {
        _timeZones.length++;
        TimeZone storage tz = _timeZones[_timeZones.length - 1];
        Country memory country;
        country.nDivisions = 1;
        tz.countries.push(country);
        tz.countries[tz.countries.length -1].divisonIdxToRound[0] = 1; 
    }
        
    function getNCountriesInTZ(uint8 timeZone) public view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].countries.length;
    }
        
    function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return _timeZones[timeZone].countries[countryIdxInTZ].nDivisions;
    }

    function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNDivisionsInCountry(timeZone, countryIdxInTZ) * LEAGUES_PER_DIV;
    }

    function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNLeaguesInCountry(timeZone, countryIdxInTZ) * TEAMS_PER_LEAGUE;
    }

    function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return (teamIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ));
    }

    function teamExists(uint256 teamId) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        return _teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry) == address(0);
    }

    function isBotTeam(uint256 teamId) public view returns(bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        return isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(address) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner;
    }

    function getOwnerTeam(uint256 teamId) public view returns(address) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function _wasPlayerCreatedInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) private view returns(bool) {
        return (playerIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ) * PLAYERS_PER_TEAM_INIT);
    }

    function _assertTZExists(uint8 timeZone) private pure {
        require(timeZone > 0 && timeZone < 25, "timeZone does not exist");
    }

    function _assertCountryInTZExists(uint8 timeZone, uint256 countryIdxInTZ) private view {
        require(countryIdxInTZ < _timeZones[timeZone].countries.length, "country does not exist in this timeZone");
    }

    function playerExists(uint256 playerId) public view returns (bool) {
        if (playerId == 0) return false;
        if (_playerIdToState[playerId] != 0) return true;
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(playerId);
        return _wasPlayerCreatedInCountry(timeZone, countryIdxInTZ, playerIdxInCountry);
    }

    function isVirtual(uint256 playerId) public view returns (bool) {
        require(playerExists(playerId), "unexistent player");
        return _playerIdToState[playerId] == 0;
    }

    function transferBotInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) public {
        require(isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a non-bot team");
        require(addr != address(0));
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = addr;
    }

    function transferBotToAddr(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        transferBotInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
    }


    // /// @dev Transfers a team to a new owner. 
    // /// @dev This function should be called only when the transfer is legit, as checked elsewhere.
    // function transferTeam(uint256 teamId, address newOwner) public {
    //     require(_teamExists(teamId), "invalid team id");
    //     require(newOwner != address(0), "meaningless adress");
    //     require(newOwner != _getTeamOwner(teams[teamId].name), "unable to transfer between the same user");
    //     bytes32 nameHash = keccak256(abi.encode(teams[teamId].name));
    //     _teamNameHashToOwner[nameHash] = newOwner;
    //     emit TeamTransfer(teamId, newOwner);
    // }
    

    // function getTeamPlayerIds(uint256 teamId) public view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
    //     (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
    //     require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
    //     for (uint8 pos = 0 ; pos < PLAYERS_PER_TEAM_MAX ; pos++){
    //         if (_timeZones[timeZone].countries[teamId]. teams[teamId].playerIds[pos] == 0) // virtual player
    //             playerIds[pos] = generateVirtualPlayerId(teamId, pos);
    //         else
    //             playerIds[pos] = teams[teamId].playerIds[pos];
    //     }
    // }
        
    // function getTeamCreationTimestamp(uint256 teamId) public view returns (uint256) {
    //     require(_teamExists(teamId), "invalid team id");
    //     return teams[teamId].creationTimestamp;
    // }

    // function getCurrentLeagueId(uint256 teamId) external view returns (uint256) {
    //     require(_teamExists(teamId), "invalid team id");
    //     return teams[teamId].currentLeagueId;
    // }

    // /// get the current and previous team league and position in league
    // function getTeamCurrentHistory(uint256 teamId) external view returns (
    //     uint256 currentLeagueId,
    //     uint8 posInCurrentLeague,
    //     uint256 prevLeagueId,
    //     uint8 posInPrevLeague
    //     )
    // {
    //     require(_teamExists(teamId), "invalid team id");
    //     return (
    //         teams[teamId].currentLeagueId,
    //         teams[teamId].posInCurrentLeague,
    //         teams[teamId].prevLeagueId,
    //         teams[teamId].posInPrevLeague);
    // }



    // function getFreeShirt(uint256 teamId) public view returns(uint8) {
    //     for (uint8 shirtNum = PLAYERS_PER_TEAM_MAX-1; shirtNum >= 0; shirtNum--) {
    //         if (isFreeShirt(teamId, shirtNum)) {
    //             return shirtNum;
    //         }
    //     }
    //     return PLAYERS_PER_TEAM_MAX;
    // }

    // function _transferPlayer(uint256 playerId, uint256 teamIdTarget) internal  {
    //     // warning: check of ownership of players and teams should be done before calling this function
    //     require(_playerExists(playerId) && _teamExists(teamIdTarget), "unexistent player or team");
    //     uint256 state = getPlayerState(playerId);
    //     uint256 newState = state;
    //     uint256 teamIdOrigin = _playerStateLib.getCurrentTeamId(state);
    //     require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
    //     uint256 shirtOrigin = _playerStateLib.getCurrentShirtNum(state);
    //     uint8 shirtTarget = getFreeShirt(teamIdTarget);
    //     require(shirtTarget != PLAYERS_PER_TEAM_MAX, "target team for transfer is already full");
        
    //     newState = _playerStateLib.setCurrentTeamId(newState, teamIdTarget);
    //     newState = _playerStateLib.setCurrentShirtNum(newState, shirtTarget);
    //     newState = _playerStateLib.setLastSaleBlock(newState, block.number);

    //     teams[teamIdTarget].playerIds[shirtTarget] = playerId;
    //     teams[teamIdOrigin].playerIds[shirtOrigin] = FREE_PLAYER_ID;

    //     _setPlayerState(newState);
    // }



    // // TODO: exchange fails on playerId0 & playerId1 of the same team
    // function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) public {
    //     // TODO: check ownership address
    //     require(_playerExists(playerId0) && _playerExists(playerId1), "unexistent playerId");
    //     uint256 state0 = getPlayerState(playerId0);
    //     uint256 state1 = getPlayerState(playerId1);
    //     uint256 newState0 = state0;
    //     uint256 teamId0 = _playerStateLib.getCurrentTeamId(state0);
    //     uint256 teamId1 = _playerStateLib.getCurrentTeamId(state1);
    //     uint256 playerShirt0 = _playerStateLib.getCurrentShirtNum(state0);
    //     uint256 playerShirt1 = _playerStateLib.getCurrentShirtNum(state1);
    //     newState0 = _playerStateLib.setCurrentTeamId(newState0, _playerStateLib.getCurrentTeamId(state1));
    //     newState0 = _playerStateLib.setCurrentShirtNum(newState0, _playerStateLib.getCurrentShirtNum(state1));
    //     state1 = _playerStateLib.setCurrentTeamId(state1,_playerStateLib.getCurrentTeamId(state0));
    //     state1 = _playerStateLib.setCurrentShirtNum(state1,_playerStateLib.getCurrentShirtNum(state0));
    //     newState0 = _playerStateLib.setLastSaleBlock(newState0, block.number);
    //     state1 = _playerStateLib.setLastSaleBlock(state1, block.number);

    //     teams[teamId0].playerIds[playerShirt0] = playerId1;
    //     teams[teamId1].playerIds[playerShirt1] = playerId0;

    //     // TODO
    //     // if getBlockNumForLastLeagueOfTeam(teamIdx1, ST) > state1.getLastSaleBlocknum():
    //     //     state1.prevLeagueIdx = ST.teams[teamIdx1].currentLeagueIdx
    //     //     state1.prevTeamPosInLeague = ST.teams[teamIdx1].teamPosInCurrentLeague

    //     // if getBlockNumForLastLeagueOfTeam(teamIdx2, ST) > state2.getLastSaleBlocknum():
    //     //     state2.prevLeagueIdx = ST.teams[teamIdx2].currentLeagueIdx
    //     //     state2.prevTeamPosInLeague = ST.teams[teamIdx2].teamPosInCurrentLeague

    //     _setPlayerState(newState0);
    //     _setPlayerState(state1);
    // }

    // function createTeam(string memory name, address owner) public {
    //     bytes32 nameHash = keccak256(abi.encode(name));
    //     require(_teamNameHashToOwner[nameHash] == address(0), "team already exists");
    //     _teamNameHashToOwner[nameHash] = owner;
    //     uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
    //     for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
    //         playerIds[p] = FREE_PLAYER_ID;
    //     }
    //     teams.push(Team(name, 0, 0, 0, 0, playerIds, block.timestamp));
    //     uint256 id = teams.length - 1;
    //     emit TeamCreated(id);
    // }

    // function signToLeague(
    //     uint256 teamId,
    //     uint256 leagueId,
    //     uint8 posInLeague
    // )
    // public
    // {
    //     require(_teamExists(teamId), "invalid team id");
    //     require(teams[teamId].currentLeagueId != leagueId, "cannot sign to a league twice");
    //     teams[teamId].prevLeagueId = teams[teamId].currentLeagueId;
    //     teams[teamId].posInPrevLeague = teams[teamId].posInCurrentLeague;
    //     teams[teamId].currentLeagueId = leagueId;
    //     teams[teamId].posInCurrentLeague = posInLeague;
    // }

    // // TODO: exception when not existent team
    // function _getTeamOwner(string memory name) internal view returns (address) {
    //     bytes32 nameHash = keccak256(abi.encode(name));
    //     return _teamNameHashToOwner[nameHash];
    // }

    // function countTeams() public view returns (uint256){
    //     return teams.length - 1;
    // }

    // function getTeamName(uint256 teamId) public view returns (string memory) {
    //     require(_teamExists(teamId), "invalid team id");
    //     return teams[teamId].name;
    // }


    // function getPlayerState(uint256 playerId) public view returns (uint256) {
    //     require(_playerExists(playerId), "unexistent player");
    //     if (_isVirtual(playerId))
    //         return generateVirtualPlayerState(playerId);
    //     else
    //         return _playerIdToState[playerId];
    // }

    // function generateVirtualPlayerId(uint256 teamId, uint8 posInTeam) public view returns (uint256) {
    //     require(_teamExists(teamId), "unexistent team");
    //     require(posInTeam < PLAYERS_PER_TEAM_MAX, "invalid player pos");
    //     return PLAYERS_PER_TEAM_INIT * (teamId - 1) + 1 + posInTeam;
    // }

    // function generateVirtualPlayerState(uint256 playerId) public view returns (uint256) {
    //         uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM_INIT;
    //         uint256 posInTeam = playerId - PLAYERS_PER_TEAM_INIT * (teamId - 1) - 1;
    //         string memory teamName = getTeamName(teamId);
    //         uint256 seed = _computeSeed(teamName, posInTeam);
    //         uint16[5] memory skills = _computeSkills(seed);
    //         uint16 birth = _computeBirth(seed, getTeamCreationTimestamp(teamId));
    //         return _playerStateLib.playerStateCreate(
    //             skills[0], // defence,
    //             skills[1], // speed,
    //             skills[2], // pass,
    //             skills[3], // shoot,
    //             skills[4], // endurance,
    //             birth, // monthOfBirthInUnixTime,
    //             playerId,
    //             teamId,
    //             posInTeam, // currentShirtNum,
    //             0, // prevLeagueId,
    //             0, // prevTeamPosInLeague,
    //             0, // prevShirtNumInLeague,
    //             0 // lastSaleBloc
    //         );
    // }

    // function _setPlayerState(uint256 state) internal {
    //     uint256 playerId = _playerStateLib.getPlayerId(state);
    //     require(_playerExists(playerId), "unexistent player");
    //     uint256 teamId = _playerStateLib.getCurrentTeamId(state);
    //     require(_teamExists(teamId), "unexistent team");
    //     uint256 shirtNumber = _playerStateLib.getCurrentShirtNum(state);
    //     require(shirtNumber < PLAYERS_PER_TEAM_MAX, "invalid shirt number");
    //     shirtNumber = _playerStateLib.getPrevShirtNumInLeague(state);
    //     require(shirtNumber < PLAYERS_PER_TEAM_MAX, "invalid shirt number");
    //     uint256 saleBlock = _playerStateLib.getLastSaleBlock(state);
    //     require(saleBlock != 0 && saleBlock <= block.number, "invalid sale block");
    //     _playerIdToState[playerId] = state;
    // }


    // function isFreeShirt(uint256 teamId, uint8 shirtNum) public view returns (bool) {
    //     return teams[teamId].playerIds[shirtNum] == FREE_PLAYER_ID;
    // }


    // /// Compute a random age between 16 and 35
    // /// @param rnd is a random number used as seed of the skills
    // /// @param currentTime in seconds since unix epoch
    // /// @return monthOfBirth in monthUnixTime
    // function _computeBirth(uint256 rnd, uint256 currentTime) internal pure returns (uint16) {
    //     rnd >>= BITS_PER_SKILL*N_SKILLS;
    //     uint16 seed = uint16(rnd & SKILL_MASK);
    //     /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
    //     uint16 age = 16 + (seed % 20);

    //     /// @dev Convert age to monthOfBirthAfterUnixEpoch.
    //     /// @dev I leave it this way for clarity, for the time being.
    //     uint years2secs = 365 * 24 * 3600; // TODO: make it a constant
    //     uint month2secs = 30 * 24 * 3600; // TODO: make it a constant

    //     return uint16((currentTime - age * years2secs) / month2secs);
    // }

    // /// Compute the pseudorandom skills, sum of the skills is 250
    // /// @param rnd is a random number used as seed of the skills
    // /// @return 5 skills
    // function _computeSkills(uint256 rnd) internal pure returns (uint16[N_SKILLS] memory) {
    //     uint16[5] memory skills;
    //     for (uint8 i = 0; i<5; i++) {
    //         skills[i] = uint16(rnd & SKILL_MASK);
    //         rnd >>= BITS_PER_SKILL;
    //     }

    //     /// The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
    //     uint16 excess;
    //     for (uint8 i = 0; i < 5; i++) {
    //         skills[i] = skills[i] % 50;
    //         excess += skills[i];
    //     }

    //     /// At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
    //     uint16 delta = (250 - excess) / 5;
    //     for (uint8 i = 0; i < 5; i++)
    //         skills[i] = skills[i] + delta;

    //     uint16 remainder = (250 - excess) % 5;
    //     for (uint8 i = 0 ; i < remainder ; i++)
    //         skills[i]++;

    //     return skills;
    // }

    // /// @return seed
    // function _computeSeed(string memory teamName, uint256 posInTeam) internal pure returns (uint256) {
    //     return uint256(keccak256(abi.encode(teamName, posInTeam)));
    // }

    // /// @return hashed arg casted to uint256
    // function _intHash(string memory arg) internal pure returns (uint256) {
    //     return uint256(keccak256(abi.encode(arg)));
    // }
}
