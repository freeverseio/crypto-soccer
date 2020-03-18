pragma solidity >=0.5.12 <=0.6.3;

import "./UpdatesView.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Updates is UpdatesView {
    event TeamTransfer(uint256 teamId, address to);
    event ActionsSubmission(uint256 verse, uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime, bytes32 root, string ipfsCid);
    event TimeZoneUpdate(uint8 timeZone, bytes32 root, uint256 submissionTime);

    function initUpdates() public {
        require(timeZoneForRound1 == 0, "cannot initialize updates twice");
        // the game starts at verse = 0. The transition to verse = 1 will be at the next exact hour.
        // that will be the begining of Round = 1. So Round 1 starts at some timezone that depends on
        // the call to the contract init() function.
        // TZ = 1 => starts at 1:00... TZ = 23 => starts at 23:00, TZ = 24, starts at 0:00
        uint256 secsOfDay   = now % (3600 * 24);
        uint256 hour        = secsOfDay / 3600;  // 0, ..., 23
        uint256 minute      = (secsOfDay - hour * 3600) / 60; // 0, ..., 59
        uint256 secs        = (secsOfDay - hour * 3600 - minute * 60); // 0, ..., 59
        if (minute < 27) {
            timeZoneForRound1 = normalizeTZ(uint8(hour));
            nextVerseTimestamp = now + (29-minute)*60 + (60 - secs);
        } else {
            timeZoneForRound1 = normalizeTZ(1+uint8(hour));
            nextVerseTimestamp = now + (29-minute)*60 + (60 - secs) + 3600;
        }
    }
 
    function _incrementVerse() private {
        currentVerse += 1;
        nextVerseTimestamp += SECS_BETWEEN_VERSES;
    }
    
    function submitActionsRoot(bytes32 actionsRoot, string memory ipfsCid) public {
        require(now > nextVerseTimestamp, "too early to accept actions root");
        (uint8 newTZ, uint8 day, uint8 turnInDay) = nextTimeZoneToUpdate();
        (uint8 prevTz,,) = prevTimeZoneToUpdate();
        // make sure the last verse is settled
        if (prevTz != NULL_TIMEZONE) {
            require(now > getLastUpdateTime(prevTz)+ CHALLENGE_TIME, "last verse is still under challenge period");
        }
        if(newTZ != NULL_TIMEZONE) {
            _setActionsRoot(newTZ, actionsRoot);
        }
        _incrementVerse();
        _setCurrentVerseSeed(blockhash(block.number-1));
        emit ActionsSubmission(currentVerse, newTZ, day, turnInDay, blockhash(block.number-1), now, actionsRoot, ipfsCid);
    }
    
    function _setActionsRoot(uint8 timeZone, bytes32 root) public returns(uint256) {
        _assertTZExists(timeZone);
        actionsRoot[timeZone] = root;
        lastActionsSubmissionTime[timeZone] = now;
    }

    // accepts an update about the root of the current state of a timezone. 
    // in order to accept it, either:
    //  - timezone is null,
    //  - timezone has not been updated yet (lastUpdate < lastActionsSubmissionTime)
    function updateTZ(bytes32 root) public {
        (uint8 tz,,) = prevTimeZoneToUpdate();
        bool accept = (tz == NULL_TIMEZONE) || (getLastUpdateTime(tz) < getLastActionsSubmissionTime(tz));
        require(accept, "TZ has already been updated once");
        _setTZRoot(tz, root); // first time that we update this TZ
        emit TimeZoneUpdate(tz, root, now);
    }
    
    
    
    function _setTZRoot(uint8 tz, bytes32 root) internal returns(uint256) {
        uint8 newIdx = 1 - newestSkillsIdx[tz];
        newestSkillsIdx[tz] = newIdx;
        roots[tz][newIdx][0] = root;
        lastUpdateTime[tz] = now;
    }

    function _setCurrentVerseSeed(bytes32 seed) private {
        currentVerseSeed = seed;
    }


}
