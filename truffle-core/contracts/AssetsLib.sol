pragma solidity >=0.5.12 <0.6.2;

import "./Storage.sol";
import "./EncodingIDs.sol";
// import "./EncodingState.sol";
import "./EncodingSkillsGetters.sol";
// import "./EncodingSkillsSetters.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract AssetsLib is Storage, EncodingSkillsGetters, EncodingIDs {
    
    address constant public NULL_ADDR = address(0);
    uint8 constant public LEAGUES_PER_DIV = 16;
    uint8 constant public TEAMS_PER_LEAGUE = 8;

    function _assertTZExists(uint8 timeZone) internal pure {
        require(timeZone > 0 && timeZone < 25, "timeZone does not exist");
    }
    
    function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry) == NULL_ADDR;
    }

    // returns NULL_ADDR if team is bot
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(address) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return _timeZones[timeZone].countries[countryIdxInTZ].teamIdxInCountryToTeam[teamIdxInCountry].owner;
    }

    function _assertCountryInTZExists(uint8 timeZone, uint256 countryIdxInTZ) internal view {
        require(countryIdxInTZ < _timeZones[timeZone].countries.length, "country does not exist in this timeZone");
    }
    
    function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return (teamIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ));
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
    
}
