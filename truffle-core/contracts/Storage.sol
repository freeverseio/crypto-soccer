pragma solidity >=0.5.12 <0.6.2;

/**
 * @title Storage for all assets and ownerships
 */

contract Storage {

    // TODO: perhaps make these input params from functions, or put setters that allow changing them for future proofness.
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public PLAYERS_PER_TEAM_INIT = 18;
    uint256 constant public FREE_PLAYER_ID      = 1; // it never corresponds to a legit playerId due to its TZ = 0

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
        bytes32[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        uint256 lastActionsSubmissionTime;
        uint256 lastUpdateTime;
        bytes32 actionsRoot;
    }    
    
    TimeZone[25] public _timeZones;  // timeZone = 0 is a dummy one, without any country. Forbidden to use timeZone[0].
    uint256 public _gameDeployDay;
    mapping(uint256 => uint256) private _playerIdToState;

    function getGameDeployDay() external view returns (uint256) { return _gameDeployDay; }
    function setGameDeployDay(uint256 val) external { _gameDeployDay = val; }

    function pushCountryToTZ(uint8 tz, uint256 nDivisions) external { 
        Country memory country;
        country.nDivisions = nDivisions;
        _timeZones[tz].countries.push(country); 
        for (uint8 division = 0 ; division < country.nDivisions ; division++){
            _timeZones[tz].countries[0].divisonIdxToRound[division] = 1;
        }
    }


    function getDivisonIdxToRound(uint8 tz, uint256 countryIdxInTZ, uint256 division) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].divisonIdxToRound[division] ;
    }

    function getLastUpdateTime(uint8 tz) external view returns(uint256) {
        return _timeZones[tz].lastUpdateTime;
    }

    function getLastActionsSubmissionTime(uint8 tz) external view returns(uint256) {
        return _timeZones[tz].lastActionsSubmissionTime;
    }

    function setSkillsRoot(uint8 tz, bytes32 root, bool newTZ) external returns(uint256) {
        if (newTZ) _timeZones[tz].newestSkillsIdx = 1 - _timeZones[tz].newestSkillsIdx;
        _timeZones[tz].skillsHash[_timeZones[tz].newestSkillsIdx] = root;
        _timeZones[tz].lastUpdateTime = now;
    }

    function setActionsRoot(uint8 tz, bytes32 root, uint256 time) external returns(uint256) {
        _timeZones[tz].actionsRoot = root;
        _timeZones[tz].lastActionsSubmissionTime = time;
    }

    function getNCountriesInTZ(uint8 tz) public view returns(uint256) {
        return _timeZones[tz].countries.length;
    }

    function getNDivisionsInCountry(uint8 tz, uint256 countryIdxInTZ) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].nDivisions;
    }

    // returns NULL_ADDR if team is bot
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) external view returns(address) {
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner;
    }

    function setOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address newOwner) external {
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = newOwner;
    }

    function getNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].nHumanTeams;
    }

    function setNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ, uint256 val) external {
        _timeZones[tz].countries[countryIdxInTZ].nHumanTeams = val;
    }
    
    function incrementNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ) external {
        _timeZones[tz].countries[countryIdxInTZ].nHumanTeams++;
    }
    
    function assignBotToAddr(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) external {
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
        for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
            playerIds[p] = FREE_PLAYER_ID;
        }
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry] = Team(playerIds, addr);
    }    

    function changeTeamAddr(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) private {
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = addr;
    }

    function getPlayerIdsInTeam(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry) external view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
        return _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds;
    }

    function getPlayerIdFromShirt(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) external view returns (uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
    }

    function setPlayerIdFromShirt(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum, uint256 newPlayerId) external {
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum] = newPlayerId;
    }

    function getPlayerState(uint256 playerId) external view returns (uint256) { return _playerIdToState[playerId]; }

    function setPlayerState(uint256 playerId, uint256 val) external { _playerIdToState[playerId] = val; }    

    function countCountries(uint8 tz) external view returns (uint256){
        return _timeZones[tz].countries.length;
    }

}
