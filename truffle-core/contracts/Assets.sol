pragma solidity >=0.4.21 <0.6.0;

import "./Encoding.sol";
 
/// teamId == 0 is invalid and represents the null team
/// TODO: fix the playerPos <=> playerShirt doubt
contract Assets is Encoding {
    event TeamTransfer(uint256 teamId, address to);
    event ActionsSubmission(uint8 timeZone, bytes32 seed, uint256 submissionTime);

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
        bytes32[2] orgMapHash;
        bytes32[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        uint256 scoresRoot;
        uint8 updateCycleIdx;
        uint256 lastActionsSubmissionTime;
        uint256 lastUpdateTime;
        bytes32 actionsRoot;
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
    // uint256 constant public SEPT2019 = 1567296000; // UTC 1st of September, 2019, midnight, expressed in Unix Time
    uint16 constant public SECS_BETWEEN_VERSES = 900; // 15 mins
    address constant public NULL_ADDR = address(0);
    bytes32 constant INIT_ORGMAP_HASH = bytes32(0); // to compute externally and place here
    uint8 constant VERSES_PER_DAY = 96; // 24 * 4
    uint16 constant VERSES_PER_ROUND = 1536; // 96 * 16
    uint8 constant NULL_TIMEZONE = 0;
    uint8 constant CHALLENGE_TIME = 60; // in secs
    bool dummy = false;
    
    mapping(uint256 => uint256) private _playerIdToState;

    TimeZone[25] internal _timeZones;  // the first timeZone is a dummy one, without any country. Forbidden to use timeZone[0].
    mapping (uint256 => uint256) internal _playerIdxToPlayerState;
    uint256 public gameDeployMonth;
    uint256 public currentVerse;
    uint256 public currentRound;
    bytes32 private _currentVerseSeed;
    bool private _needsInit = true;
    uint256 public nextVerseTimestamp;
    uint8 public timeZoneForRound1;
    
    function init() public {
        require(_needsInit == true, "cannot initialize twice");
        // the game starts at verse = 0. The transition to verse = 1 will be at the next exact hour.
        // that will be the begining of Round = 1. So Round 1 starts at some timezone that depends on
        // the call to the contract init() function.
        uint256 secsOfDay   = now % (3600 * 24);
        uint256 hour        = secsOfDay / 3600;  // 0, ..., 23
        uint256 minute      = (secsOfDay - hour * 3600) / 60; // 0, ..., 59
        uint256 secs        = (secsOfDay - hour * 3600 - minute * 60); // 0, ..., 59
        if (minute < 42) {
            timeZoneForRound1 = 1 + uint8(hour);
            nextVerseTimestamp = now + (44-minute)*60 + (60 - secs);
        } else {
            timeZoneForRound1 = 2 + uint8(hour);
            nextVerseTimestamp = now + (44-minute)*60 + (60 - secs) + 3600;
        }
        setCurrentVerseSeed(blockhash(block.number-1)); 
        gameDeployMonth = secsToMonths(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        _needsInit = false;
    }
    
    function getNow() public view returns(uint256) {
        return now;
    }

    function _initTimeZone(uint8 tz) private {
        Country memory country;
        country.nDivisions = 1;
        _timeZones[tz].countries.push(country);
        _timeZones[tz].countries[0].divisonIdxToRound[0] = 1; 
        _timeZones[tz].orgMapHash[0] = INIT_ORGMAP_HASH; 
    }
    
    
    function submitActionsRoot(bytes32 actionsRoot) public {
        require(now > nextVerseTimestamp, "too early to accept actions root");
        (uint8 tz,,) = timeZoneToUpdate();
        require(now - _timeZones[tz].lastUpdateTime > CHALLENGE_TIME, "last verse is still under challenge period");
        _timeZones[tz].actionsRoot = actionsRoot;
        _timeZones[tz].lastActionsSubmissionTime = block.number;
        setCurrentVerseSeed(blockhash(block.number-1)); 
        emit ActionsSubmission(tz, blockhash(block.number-1), now);
    }
    
    // each day has 24 hours, each with 4 verses => 96 verses per day.
    // day = 1,..16
    // turnInDay = 0, 1, 2, 3
    // so for each TZ, we go from (day, turn) = (1, 0) ... (15,3) => a total of 16*4 = 64 turns per timeZone
    // from these, all map easily to timeZones
    function timeZoneToUpdate() public view returns (uint8 timeZone, uint8 day, uint8 turnInDay) {
        return _timeZoneToUpdatePure(currentVerse, timeZoneForRound1);
    }

    function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) public pure returns (uint8 timeZone, uint8 day, uint8 turnInDay) {
        // if currentVerse = 0, we should be updating timeZoneForRound1
        // recall that timeZones range from 1...24 (not from 0...24)
        uint16 verseInRound = uint16(verse % VERSES_PER_ROUND);
        if (verseInRound < VERSES_PER_DAY) {
            timeZone = 1 + uint8((TZForRound1 - 1 + (verseInRound / 4))% 24);
            day = 1;
            turnInDay = uint8(verseInRound % 4);
        } else if (verseInRound == VERSES_PER_DAY) {
            timeZone = NULL_TIMEZONE;
        } else {
            timeZone = 1 + uint8((TZForRound1 - 1 + ((verseInRound - 1) / 4))% 24);
            day = 1 + uint8((verseInRound - 1) / VERSES_PER_DAY);
            turnInDay = uint8((verseInRound - 1) % 4);
        }
    }

    function wait() public returns (uint256) {
        dummy = !dummy;
        return now; 
    }
    
    function setCurrentVerseSeed(bytes32 seed) public {
        _currentVerseSeed = seed;
    }
        
    function getCurrentVerseSeed() public view returns (bytes32) {
        return _currentVerseSeed;
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

    function transferBotInCountryToAddr(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) public {
        require(isBotTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "cannot transfer a non-bot team");
        require(addr != NULL_ADDR);
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
        for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
            playerIds[p] = FREE_PLAYER_ID;
        }
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry] = Team(playerIds, addr);
    }

    function transferBotToAddr(uint256 teamId, address addr) public {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        transferBotInCountryToAddr(timeZone, countryIdxInTZ, teamIdxInCountry, addr);
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

    function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) public pure returns(uint256) {
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
        uint256 monthOfBirth = computeBirthMonth(dna, playerCreationMonth);
        uint16[N_SKILLS] memory skills = computeSkills(dna);
        return encodePlayerSkills(skills, monthOfBirth, playerId);
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
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
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
