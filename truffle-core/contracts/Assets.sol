pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsView.sol";

/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Assets is AssetsView {

    event AssetsInit(address creatorAddr);
    event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry);
    

    function setAcademyAddr(address addr) public {
        _academyAddr = addr;
        emit TeamTransfer(ACADEMY_TEAM, addr);        
    }
    
    function init() public {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        for (uint8 tz = 1; tz < 25; tz++) {
            _initTimeZone(tz);
        }
        emit AssetsInit(msg.sender);
    }

    // hack for testing: we can init only one timezone
    // at some point, remove this option
    function initSingleTZ(uint8 tz) public {
        require(gameDeployDay == 0, "cannot initialize twice");
        gameDeployDay = secsToDays(now);
        _initTimeZone(tz);
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

    function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) external {
        uint256 firstBotIdx = _timeZones[timeZone].countries[countryIdxInTZ].nHumanTeams;
        require(isBotTeamInCountry(timeZone, countryIdxInTZ, firstBotIdx), "cannot transfer a non-bot team");
        require(addr != NULL_ADDR, "invalid address");
        _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToOwner[firstBotIdx] = addr;
        _timeZones[timeZone].countries[countryIdxInTZ].nHumanTeams++;
        uint256 teamId = encodeTZCountryAndVal(timeZone, countryIdxInTZ, firstBotIdx);
        emit TeamTransfer(teamId, addr);
    }

}
