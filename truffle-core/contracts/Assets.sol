pragma solidity >=0.4.21 <0.6.0;

import "./Encoding.sol";
/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Assets is Encoding {
    event TeamTransfer(uint256 teamId, address to);
    event PlayerTransfer(uint256 playerId, uint256 teamIdTarget);
    event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry);

    struct Team {
        uint256[PLAYERS_PER_TEAM_MAX] playerIds; 
        address owner;
    }

    struct Country {
        uint256 nDivisions;
        uint8 nDivisionsToAddNextRound;
        mapping (uint256 => uint256) divisonIdxToRound;
        mapping (uint256 => Team) teamIdxInCountryToTeam;
        uint256 nHumanTeams;
    }

    struct TimeZone {
        Country[] countries;
        uint8 nCountriesToAdd;
        bytes32[2] orgMapHash;
        bytes32[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        bytes32 scoresRoot;
        uint8 updateCycleIdx;
        uint256 lastActionsSubmissionTime;
        uint256 lastUpdateTime;
        bytes32 actionsRoot;
        uint256 lastMarketClosureBlockNum;
    }    
    
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint8 constant public N_SKILLS = 5;
    uint8 constant public LEAGUES_PER_DIV = 16;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public TEAMS_PER_DIVISION = 128; // LEAGUES_PER_DIV * TEAMS_PER_LEAGUE
    address constant public FREEVERSE = address(1);
    uint256 constant public DAYS_PER_ROUND = 16;
    address constant public NULL_ADDR = address(0);
    bytes32 constant INIT_ORGMAP_HASH = bytes32(0); // to be computed externally once and placed here
    
   
    // skills idxs: Defence, Speed, Pass, Shoot, Endurance
    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;
    
    mapping(uint256 => uint256) private _playerIdToState;

    TimeZone[25] public _timeZones;  // timeZone = 0 is a dummy one, without any country. Forbidden to use timeZone[0].
    mapping (uint256 => uint256) internal _playerIdxToPlayerState;
    uint256 public gameDeployMonth;
    uint256 public currentRound;
    bool private _needsInit = true;

    function init() public {
        require(_needsInit == true, "cannot initialize twice");
        gameDeployMonth = secsToMonths(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        _needsInit = false;
    }

    function _initTimeZone(uint8 tz) private {
        Country memory country;
        country.nDivisions = 1;
        _timeZones[tz].countries.push(country);
        _timeZones[tz].countries[0].divisonIdxToRound[0] = 1;
        _timeZones[tz].orgMapHash[0] = INIT_ORGMAP_HASH;
        emit DivisionCreation(tz, 0, 0);
    }

    function getLastUpdateTime(uint8 timeZone) external view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastUpdateTime;
    }

    function getLastActionsSubmissionTime(uint8 timeZone) external view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastActionsSubmissionTime;
    }

    function setSkillsRoot(uint8 tz, bytes32 root) external returns(uint256) {
        _timeZones[tz].skillsHash[_timeZones[tz].newestSkillsIdx] = root;
        _timeZones[tz].newestSkillsIdx = 1 - _timeZones[tz].newestSkillsIdx;
        _timeZones[tz].lastUpdateTime = now;
    }

    function setActionsRoot(uint8 timeZone, bytes32 root) external returns(uint256) {
        _assertTZExists(timeZone);
        _timeZones[timeZone].actionsRoot = root;
        _timeZones[timeZone].lastActionsSubmissionTime = now;
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
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return _teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry) == NULL_ADDR;
    }

    function isBotTeam(uint256 teamId) public view returns(bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    // returns NULL_ADDR if team is bot
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(address) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner;
    }

    function getOwnerTeam(uint256 teamId) public view returns(address) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }

    // returns NULL_ADDR if team is bot
    function getOwnerPlayer(uint256 playerId) public view returns(address) {
        require(playerExists(playerId), "unexistent player");
        uint256 teamId = getCurrentTeamId(getPlayerState(playerId));
        return getOwnerTeam(teamId);
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
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        return _wasPlayerCreatedInCountry(timeZone, countryIdxInTZ, playerIdxInCountry);
    }

    function isVirtualPlayer(uint256 playerId) public view returns (bool) {
        require(playerExists(playerId), "unexistent player");
        return _playerIdToState[playerId] == 0;
    }

    function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) external {
        uint256 firstBotIdx = _timeZones[timeZone].countries[countryIdxInTZ].nHumanTeams;
        require(isBotTeamInCountry(timeZone, countryIdxInTZ, firstBotIdx), "cannot transfer a non-bot team");
        require(addr != NULL_ADDR, "invalid address");
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
        for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
            playerIds[p] = FREE_PLAYER_ID;
        }
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[firstBotIdx] = Team(playerIds, addr);
        _timeZones[timeZone].countries[countryIdxInTZ].nHumanTeams++;
        uint256 teamId = encodeTZCountryAndVal(timeZone, countryIdxInTZ, firstBotIdx);
        emit TeamTransfer(teamId, addr);
    }

    function transferTeamInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) private {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a non-bot team");
        require(addr != NULL_ADDR, "cannot transfer to a null address");
        require(_timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner != addr, "buyer and seller are the same addr");
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = addr;
    }

    function transferTeam(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        transferTeamInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
        emit TeamTransfer(teamId, addr);
    }

    function getDefaultPlayerIdForTeamInCountry(
        uint8 timeZone,
        uint256 countryIdxInTZ,
        uint256 teamIdxInCountry,
        uint8 shirtNum
    )
        public
        pure
        returns(uint256)
    {
        if (shirtNum >= PLAYERS_PER_TEAM_INIT) {
            return FREE_PLAYER_ID;
        } else {
            return encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry * PLAYERS_PER_TEAM_INIT + shirtNum);
        }
    }

    // TODO: we really don't need this function. Only for external use. Consider removal
    function getPlayerIdsInTeam(uint256 teamId) public view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
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
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        uint256 division = teamIdxInCountry / TEAMS_PER_DIVISION;
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        uint256 dna = uint256(keccak256(abi.encode(timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)));
        uint256 playerCreationMonth = (gameDeployMonth * 30 + _timeZones[timeZone].countries[countryIdxInTZ].divisonIdxToRound[division] * DAYS_PER_ROUND) / 30;
        return computeSkillsAndEncode(dna, shirtNum, playerCreationMonth, playerId);
    }

    // the next function was separated from getPlayerSkillsAtBirth only to keep stack within limits
    function computeSkillsAndEncode(uint256 dna, uint8 shirtNum, uint256 playerCreationMonth, uint256 playerId) internal pure returns (uint256) {
        uint256 monthOfBirth;
        (monthOfBirth, dna) = computeBirthMonth(dna, playerCreationMonth);
        (uint16[N_SKILLS] memory skills, uint8 potential, uint8 forwardness, uint8 leftishness) = computeSkills(dna, shirtNum);
        return encodePlayerSkills(skills, monthOfBirth, playerId, potential, forwardness, leftishness);
    }

    function getPlayerStateAtBirth(uint256 playerId) public view returns (uint256) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        uint256 currentTeamId = encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry);
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        return encodePlayerState(playerId, currentTeamId, shirtNum, 0, 0);
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
    function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) public pure returns (uint16, uint256) {
        require(playerCreationMonth > 40*12, "invalid playerCreationMonth");
        uint16 age = 16 + uint16((dna % 20));
        dna >>= 5; // log2(20) = 4.6.... ceil = 5
        return ( uint16(playerCreationMonth - age * 12), dna);
    }

    /// Compute the pseudorandom skills, sum of the skills is 250
    /// @param dna is a random number used as seed of the skills
    /// skills have currently, 16bits each, and there are 5 of them
    /// potential is a number between 0 and 9 => takes 4 bit
    /// 0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    /// @return uint16[N_SKILLS] skills, uint8 potential, uint8 forwardness, uint8 leftishness
    function computeSkills(uint256 dna, uint8 shirtNum) public pure returns (uint16[N_SKILLS] memory, uint8, uint8, uint8) {
        uint16[5] memory skills;
        uint16[N_SKILLS] memory correctFactor;
        uint8 potential = uint8(dna % 10);
        uint8 forwardness;
        uint8 leftishness;
        dna >>= 4; // log2(10) = 3.3 => ceil = 4
        if (shirtNum < 3) {
            // 3 GoalKeepers:
            correctFactor[SK_SHO] = 200;
            forwardness = IDX_GK;
            leftishness = 0;
        } else if (shirtNum < 8) {
            // 5 Defenders
            correctFactor[SK_SHO] = 40;
            correctFactor[SK_DEF] = 160;
            forwardness = IDX_D;
            leftishness = uint8(1+ ((dna + shirtNum) % 7));
        } else if (shirtNum < 10) {
            // 2 Pure Midfielders
            correctFactor[SK_PAS] = 160;
            forwardness = IDX_M;
            leftishness = uint8(1+ ((dna + shirtNum) % 7));
        } else if (shirtNum < 12) {
            // 2 Defensive Midfielders
            correctFactor[SK_PAS] = 130;
            correctFactor[SK_SHO] = 70;
            forwardness = IDX_MD;
            leftishness = uint8(1+ ((dna + shirtNum) % 7));
        } else if (shirtNum < 14) {
            // 2 Attachking Midfielders
            correctFactor[SK_PAS] = 130;
            correctFactor[SK_DEF] = 70;
            forwardness = IDX_MF;
            leftishness = uint8(1+ ((dna + shirtNum) % 7));
        } else if (shirtNum < 16) {
            // 2 Forwards that play center-left
            correctFactor[SK_SHO] = 160;
            correctFactor[SK_DEF] = 70;
            forwardness = IDX_F;
            leftishness = 6;
        } else {
            // 2 Forwards that play center-right
            correctFactor[SK_SHO] = 160;
            correctFactor[SK_DEF] = 70;
            forwardness = IDX_F;
            leftishness = 3;
        }
        dna >>= 51; // log2(7) = 2.9 => ceil = 3, times 17 players => 51                       

        /// Compute initial skills, as a random with [0, 49] 
        /// ...apply correction factor depending on preferred pos,
        //  ...and adjust skills to so that they add up to, at least, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < N_SKILLS; i++) {
            if (correctFactor[i] == 0) {
                skills[i] = uint16(dna % 50);
            } else {
                skills[i] = (uint16(dna % 50) * correctFactor[i])/100;
            }
            dna >>= 6; // los2(50) -> ceil
            excess += skills[i];
        }
        if (excess < 250) {
            uint16 delta = (250 - excess) / N_SKILLS;
            for (uint8 i = 0; i < 5; i++) skills[i] = skills[i] + delta;
        }
        return (skills, potential, forwardness, leftishness);
    }

    function isFreeShirt(uint256 teamId, uint8 shirtNum) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        require(!isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry),"cannot query about the shirt of a Bot Team");
        uint256 writtenId = _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
        if (shirtNum > PLAYERS_PER_TEAM_INIT - 1) {
            return (writtenId == 0 || writtenId == FREE_PLAYER_ID);
        } else {
            return writtenId == FREE_PLAYER_ID;
        }
    }

    function secsToMonths(uint256 secs) private pure returns (uint256) {
        return (secs * 12)/ 31536000;  // 31536000 = 3600 * 24 * 365
    }

    function monthsToSecs(uint256 months) private pure returns (uint256) {
        return (months * 31536000) / 12; // 31536000 = 3600 * 24 * 365
    }

    function getPlayerAgeInMonths(uint256 playerId) public view returns (uint256) {
        return secsToMonths(now - monthsToSecs(getMonthOfBirth(getPlayerSkillsAtBirth(playerId))));
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
        // TODO: checking if they are bots should be done outside this function
        require(playerExists(playerId) && teamExists(teamIdTarget), "unexistent player or team");
        uint256 state = getPlayerState(playerId);
        uint256 newState = state;
        uint256 teamIdOrigin = getCurrentTeamId(state);
        require(teamIdOrigin != teamIdTarget, "cannot transfer to original team");
        require(!isBotTeam(teamIdOrigin) && !isBotTeam(teamIdTarget), "cannot transfer player when at least one team is a bot");
        uint256 shirtOrigin = getCurrentShirtNum(state);
        uint8 shirtTarget = getFreeShirt(teamIdTarget);
        require(shirtTarget != PLAYERS_PER_TEAM_MAX, "target team for transfer is already full");
        
        newState = setCurrentTeamId(newState, teamIdTarget);
        newState = setCurrentShirtNum(newState, shirtTarget);
        newState = setLastSaleBlock(newState, block.number);
        _playerIdToState[playerId] = newState;

        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamIdOrigin);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtOrigin] = FREE_PLAYER_ID;
        (timeZone, countryIdxInTZ, teamIdxInCountry) = decodeTZCountryAndVal(teamIdTarget);
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtTarget] = playerId;
        emit PlayerTransfer(playerId, teamIdTarget);
    }

    function countCountries(uint8 timeZone) public view returns (uint256){
        _assertTZExists(timeZone);
        return _timeZones[timeZone].countries.length;
    }

    function countTeams(uint8 timeZone, uint256 countryIdxInTZ) public view returns (uint256){
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return _timeZones[timeZone].countries[countryIdxInTZ].nDivisions * TEAMS_PER_DIVISION;
    }
}
