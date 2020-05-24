pragma solidity >= 0.6.3;

import "./Storage.sol";
import "../encoders/EncodingIDs.sol";
import "../encoders/EncodingSkillsGetters.sol";

/**
 @title Library of View/Pure functions to needed by game assets and market
 @author Freeverse.io, www.freeverse.io
*/

contract AssetsLib is Storage, EncodingSkillsGetters, EncodingIDs {
    
    event TeamTransfer(uint256 teamId, address to);

    /// Modifiers for all functions that write to Storage
    
    modifier onlyMarket() {
        require(msg.sender == _market, "Only market owner is authorized.");
        _;
    }

    modifier onlyRelay() {
        require(msg.sender == _relay, "Only Relay owner is authorized.");
        _;
    }
    
    modifier onlyCOO() {
        require(msg.sender == _COO, "Only COO is authorized.");
        _;
    }
    
    /// Rest of view/pure functions 
    
    function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return getOwnerTeamInCountry(timeZone, countryIdxInTZ, teamIdxInCountry) == NULL_ADDR;
    }

    function isBotTeam(uint256 teamId) public view returns(bool) {
        if (teamId == ACADEMY_TEAM) return false;
        return teamIdToOwner[teamId] == NULL_ADDR;
    }

    /// returns NULL_ADDR if team is bot
    function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(address) {
        if (!_tzExists(timeZone) || !_countryInTZExists(timeZone, countryIdxInTZ)) return NULL_ADDR;
        return teamIdToOwner[encodeTZCountryAndVal(timeZone, countryIdxInTZ, teamIdxInCountry)];
    }
    
    function teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) public view returns(bool) {
        return (teamIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ));
    }

    function wasTeamCreatedVirtually(uint256 teamId) public view returns (bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) = decodeTZCountryAndVal(teamId);
        return teamExistsInCountry(timeZone, countryIdxInTZ, teamIdxInCountry);
    }
    
    function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        if (!_tzExists(timeZone) || !_countryInTZExists(timeZone, countryIdxInTZ)) return 0;
        return countryIdToNDivisions[encodeTZCountryAndVal(timeZone, countryIdxInTZ, 0)];
    }

    function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNDivisionsInCountry(timeZone, countryIdxInTZ) * LEAGUES_PER_DIV;
    }

    function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) public view returns(uint256) {
        return getNLeaguesInCountry(timeZone, countryIdxInTZ) * TEAMS_PER_LEAGUE;
    }
    
    function wasPlayerCreatedVirtually(uint256 playerId) public view returns(bool) {
        (uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        return _wasPlayerCreatedInCountry(timeZone, countryIdxInTZ, playerIdxInCountry);
    }

    function getCurrentRound(uint8 tz) public view returns (uint256) {
        return getCurrentRoundPure(tz, timeZoneForRound1, currentVerse);
    }

    function getCurrentRoundPure(uint8 tz, uint8 tz1, uint256 verse) public pure returns (uint256) { 
        /// first, compute "roundTZ1" for the first timezone that played a match
        /// first, ensure that round is always >= 1.
        if (verse < VERSES_PER_ROUND) return 0;
        uint256 roundTZ1 = verse / VERSES_PER_ROUND;
        /// Next, note that verses where this tz plays first matches of rounds:
        ///   verses(round) = deltaN * 4 + VERSES_PER_ROUND * round
        uint256 deltaN = (tz >= tz1) ? (tz - tz1) : ((tz + 24) - tz1);
        if (verse < 4 * deltaN + roundTZ1 * VERSES_PER_ROUND) {
            return roundTZ1 - 1;
        } else {
            return roundTZ1;
        }
    }
    
    function _countryInTZExists(uint8 timeZone, uint256 countryIdxInTZ) internal view returns(bool) {
        return(countryIdxInTZ < tzToNCountries[timeZone]);
    }

    function _wasPlayerCreatedInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 playerIdxInCountry) private view returns(bool) {
        return (playerIdxInCountry < getNTeamsInCountry(timeZone, countryIdxInTZ) * PLAYERS_PER_TEAM_INIT);
    }
    
    function market() public view returns (address) { return _market; }
    function COO() public view returns (address) { return _COO; }
    function relay() public view returns (address) { return _relay; }
    function cryptoMktAddr() public view returns (address) { return _cryptoMktAddr; }
    
    function _tzExists(uint8 timeZone) internal pure returns(bool) {
        return(timeZone > NULL_TIMEZONE && timeZone < 25);
    }
}
