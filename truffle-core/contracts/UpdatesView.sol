pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsLib.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract UpdatesView is Storage, AssetsLib {

    function getNow() public view returns(uint256) {
        return now;
    }

    function getLastUpdateTime(uint8 timeZone) internal view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastUpdateTime;
    }
    
    function getLastActionsSubmissionTime(uint8 timeZone) public view returns(uint256) {
        _assertTZExists(timeZone);
        return _timeZones[timeZone].lastActionsSubmissionTime;
    }

    
    // each day has 24 hours, each with 4 verses => 96 verses per day.
    // day = 0,..13
    // turnInDay = 0, 1, 2, 3
    // so for each TZ, we go from (day, turn) = (0, 0) ... (13,3) => a total of 14*4 = 56 turns per timeZone
    // from these, all map easily to timeZones
    function nextTimeZoneToUpdate() public view returns (uint8 timeZone, uint8 day, uint8 turnInDay) {
        return _timeZoneToUpdatePure(currentVerse, timeZoneForRound1);
    }

    function prevTimeZoneToUpdate() public view returns (uint8 timeZone, uint8 day, uint8 turnInDay) {
        if (currentVerse == 0) {
            return (NULL_TIMEZONE, 0, 0);
        }
        return _timeZoneToUpdatePure(currentVerse - 1, timeZoneForRound1);
    }

    function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) public pure returns (uint8 timeZone, uint8 day, uint8 turnInDay) {
        // if currentVerse = 0, we should be updating timeZoneForRound1
        // recall that timeZones range from 1...24 (not from 0...24)
        uint16 verseInRound = uint16(verse % VERSES_PER_ROUND);
        turnInDay = uint8(verseInRound % 4);
        if (turnInDay < 2) {
            timeZone = normalizeTZ(TZForRound1 + verseInRound/4);
            day = 2 * uint8(verseInRound / VERSES_PER_DAY);
        } else {
            timeZone = normalizeTZ(TZForRound1 + 9 + verseInRound/4);
            day = 1 + 2 * uint8(verseInRound / VERSES_PER_DAY);
            turnInDay -= 2;
        }
        
    }
    
    function normalizeTZ(uint16 tz) public pure returns (uint8) {
        return uint8(1 + ((tz - 1)% 24));
    }

    function getNextVerseTimestamp() public view returns (uint256) { return nextVerseTimestamp; }
    function getTimeZoneForRound1() public view returns (uint8) { return timeZoneForRound1; }
    function getCurrentVerse() public view returns (uint256) { return currentVerse; }
    function getCurrentVerseSeed() public view returns (bytes32) { return currentVerseSeed; }

    // function getTimeZonesPlayingInNextRound() public view returns (uint8[] memory) {
    //     // morning 11am - evening 20 pm
    //     // so if the first match is in the morning => 11am, then the next one is at +9h
    //     uint8[] memory tzs = new uint8[](2);
    //     (tzs[0], , ) = nextTimeZoneToUpdate();
    //     tzs[1] = normalizeTZ(tzs[0] - 9);
    //     return tzs; 
    // }

}
