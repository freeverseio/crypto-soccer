pragma solidity >=0.5.12 <=0.6.3;

import "./Storage.sol";
import "./EncodingIDs.sol";
import "./EncodingSkillsGetters.sol";
/**
 * @title Entry point for changing ownership of assets, and managing bids and auctions.
 */

contract AssetsLib is Storage, EncodingSkillsGetters, EncodingIDs {
    
    event TeamTransfer(uint256 teamId, address to);

    function _assertTZExists(uint8 timeZone) internal pure {
        require(timeZone > NULL_TIMEZONE && timeZone < 25, "timeZone does not exist");
    }
    
    function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry) == NULL_ADDR;
    }

    function isBotTeam(uint256 teamId) public view returns(bool) {
        if (teamId == ACADEMY_TEAM) return false;
        return teamIdToOwner[teamId] == NULL_ADDR;
    }

    // returns NULL_ADDR if team is bot
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(address) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return teamIdToOwner[encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry)];
    }

    function _assertCountryInTZExists(uint8 timeZone, uint256 countryIdxInTZ) internal view {
        require(countryIdxInTZ < tzToNCountries[timeZone], "country does not exist in this timeZone");
    }
    
    function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return (teamIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ));
    }

    function teamExists(uint256 teamId) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return _teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }
    
    function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        _assertTZExists(timeZone);
        _assertCountryInTZExists(timeZone, countryIdxInTZ);
        return countryIdToNDivisions[encodeTZCountryAndVal(timeZone, countryIdxInTZ, 0)];
    }

    function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNDivisionsInCountry(timeZone, countryIdxInTZ) * LEAGUES_PER_DIV;
    }

    function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNLeaguesInCountry(timeZone, countryIdxInTZ) * TEAMS_PER_LEAGUE;
    }
    
    function playerExists(uint256 playerId) public view returns (bool) {
        if (playerId == 0) return false;
        if (getIsSpecial(playerId)) return (_playerIdToState[playerId] != 0);
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        return _wasPlayerCreatedInCountry(timeZone, countryIdxInTZ, playerIdxInCountry);
    }

    function _wasPlayerCreatedInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) private view returns(bool) {
        return (playerIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ) * PLAYERS_PER_TEAM_INIT);
    }
    
    function getCurrentRound(uint8 tz) public view returns (uint256) {
        return getCurrentRoundPure(tz, timeZoneForRound1, currentVerse);
    }

    function getCurrentRoundPure(uint8 tz, uint8 tz1, uint256 verse) public pure returns (uint256) { 
        // first, compute "roundTZ1" for the first timezone that played a match
        // first, ensure that round is always >= 1.
        if (verse < VERSES_PER_ROUND) return 0;
        uint256 roundTZ1 = verse / VERSES_PER_ROUND;
        // Next, note that verses where this tz plays first matches of rounds:
        //   verses(round) = deltaN * 4 + VERSES_PER_ROUND * round
        uint256 deltaN = (tz >= tz1) ? (tz - tz1) : ((tz + 24) - tz1);
        if (verse < 4 * deltaN + roundTZ1 * VERSES_PER_ROUND) {
            return roundTZ1 - 1;
        } else {
            return roundTZ1;
        }
    }

}
