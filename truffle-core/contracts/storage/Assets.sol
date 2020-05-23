pragma solidity >= 0.6.3;

import "../storage/AssetsView.sol";

/**
 @title Creation of all "default" game assets via creation of timezones, countries and divisions
 @author Freeverse.io, www.freeverse.io
 @dev Only other way of creating assets is via BuyNow pattern, in the Market contract.
 @dev All functions in this file modify storage. All view/pure funcions are inherited from AssetsView.
 @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 @dev All storage is govenrned by Proxy, via the Storage contract.
*/

contract Assets is AssetsView {

    event AssetsInit(address creatorAddr);
    event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry);
    
    //// Setter for main roles: COO, Market owner, Relay owner
    function setCOO(address addr) external onlySuperUser { _COO = addr; }
    
    function setMarket(address addr) external onlySuperUser {
        _market = addr;
        teamIdToOwner[ACADEMY_TEAM] = addr;
        if (gameDeployDay == 0) { emit AssetsInit(msg.sender); }
        emit TeamTransfer(ACADEMY_TEAM, addr);        
    }
    
    function setRelay(address addr) external onlySuperUser { _relay = addr; }
   

    /// External Functions

    /// Inits all 24 timezones, each with one country, each with one division
    function initTZs() external onlyCOO {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        if (_market == NULL_ADDR) { emit AssetsInit(msg.sender); }
    }

    /// Next function is only used for testing: it inits only one timezone
    function initSingleTZ(uint8 tz) external onlyCOO {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        _initTimeZone(tz);
        if (_market == NULL_ADDR) { emit AssetsInit(msg.sender); }
    }

    function addDivisionManually(uint8 tz, uint256 countryIdxInTZ) external onlyCOO { _addDivision(tz, countryIdxInTZ); }

    function addCountryManually(uint8 tz) external onlyCOO { _addCountry(tz); }

    /// Entry point for new users: acquiring a bot team
    function transferFirstBotToAddr(uint8 tz, uint256 countryIdxInTZ, address addr) public onlyRelay {
        require(tzToNCountries[tz] != 0, "Timezone has not been initialized");
        uint256 countryId = encodeTZCountryAndVal(tz, countryIdxInTZ, 0); 
        uint256 firstBotIdx = countryIdToNHumanTeams[countryId];
        uint256 teamId = encodeTZCountryAndVal(tz, countryIdxInTZ, firstBotIdx);
        require(isBotTeam(teamId), "cannot transfer a non-bot team");
        require(addr != NULL_ADDR, "invalid address");
        if ((firstBotIdx % TEAMS_PER_DIVISION) == (TEAMS_PER_DIVISION-1)) { _addDivision(tz, countryIdxInTZ); }
        teamIdToOwner[teamId] = addr;
        countryIdToNHumanTeams[countryId] = firstBotIdx + 1;
        emit TeamTransfer(teamId, addr);
    }
    
    /// Speeds up assignment to new users. 
    /// This function will crash if it cannot handle all transfers in one single TX
    /// It is the responsibility of the caller to ensure that the arrays length match correctly
    function transferFirstBotsToAddresses(uint8[] calldata tz, uint256[] calldata countryIdxInTZ, address[] calldata addr) external onlyRelay {
        for (uint256 i = 0; i < tz.length; i++) {
            transferFirstBotToAddr(tz[i], countryIdxInTZ[i], addr[i]); 
        }            
    }
    

    // Private Functions

    function _initTimeZone(uint8 tz) private {
        _orgMapRoot[tz][0] = INIT_ORGMAP_HASH;
        _addCountry(tz);
    }
    
    function _addCountry(uint8 tz) private {
        uint256 countryIdxInTZ = tzToNCountries[tz];
        tzToNCountries[tz] = countryIdxInTZ + 1;
        for (uint8 division = 0 ; division < DIVS_PER_LEAGUE_AT_START; division++){
            _addDivision(tz, countryIdxInTZ); 
        }
    }

    function _addDivision(uint8 tz, uint256 countryIdxInTZ) private {
        uint256 countryId = encodeTZCountryAndVal(tz, countryIdxInTZ, 0);
        uint256 nDivs = countryIdToNDivisions[countryId];
        uint256 divisionId = encodeTZCountryAndVal(tz, countryIdxInTZ, nDivs);
        countryIdToNDivisions[countryId] = nDivs + 1;
        divisionIdToRound[divisionId] = getCurrentRound(tz) + 1;
        emit DivisionCreation(tz, countryIdxInTZ, nDivs);
    }
}
