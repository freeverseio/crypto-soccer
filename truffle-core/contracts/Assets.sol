pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsView.sol";

/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Assets is AssetsView {

    event AssetsInit(address creatorAddr);
    event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry);
    
    function setCOO(address addr) public onlySuperUser { _COO = addr; }
    
    function setMarket(address addr) public onlySuperUser {
        _market = addr;
        teamIdToOwner[ACADEMY_TEAM] = addr;
        emit TeamTransfer(ACADEMY_TEAM, addr);        
    }
    
    function init() public onlyCOO {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        emit AssetsInit(msg.sender);
    }

    // hack for testing: we can init only one timezone
    // at some point, remove this option
    function initSingleTZ(uint8 tz) public onlyCOO {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        _initTimeZone(tz);
        emit AssetsInit(msg.sender);
    }
    

    function _initTimeZone(uint8 tz) private {
        _orgMapRoot[tz][0] = INIT_ORGMAP_HASH;
        addCountry(tz);
    }
    
    function addCountry(uint8 tz) public onlyCOO {
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

    function addDivisionManually(uint8 tz, uint256 countryIdxInTZ) external onlyCOO { _addDivision(tz, countryIdxInTZ); }

    // this function will crash if it cannot handle all transfers in one single TX
    // it is the responsibility of the caller to ensure that the arrays match correctly
    function transferFirstBotsToAddresses(uint8[] calldata tz, uint256[] calldata countryIdxInTZ, address[] calldata addr) external onlyMarket {
        for (uint256 i = 0; i < tz.length; i++) {
            transferFirstBotToAddr(tz[i], countryIdxInTZ[i], addr[i]); 
        }            
    }

    // Entry point for new users: acquiring a bot team
    function transferFirstBotToAddr(uint8 tz, uint256 countryIdxInTZ, address addr) public onlyMarket {
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

}
