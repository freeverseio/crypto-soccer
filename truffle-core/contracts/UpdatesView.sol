pragma solidity >=0.5.12 <=0.6.3;

import "./AssetsLib.sol";
 /**
 * @title Entry point to submit user actions, and tz root updates, which makes time evolve.
 */

contract UpdatesView is AssetsLib {

    function getNow() public view returns(uint256) {
        return now;
    }

    function getLastUpdateTime(uint8 tz) public view returns(uint256) {
        _tzExists(tz);
        return _lastUpdateTime[tz];
    }
    
    function getLastActionsSubmissionTime(uint8 tz) public view returns(uint256) {
        _tzExists(tz);
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

    function getChallengeTime() public view returns (uint256) { return _challengeTime; }

    function getChallengeData(uint8 tz, bool current) public view returns(uint8, uint8, uint8) { 
        uint8 idx = current ? _newestRootsIdx[tz] : 1 - _newestRootsIdx[tz];
        return (idx, _challengeLevel[tz][idx], _levelVerifiableByBC[tz][idx]);
    }

    function getStatus(uint8 tz, bool current) public view returns(uint8, uint8, bool) { 
        uint8 idx = current ? _newestRootsIdx[tz] : 1 - _newestRootsIdx[tz];
        uint8 writtenLevel = _challengeLevel[tz][idx];
        return getStatusPure(now, _lastUpdateTime[tz], _challengeTime, writtenLevel);
    }
    
    function isTimeToUpdate() public view returns(bool) {
        (uint8 tz,,) = prevTimeZoneToUpdate();
        if (tz == NULL_TIMEZONE) return true;
        if (!(getLastUpdateTime(tz) < getLastActionsSubmissionTime(tz))) return false;
        (,, bool isSettled) = getStatus(tz, true);
        return isSettled;
    }
    

    
    function getStatusPure(uint256 nowTime, uint256 lastUpdate, uint256 challengeTime, uint8 writtenLevel) public pure returns(uint8 finalLevel, uint8 nJumps, bool isSettled) {
        if (challengeTime == 0) return (writtenLevel, 0, nowTime > lastUpdate);
        uint256 numChallPeriods = (nowTime > lastUpdate) ? (nowTime - lastUpdate)/challengeTime : 0;
        finalLevel = (writtenLevel >= 2 * numChallPeriods) ? uint8(writtenLevel - 2 * numChallPeriods) : (writtenLevel % 2);
        nJumps = (writtenLevel - finalLevel) / 2;
        isSettled = nowTime > lastUpdate + (nJumps + 1) * challengeTime;
    }
    
    // tz(n0)   : 11.30 = 0
    //          : 21.00 = 11.30 + 9.30h
    // tz(n)    : 11.30 + (n-n0)*1h
    //          : 21.00 + (n-n0)*1h
    //          : 11.30 + (n-n0)*1h + 24h * day (day = mDay/2)
    //              = 11.30 + ( deltaN + 12 * mDay ) * 1h
    //          : 21.00 + (n-n0) + 24h * day (day = (mDay-1)/2)
    //              = 11.30 + (9.5 + deltaN + 12 * (mDay-1) ) * 1h
    // add round * T_round = round * 7 * 24 * 3600 = round * DAYS_PER_ROUND * 24 * 1h
    // 
    // if even: 11.30 + ( deltaN + 12 * mDay + 24 * round * DAYS_PER_ROUND ) * 1h
    // if odd:  11.30 + (9.5 + deltaN + 12 * (mDay-1) + 24 * round * DAYS_PER_ROUND ) * 1h
    //        = 11.30 + (19 + 2*deltaN + 24 * (mDay-1) + 48 * round * DAYS_PER_ROUND ) * (1h/2)
    
    function getMatchUTC(uint8 tz, uint256 round, uint256 matchDay) public view returns(uint256 timeUTC) {
        require(tz > 0 && tz < 25, "timezone out of range");
        uint256 deltaN = (tz >= timeZoneForRound1) ? (tz - timeZoneForRound1) : ((tz + 24) - timeZoneForRound1);
        if (matchDay % 2 == 0) {
            return firstVerseTimeStamp + (deltaN + 12 * matchDay + 24 * DAYS_PER_ROUND * round) * 3600;
        } else {
            return firstVerseTimeStamp + (19 + 2*deltaN + 24 * (matchDay-1) + 48 * DAYS_PER_ROUND * round) * 1800;
        }
    }

    function getMatchUTCInCurrentRound(uint8 tz, uint256 matchDay) public view returns(uint256 timeUTC) {
        return getMatchUTC(tz, getCurrentRound(tz), matchDay);
    }
    
    function getAllMatchdaysUTCInRound(uint8 tz, uint256 round) public view returns(uint256[MATCHDAYS_PER_ROUND] memory timesUTC) {
        for (uint8 m = 0; m < MATCHDAYS_PER_ROUND; m++) timesUTC[m] = getMatchUTC(tz, round, m);
    }

    function getAllMatchdaysUTCInCurrentRound(uint8 tz) public view returns(uint256[MATCHDAYS_PER_ROUND] memory timesUTC) {
        for (uint8 m = 0; m < MATCHDAYS_PER_ROUND; m++) timesUTC[m] = getMatchUTC(tz, getCurrentRound(tz), m);
    }

}
