pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsLib.sol";
 /**
 * @title Entry point to submit user actions, and tz root updates, which makes time evolve.
 */

contract UpdatesView is Storage, AssetsLib {

    function getNow() public view returns(uint256) {
        return now;
    }

    function getLastUpdateTime(uint8 tz) public view returns(uint256) {
        _assertTZExists(tz);
        return _lastUpdateTime[tz];
    }
    
    function getLastActionsSubmissionTime(uint8 tz) public view returns(uint256) {
        _assertTZExists(tz);
        return _lastActionsSubmissionTime[tz];
    }

    
    // each day has 24 hours, each with 4 verses => 96 verses per day.
    // day = 0,..13
    // turnInDay = 0, 1, 2, 3
    // so for each TZ, we go from (day, turn) = (0, 0) ... (13,3) => a total of 14*4 = 56 turns per tz
    // from these, all map easily to timeZones
    function nextTimeZoneToUpdate() public view returns (uint8 tz, uint8 day, uint8 turnInDay) {
        return _timeZoneToUpdatePure(currentVerse, timeZoneForRound1);
    }

    function prevTimeZoneToUpdate() public view returns (uint8 tz, uint8 day, uint8 turnInDay) {
        if (currentVerse == 0) {
            return (NULL_TIMEZONE, 0, 0);
        }
        return _timeZoneToUpdatePure(currentVerse - 1, timeZoneForRound1);
    }


    // tz0  : v = 0, V_DAY, 2 * V_DAY...
    // tzN  : v = 4N + V_DAY * day,  day = 0,...6
    // tzN  : v = 4N + V_DAY * day,  day = 0,...6
    //  => tzN - tz0 = (v - V_DAY*day)
    //  => 4 tzN = 4 tz0 + v % VERSES_PER_DAY
    // last : v = V_DAY + DELTA + V_DAY * 6 
    // Imagine 2 tzs:
    // 0:00 - tz0; 0:30 - NUL; 1:00 - tz1; 1:30 - tz0; 0:00 - tz0; 0:30 - tz1;
    // So the last
    function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) public pure returns (uint8 timezone, uint8 day, uint8 turnInDay) {
        // if currentVerse = 0, we should be updating timeZoneForRound1
        // recall that timeZones range from 1...24 (not from 0...24)
        turnInDay = uint8(verse % 4);
        uint256 delta = 9 * 4 + turnInDay;
        uint256 tz;
        uint256 dia;        
        if (turnInDay >=2 && verse < delta) return (NULL_TIMEZONE, 0, 0);
        if (turnInDay < 2) {
            tz = TZForRound1 + ((verse - turnInDay) % VERSES_PER_DAY)/4;
            dia = 2 * uint8((verse - 4 * (tz - TZForRound1) - turnInDay)/VERSES_PER_DAY);
        } else {
            tz = TZForRound1 + ((verse - delta) % VERSES_PER_DAY)/4;
            dia = 1 + 2 * uint8((verse - 4 * (tz - TZForRound1) - delta)/VERSES_PER_DAY);
            turnInDay -= 2;
        }
        timezone = normalizeTZ(tz);
        day = uint8(dia % MATCHDAYS_PER_ROUND);
    }
    
    function normalizeTZ(uint256 tz) public pure returns (uint8) {
        return uint8(1 + ((24 + tz - 1)% 24));
    }

    function getNextVerseTimestamp() public view returns (uint256) { return nextVerseTimestamp; }
    function getTimeZoneForRound1() public view returns (uint8) { return timeZoneForRound1; }
    function getCurrentVerse() public view returns (uint256) { return currentVerse; }
    function getCurrentVerseSeed() public view returns (bytes32) { return currentVerseSeed; }

    function getRoot(uint8 tz, uint8 level, bool current) public view returns(bytes32) { 
        return (current) ? _roots[tz][_newestRootsIdx[tz]][level] : _roots[tz][1-_newestRootsIdx[tz]][level];
    }

    function getChallengeData(uint8 tz, bool current) public view returns(uint8, uint8, uint8) { 
        uint8 idx = current ? _newestRootsIdx[tz] : 1 - _newestRootsIdx[tz];
        return (idx, _challengeLevel[tz][idx], _levelVerifiableByBC);
    }

    function getStatus(uint8 tz, bool current) public view returns(uint8, uint8, bool) { 
        uint8 idx = current ? _newestRootsIdx[tz] : 1 - _newestRootsIdx[tz];
        uint8 writtenLevel = _challengeLevel[tz][idx];
        return getStatusPure(now, _lastUpdateTime[tz], writtenLevel);
    }
    
    function getStatusPure(uint256 nowTime, uint256 lastUpdate, uint8 writtenLevel) public pure returns(uint8 finalLevel, uint8 nJumps, bool isSettled) {
        uint256 numChallPeriods = (nowTime > lastUpdate) ? (nowTime - lastUpdate)/CHALLENGE_TIME : 0;
        finalLevel = (writtenLevel >= 2 * numChallPeriods) ? uint8(writtenLevel - 2 * numChallPeriods) : (writtenLevel % 2);
        nJumps = (writtenLevel - finalLevel) / 2;
        isSettled = nowTime > lastUpdate + (nJumps + 1) * CHALLENGE_TIME;
    }

}
