pragma solidity >=0.5.12 <0.6.2;

/**
 * @title Storage for all assets, their ownership, and the evolution of the game
 * security is split into three addresses:
 *      - storageOwner: owner of this contract, required to change any of the 3 addresses
 *      - assetsOwner:  required to use any setter related to assets data
 *      - updatesOwner: required to use any setter related to updates data
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

    // Storage data:
    // ...Internal Storage Security:
    address public _storageOwner;
    address public _assetsOwner;
    address public _marketOwner;
    address public _updatesOwner;
    
    // ...Storage for other Assets and Updates:
    TimeZone[25] public _timeZones;  // note: _timeZone[0] is a dummy one, without any country
    uint256 public _gameDeployDay;
    mapping(uint256 => uint256) private _playerIdToState;
    
    // ...Storage for other Market:
    mapping (uint256 => uint256) private _playerIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToAuctionData;
    mapping (uint256 => uint256) private _teamIdToAcqConstraints;

    // Contructor and Storage Security modifiers:
    constructor() public { _storageOwner = msg.sender; }

    modifier onlyOwner {
        require(msg.sender == _storageOwner, "only owner of Storage can set a new Storage owner");
        _;
    }
    modifier onlyAssets {
        require(msg.sender == _assetsOwner, "only owner of Storage can set a new Assets owner");
        _;
    }
    modifier onlyMarket {
        require(msg.sender == _marketOwner, "only owner of Storage can set a new Market owner");
        _;
    }
    modifier onlyUpdates {
        require(msg.sender == _updatesOwner, "only owner of Storage can set a new Updates owner");
        _;
    }

    // Internal Storage Security Functions:
    function setStorageOwner(address newOwner) external onlyOwner { _storageOwner = newOwner; }
    function setAssetsOwner(address newOwner) external onlyOwner { _assetsOwner = newOwner; }
    function setMarketOwner(address newOwner) external onlyOwner { _marketOwner = newOwner; }
    function setUpdatesOwner(address newOwner) external onlyOwner { _updatesOwner = newOwner; }

    // Assets Setters:
    function setPlayerState(uint256 playerId, uint256 val) external onlyAssets { _playerIdToState[playerId] = val; }    
    function setGameDeployDay(uint256 val) external onlyAssets { _gameDeployDay = val; }

    function setOwnerTeamInCountry(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address newOwner) external onlyAssets {
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner = newOwner;
    }
    function setNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ, uint256 val) external onlyAssets {
        _timeZones[tz].countries[countryIdxInTZ].nHumanTeams = val;
    }
    function incrementNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ) external onlyAssets {
        _timeZones[tz].countries[countryIdxInTZ].nHumanTeams++;
    }
    function setPlayerIdFromShirt(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum, uint256 newPlayerId) external onlyAssets {
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum] = newPlayerId;
    }
    function pushCountryToTZ(uint8 tz, uint256 nDivisions) external onlyAssets { 
        Country memory country;
        country.nDivisions = nDivisions;
        _timeZones[tz].countries.push(country); 
        for (uint8 division = 0 ; division < country.nDivisions ; division++){
            _timeZones[tz].countries[0].divisonIdxToRound[division] = 1;
        }
    }
    function assignBotToAddr(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, address addr) external {
        uint256[PLAYERS_PER_TEAM_MAX] memory playerIds;
        for (uint p = PLAYERS_PER_TEAM_INIT; p < PLAYERS_PER_TEAM_MAX; p++) {
            playerIds[p] = FREE_PLAYER_ID;
        }
        _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry] = Team(playerIds, addr);
    }    

    // Market Setters:
    function setAcquisitionConstraint(uint256 teamId, uint256 constraint) external onlyMarket returns(uint256) {
        _teamIdToAcqConstraints[teamId] = constraint;
    }
    function setPlayerIdToAuctionData(uint256 playerId, uint256 auctionData) external onlyMarket returns(uint256) {
        _playerIdToAuctionData[playerId] = auctionData;
    }
    function setTeamIdToAuctionData(uint256 teamId, uint256 auctionData) external onlyMarket returns(uint256) {
        _teamIdToAuctionData[teamId] = auctionData;
    }

    // Updates Setters:
    function setSkillsRoot(uint8 tz, bytes32 root, bool newTZ) external onlyUpdates returns(uint256) {
        if (newTZ) _timeZones[tz].newestSkillsIdx = 1 - _timeZones[tz].newestSkillsIdx;
        _timeZones[tz].skillsHash[_timeZones[tz].newestSkillsIdx] = root;
        _timeZones[tz].lastUpdateTime = now;
    }
    function setActionsRoot(uint8 tz, bytes32 root, uint256 time) external onlyUpdates returns(uint256) {
        _timeZones[tz].actionsRoot = root;
        _timeZones[tz].lastActionsSubmissionTime = time;
    }
    
    // Assets Getters:
    function getPlayerState(uint256 playerId) external view returns (uint256) { return _playerIdToState[playerId]; }
    function getGameDeployDay() external view returns (uint256) { return _gameDeployDay; }
    function getPlayerIdsInTeam(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry) external view returns (uint256[PLAYERS_PER_TEAM_MAX] memory playerIds) {
        return _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds;
    }
    function getPlayerIdFromShirt(uint8 tz, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) external view returns (uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].playerIds[shirtNum];
    }
    function countCountries(uint8 tz) external view returns (uint256){
        return _timeZones[tz].countries.length;
    }
    function getDivisonIdxToRound(uint8 tz, uint256 countryIdxInTZ, uint256 division) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].divisonIdxToRound[division] ;
    }
    function getNCountriesInTZ(uint8 tz) external view returns(uint256) {
        return _timeZones[tz].countries.length;
    }
    function getNDivisionsInCountry(uint8 tz, uint256 countryIdxInTZ) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].nDivisions;
    }
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) external view returns(address) {
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner;
    }
    function getNHumanTeamsInCountry(uint8 tz, uint256 countryIdxInTZ) external view returns(uint256) {
        return _timeZones[tz].countries[countryIdxInTZ].nHumanTeams;
    }
    
    // Market Getters:
    function getAcquisitionConstraint(uint256 teamId) external view returns(uint256) {
        return _teamIdToAcqConstraints[teamId];
    }
    function getAuctionDataForTeam(uint256 teamId) external view returns(uint256) {
        return _teamIdToAuctionData[teamId];
    }
    function getAuctionDataForPlayer(uint256 playerId) external view returns(uint256) {
        return _playerIdToAuctionData[playerId];
    }
    
    // Updates Getters:
    function getLastUpdateTime(uint8 tz) external view returns(uint256) {
        return _timeZones[tz].lastUpdateTime;
    }
    function getLastActionsSubmissionTime(uint8 tz) external view returns(uint256) {
        return _timeZones[tz].lastActionsSubmissionTime;
    }
    


}
