pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Updates {
    event TeamTransfer(uint256 teamId, address to);
    event ActionsSubmission(uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime);
    event TimeZoneUpdate(uint8 timeZone, bytes32 root, uint256 submissionTime);

    uint16 constant public SECS_BETWEEN_VERSES = 900; // 15 mins
    uint8 constant VERSES_PER_DAY = 96; // 24 * 4
    uint16 constant VERSES_PER_ROUND = 1536; // 96 * 16
    uint8 constant public NULL_TIMEZONE = 0;
    uint8 constant CHALLENGE_TIME = 60; // in secs
    
    uint256 public nextVerseTimestamp;
    uint8 public timeZoneForRound1;
    uint256 public currentVerse;
    bytes32 private _currentVerseSeed;
    bool private _needsInitUpdates = true;

    Assets private _assets;

    function initUpdates(address addr) public {
        require(_needsInitUpdates == true, "cannot initialize twice");
        _setAssetsAdress(addr);
        // the game starts at verse = 0. The transition to verse = 1 will be at the next exact hour.
        // that will be the begining of Round = 1. So Round 1 starts at some timezone that depends on
        // the call to the contract init() function.
        uint256 secsOfDay   = now % (3600 * 24);
        uint256 hour        = secsOfDay / 3600;  // 0, ..., 23
        uint256 minute      = (secsOfDay - hour * 3600) / 60; // 0, ..., 59
        uint256 secs        = (secsOfDay - hour * 3600 - minute * 60); // 0, ..., 59
        if (minute < 42) {
            timeZoneForRound1 = 1 + uint8(hour);
            nextVerseTimestamp = now + (44-minute)*60 + (60 - secs);
        } else {
            timeZoneForRound1 = 2 + uint8(hour);
            nextVerseTimestamp = now + (44-minute)*60 + (60 - secs) + 3600;
        }
        _needsInitUpdates = false;
    }
 
    function _setAssetsAdress(address addr) private {
        _assets = Assets(addr);
    }

    function getNow() public view returns(uint256) {
        return now;
    }

    function incrementVerse() private {
        currentVerse += 1;
        nextVerseTimestamp += SECS_BETWEEN_VERSES;
    }
    
    function submitActionsRoot(bytes32 actionsRoot) public {
        require(now > nextVerseTimestamp, "too early to accept actions root");
        (uint8 newTZ, uint8 day, uint8 turnInDay) = nextTimeZoneToUpdate();
        (uint8 prevTz,,) = prevTimeZoneToUpdate();
        // make sure the last verse is settled
        if (prevTz != NULL_TIMEZONE) {
            require(now > _assets.getLastUpdateTime(prevTz)+ CHALLENGE_TIME, "last verse is still under challenge period");
        }
        // if we are precisely at a moment where nothing needs to be done, just move ahead
        if (newTZ == NULL_TIMEZONE) {
            incrementVerse() ;
            emit ActionsSubmission(NULL_TIMEZONE, 0, 0, 0, now);
            return;
        }
        _assets.setActionsRoot(newTZ, actionsRoot);
        incrementVerse() ;
        setCurrentVerseSeed(blockhash(block.number-1)); 
        emit ActionsSubmission(newTZ, day, turnInDay, blockhash(block.number-1), now);
    }

    function updateTZ(bytes32 root) public {
        (uint8 tz,,) = prevTimeZoneToUpdate();
        require(tz != NULL_TIMEZONE, "nothing to update in the current timeZone");
        uint256 lastUpdate = _assets.getLastUpdateTime(tz);
        uint256 lastActionsSubmissionTime = _assets.getLastActionsSubmissionTime(tz);
        if (lastUpdate > lastActionsSubmissionTime) {
            require(now < lastUpdate + CHALLENGE_TIME, "challenging period is already over for this timezone");
        } else {
            require(now < lastActionsSubmissionTime + CHALLENGE_TIME, "challenging period is already over for this timezone");
        }
        _assets.setSkillsRoot(tz, root);
        emit TimeZoneUpdate(tz, root, now);
    }
    
    // each day has 24 hours, each with 4 verses => 96 verses per day.
    // day = 1,..16
    // turnInDay = 0, 1, 2, 3
    // so for each TZ, we go from (day, turn) = (1, 0) ... (15,3) => a total of 16*4 = 64 turns per timeZone
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
        if (verseInRound < VERSES_PER_DAY) {
            timeZone = 1 + uint8((TZForRound1 - 1 + (verseInRound / 4))% 24);
            day = 1;
            turnInDay = uint8(verseInRound % 4);
        } else if (verseInRound == VERSES_PER_DAY) {
            timeZone = NULL_TIMEZONE;
        } else {
            timeZone = 1 + uint8((TZForRound1 - 1 + ((verseInRound - 1) / 4))% 24);
            day = 1 + uint8((verseInRound - 1) / VERSES_PER_DAY);
            turnInDay = uint8((verseInRound - 1) % 4);
        }
    }

    function setCurrentVerseSeed(bytes32 seed) public {
        _currentVerseSeed = seed;
    }
        
    function getCurrentVerseSeed() public view returns (bytes32) {
        return _currentVerseSeed;
    }
        

}
