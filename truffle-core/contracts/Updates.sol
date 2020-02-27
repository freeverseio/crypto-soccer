pragma solidity >=0.5.12 <=0.6.3;

import "./UpdatesView.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Updates is UpdatesView {
    event TeamTransfer(uint256 teamId, address to);
    event ActionsSubmission(uint256 verse, uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime, bytes32 root, string cid);
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
 
    function incrementVerse() private {
        currentVerse += 1;
        nextVerseTimestamp += SECS_BETWEEN_VERSES;
    }
    
    function submitActionsRoot(bytes32 actionsRoot, string memory cid) public {
        require(now > nextVerseTimestamp, "too early to accept actions root");
        (uint8 newTZ, uint8 day, uint8 turnInDay) = nextTimeZoneToUpdate();
        // (uint8 prevTz,,) = prevTimeZoneToUpdate();
        // make sure the last verse is settled
        // if (prevTz != NULL_TIMEZONE) {
        //     require(now > _storageProxy.getLastUpdateTime(prevTz)+ CHALLENGE_TIME, "last verse is still under challenge period");
        // }
        if(newTZ != NULL_TIMEZONE) {
            setActionsRoot(newTZ, actionsRoot);
        }
        incrementVerse();
        setCurrentVerseSeed(blockhash(block.number-1));
        emit ActionsSubmission(currentVerse, newTZ, day, turnInDay, blockhash(block.number-1), now, actionsRoot, cid);
    }
    
    function setActionsRoot(uint8 timeZone, bytes32 root) public returns(uint256) {
        _assertTZExists(timeZone);
        _timeZones[timeZone].actionsRoot = root;
        _timeZones[timeZone].lastActionsSubmissionTime = now;
    }

    function updateTZ(bytes32 root) public {
        (uint8 tz,,) = prevTimeZoneToUpdate();
        if(tz != NULL_TIMEZONE) {
            uint256 lastUpdate = getLastUpdateTime(tz);
            uint256 lastActionsSubmissionTime = getLastActionsSubmissionTime(tz);
            if (lastUpdate > lastActionsSubmissionTime) {
                require(now < lastUpdate + CHALLENGE_TIME, "challenging period is already over for this timezone");
                setSkillsRoot(tz, root, false); // this is a challenge to a previous update
            } else {
                require(now < lastActionsSubmissionTime + CHALLENGE_TIME, "challenging period is already over for this timezone");
                setSkillsRoot(tz, root, true); // first time that we update this TZ
            }
        }
        emit TimeZoneUpdate(tz, root, now);
    }
    
    function setSkillsRoot(uint8 tz, bytes32 root, bool newTZ) internal returns(uint256) {
        if (newTZ) _timeZones[tz].newestSkillsIdx = 1 - _timeZones[tz].newestSkillsIdx;
        _timeZones[tz].skillsHash[_timeZones[tz].newestSkillsIdx] = root;
        _timeZones[tz].lastUpdateTime = now;
    }

    function setCurrentVerseSeed(bytes32 seed) public {
        currentVerseSeed = seed;
    }

}
