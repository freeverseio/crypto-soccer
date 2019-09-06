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
    uint8 constant public TEAMS_PER_DIVISION = 128; // LEAGUES_PER_DIV * TEAMS_PER_LEAGUE
    address constant public FREEVERSE = address(1);
    uint256 constant public DAYS_PER_ROUND = 16;
    
    mapping(uint256 => uint256) private _playerIdToState;

    PlayerState internal _playerStateLib;
    TimeZone[25] internal _timeZones;  // the first timeZone is a dummy one, without any country. Forbidden to use timeZone[0].
    mapping (uint256 => uint256) internal _playerIdxToPlayerState;
    uint256 public gameDeployMonth;

    constructor(address playerState) public {
        _playerStateLib = PlayerState(playerState);
        gameDeployMonth = secsToMonths(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
    }

    function _initTimeZone(uint8 tz) private {
        Country memory country;
        country.nDivisions = 1;
        _timeZones[tz].countries.push(country);
        _timeZones[tz].countries[0].divisonIdxToRound[0] = 1; 
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

    function isVirtualPlayer(uint256 playerId) public view returns (bool) {
        require(playerExists(playerId), "unexistent player");
        return _playerIdToState[playerId] == 0;
    }

    function transferBotInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) public {
        require(isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a non-bot team");
        require(addr != address(0));
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
        for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
            playerIds[p] = FREE_PLAYER_ID;
        }
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry] = Team(playerIds, addr);
    }

    function transferBotToAddr(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        transferBotInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
    }

    function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) public view returns(uint256) {
        if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
            return FREE_PLAYER_ID;
        } else {
            return _playerStateLib.encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry * PLAYERS_PER_TEAM_INIT + shirtNum);
        }
    }
  

    // TODO: we really don't need this function. Only for external use. Consider removal
    function getPlayerIdsInTeam(uint256 teamId) public view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        if (isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry)) {
            for (uint8 shirtNum = 0 ; shirtNum < PLAYERS_PER_TEAM_MAX ; shirtNum++){
                playerIds[shirtNum] = getDefaultPlayerIdForTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum);
            }
        } else {
            for (uint8 shirtNum = 0 ; shirtNum < PLAYERS_PER_TEAM_MAX ; shirtNum++){
                uint256 writtenId = _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
                if (writtenId == 0) {
                    playerIds[shirtNum] = getDefaultPlayerIdForTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum);
                } else {
                    playerIds[shirtNum] = writtenId;
                }
            }
        }
    }

    function getPlayerSkillsAtBirth(uint256 playerId) public view returns (uint256) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        uint256 division = teamIdxInCountry / TEAMS_PER_DIVISION;
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        uint256 dna = uint256(keccak256(abi.encode(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)));
        uint256 playerCreationMonth = (gameDeployMonth * 30 + _timeZones[timeZone].countries[countryIdxInTZ].divisonIdxToRound[division] * DAYS_PER_ROUND) / 30;
        uint256 monthOfBirth = computeBirthMonth(dna, playerCreationMonth);
        uint16[N_SKILLS] memory skills = computeSkills(dna);
        return _playerStateLib.encodePlayerSkills(skills, monthOfBirth, playerId);
    }

    function getPlayerStateAtBirth(uint256 playerId) public view returns (uint256) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        uint256 currentTeamId = _playerStateLib.encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry);
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        return _playerStateLib.encodePlayerState(playerId, currentTeamId, shirtNum, 0, 0);
    }

    function getPlayerState(uint256 playerId) public view returns (uint256) {
        if (isVirtualPlayer(playerId)) { 
            return getPlayerStateAtBirth(playerId);
        } else {
            return _playerIdToState[playerId];
        }
    }


    /// Compute a random age between 16 and 35
    /// @param dna is a random number used as seed of the skills
    /// @param playerCreationMonth since unix epoch
    /// @return monthOfBirth since unix epoch
    function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) public pure returns (uint16) {
        require(playerCreationMonth > 40*12, "invalid playerCreationMonth");
        dna >>= BITS_PER_SKILL*N_SKILLS;
        uint16 seed = uint16(dna & SKILL_MASK);
        uint16 age = 16 + (seed % 20);
        return uint16(playerCreationMonth - age * 12);
    }

    /// Compute the pseudorandom skills, sum of the skills is 250
    /// @param dna is a random number used as seed of the skills
    /// @return 5 skills
    function computeSkills(uint256 dna) public pure returns (uint16[N_SKILLS] memory) {
        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(dna & SKILL_MASK);
            dna >>= BITS_PER_SKILL;
        }
        /// Adjust skills to so that they add up to, maximum, 5*50 = 250.
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
    
    function isFreeShirt(uint256 teamId, uint8 shirtNum) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry),"cannot query about the shirt of a Bot Team");
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum] == FREE_PLAYER_ID;
    }
    
    function secsToMonths(uint256 secs) private pure returns (uint256) {
        return (secs * 12)/ 31536000;  // 31536000 = 3600 * 24 * 365
    }

    function monthsToSecs(uint256 months) private pure returns (uint256) {
        return (months * 31536000) / 12; // 31536000 = 3600 * 24 * 365
    }

    function getPlayerAgeInMonths(uint256 playerId) public view returns (uint256) {
        return secsToMonths(now - monthsToSecs(_playerStateLib.getMonthOfBirth(getPlayerSkillsAtBirth(playerId))));
    }

    function getFreeShirt(uint256 teamId) public view returns(uint8) {
        for (uint8 shirtNum = PLAYERS_PER_TEAM_MAX-1; shirtNum >= 0; shirtNum--) {
            if (isFreeShirt(teamId, shirtNum)) {
                return shirtNum;
            }
        }
        return PLAYERS_PER_TEAM_MAX;
    }

    function transferPlayer(uint256 playerId, uint256 teamIdTarget) public  {
        // warning: check of ownership of players and teams should be done before calling this function
        require(playerExists(playerId) && teamExists(teamIdTarget), "unexistent player or team");
        uint256 state = getPlayerState(playerId);
        uint256 newState = state;
        uint256 teamIdOrigin = _playerStateLib.getCurrentTeamId(state);
        require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
        uint256 shirtOrigin = _playerStateLib.getCurrentShirtNum(state);
        uint8 shirtTarget = getFreeShirt(teamIdTarget);
        require(shirtTarget != PLAYERS_PER_TEAM_MAX, "target team for transfer is already full");
        
        newState = _playerStateLib.setCurrentTeamId(newState, teamIdTarget);
        newState = _playerStateLib.setCurrentShirtNum(newState, shirtTarget);
        newState = _playerStateLib.setLastSaleBlock(newState, block.number);
        _playerIdToState[playerId] = newState;

        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamIdOrigin);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtOrigin] = FREE_PLAYER_ID;
        (timeZone, countryIdxInTZ, teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamIdTarget);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtTarget] = playerId;
    }
    
    
    // function transferTeam(uint256 teamId, address newOwnerAddr) public {
    //     (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
    //     transferBotInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
    // }


    // function getDivisionCreationRound(uint8 timeZone, uint256 countryIdx, uint256 divisionIdx) public view returns(uint256) {
    //     uint256 creationRound = _timeZones[timeZone].countries[countryIdx].divisonIdxToRound[divisionIdx];
    //     return (creationRound - 1)* DAYS_PER_ROUND;
    // }

    // function getTeamCreationTimestamp(uint256 teamId) public view returns (uint256) {
    //     (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = _playerStateLib.decodeTZCountryAndVal(teamId);
    //     require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
    //     return getDivisionCreationDay(timeZone, countryIdx, divisionIdx);
    // }

    // function getCurrentLeagueId(uint256 teamId) external view returns (uint256) {
    //     require(_teamExists(teamId), "invalid team id");
    //     return teams[teamId].currentLeagueId;
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
    //         uint16 birth = _computeBirthMonth(seed, getTeamCreationTimestamp(teamId));
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


}
