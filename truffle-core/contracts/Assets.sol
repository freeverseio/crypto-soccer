pragma solidity >=0.5.12 <0.6.2;

import "./EncodingSkills.sol";
import "./EncodingIDs.sol";
import "./EncodingState.sol";
import "./Storage.sol";
import "./StorageProxy.sol";
import "./AssetsLib.sol";

/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Assets is AssetsLib, EncodingSkills, EncodingState {

    event AssetsInit(address creatorAddr);
    event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry);
    
    uint8 constant public TEAMS_PER_DIVISION = 128; // LEAGUES_PER_DIV * TEAMS_PER_LEAGUE
    uint256 constant public DAYS_PER_ROUND = 16;
    bytes32 constant INIT_ORGMAP_HASH = bytes32(0); // to be computed externally once and placed here

    function init() public {
        require(_wasInited == false, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        _wasInited = true;
        emit AssetsInit(msg.sender);
    }

    // hack for testing: we can init only one timezone
    // at some point, remove this option
    function initSingleTZ(uint8 tz) public {
        require(_wasInited == false, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        _initTimeZone(tz);
        _wasInited = true;
        emit AssetsInit(msg.sender);
    }
    

    function _initTimeZone(uint8 tz) private {
        Country memory country;
        country.nDivisions = 1;
        _timeZones[tz].countries.push(country);
        _timeZones[tz].orgMapHash[0] = INIT_ORGMAP_HASH;
        for (uint8 division = 0 ; division < country.nDivisions ; division++){
            _timeZones[tz].countries[0].divisonIdxToRound[division] = 1;
            emit DivisionCreation(tz, 0, division);
        }
    }

    function getLastUpdateTime(uint8 timeZone) external view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastUpdateTime;
    }

    function getLastActionsSubmissionTime(uint8 timeZone) external view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastActionsSubmissionTime;
    }

    function setSkillsRoot(uint8 tz, bytes32 root, bool newTZ) external returns(uint256) {
        if (newTZ) _timeZones[tz].newestSkillsIdx = 1 - _timeZones[tz].newestSkillsIdx;
        _timeZones[tz].skillsHash[_timeZones[tz].newestSkillsIdx] = root;
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

    function getPlayerSkillsAtBirth(uint256 playerId) public view returns (uint256) {
        if (getIsSpecial(playerId)) return getSpecialPlayerSkillsAtBirth(playerId);
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        uint8 shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
        uint256 division = teamIdxInCountry / TEAMS_PER_DIVISION;
        require(_teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry), "invalid team id");
        // compute a dna that is unique to this player, since it is made of a unique playerId:
        uint256 playerCreationDay = gameDeployDay + _timeZones[timeZone].countries[countryIdxInTZ].divisonIdxToRound[division] * DAYS_PER_ROUND;
        return computeSkillsAndEncode(shirtNum, playerCreationDay, playerId);
    }

    function getSpecialPlayerSkillsAtBirth(uint256 playerId) internal pure returns (uint256) {
        return playerId;
    }

    // the next function was separated from getPlayerSkillsAtBirth only to keep stack within limits
    function computeSkillsAndEncode(uint8 shirtNum, uint256 playerCreationDay, uint256 playerId) internal pure returns (uint256) {
        uint256 dna = uint256(keccak256(abi.encode(playerId)));
        uint256 dayOfBirth;
        (dayOfBirth, dna) = computeBirthDay(dna, playerCreationDay);
        (uint16[N_SKILLS] memory skills, uint8[4] memory birthTraits, uint32 sumSkills) = computeSkills(dna, shirtNum);
        return encodePlayerSkills(skills, dayOfBirth, 0, playerId, birthTraits, false, false, 0, 0, false, sumSkills);
    }


    /// Compute a random age between 16 and 35
    /// @param dna is a random number used as seed of the skills
    /// @param playerCreationDay since unix epoch
    /// @return dayOfBirth since unix epoch
    function computeBirthDay(uint256 dna, uint256 playerCreationDay) public pure returns (uint16, uint256) {
        uint256 ageInDays = 5840 + (dna % 7300);  // 5840 = 16*365, 7300 = 20 * 365
        dna >>= 13; // log2(7300) = 12.8
        return (uint16(playerCreationDay - ageInDays / 7), dna); // 1095 = 3 * 365
    }
    
    /// Compute the pseudorandom skills, sum of the skills is 5K (1K each skill on average)
    /// @param dna is a random number used as seed of the skills
    /// skills have currently, 16bits each, and there are 5 of them
    /// potential is a number between 0 and 9 => takes 4 bit
    /// 0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    /// @return uint16[N_SKILLS] skills, uint8 potential, uint8 forwardness, uint8 leftishness
    function computeSkills(uint256 dna, uint8 shirtNum) public pure returns (uint16[N_SKILLS] memory, uint8[4] memory, uint32) {
        uint16[5] memory skills;
        uint16[N_SKILLS] memory correctFactor;
        uint8 potential = uint8(dna % 10);
        dna >>= 4; // log2(10) = 3.3 => ceil = 4
        uint8 forwardness;
        uint8 leftishness;
        uint8 aggressiveness = uint8(dna % 4);
        dna >>= 2; // log2(4) = 2
        // correctFactor/1000 increases a particular skill depending on player's forwardness
        if (shirtNum < 3) {
            // 3 GoalKeepers:
            correctFactor[SK_SHO] = 2000;
            forwardness = IDX_GK;
            leftishness = 0;
        } else if (shirtNum < 8) {
            // 5 Defenders
            correctFactor[SK_SHO] = 400;
            correctFactor[SK_DEF] = 1600;
            forwardness = IDX_D;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 10) {
            // 2 Pure Midfielders
            correctFactor[SK_PAS] = 1600;
            forwardness = IDX_M;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 12) {
            // 2 Defensive Midfielders
            correctFactor[SK_PAS] = 1300;
            correctFactor[SK_SHO] = 700;
            forwardness = IDX_MD;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 14) {
            // 2 Attachking Midfielders
            correctFactor[SK_PAS] = 1300;
            correctFactor[SK_DEF] = 700;
            forwardness = IDX_MF;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 16) {
            // 2 Forwards that play center-left
            correctFactor[SK_SHO] = 1600;
            correctFactor[SK_DEF] = 700;
            forwardness = IDX_F;
            leftishness = 6;
        } else {
            // 2 Forwards that play center-right
            correctFactor[SK_SHO] = 1600;
            correctFactor[SK_DEF] = 700;
            forwardness = IDX_F;
            leftishness = 3;
        }
        dna >>= 3; // log2(7) = 2.9 => ceil = 3                      

        /// Compute initial skills, as a random with [0, 49] 
        /// ...apply correction factor depending on preferred pos,
        //  ...and adjust skills to so that they add up to, at least, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < N_SKILLS; i++) {
            if (correctFactor[i] == 0) {
                skills[i] = uint16(dna % 1000);
            } else {
                skills[i] = (uint16(dna % 1000) * correctFactor[i])/1000;
            }
            excess += skills[i];
            dna >>= 10; // los2(1000) -> ceil
        }
        // at this point, excess is, at most, 5*999 = 4995, so (5000 - excess) > 0
        uint16 delta;
        delta = (5000 - excess) / N_SKILLS;
        for (uint8 i = 0; i < N_SKILLS; i++) skills[i] = skills[i] + delta;
        // note: final sum of skills = excess + N_SKILLS * delta;
        return (skills, [potential, forwardness, leftishness, aggressiveness], uint32(excess + N_SKILLS * delta));
    }


    function secsToDays(uint256 secs) private pure returns (uint256) {
        return secs / 86400;  // 86400 = 3600 * 24
    }

    function daysToSecs(uint256 dayz) private pure returns (uint256) {
        return dayz * 86400; // 86400 = 3600 * 24 * 365
    }

    function getPlayerAgeInDays(uint256 playerId) public view returns (uint256) {
        return secsToDays(7 * (now - daysToSecs(getBirthDay(getPlayerSkillsAtBirth(playerId)))));
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
