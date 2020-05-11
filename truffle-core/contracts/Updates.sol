pragma solidity >=0.5.12 <=0.6.3;

import "./UpdatesBase.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Updates is UpdatesBase {
    event ActionsSubmission(uint256 verse, uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime, bytes32 root, string ipfsCid);
    event TimeZoneUpdate(uint256 verse, uint8 timeZone, bytes32 root, uint256 submissionTime);
    event ChallengeAccepted(uint8 tz, uint8 newLevel, bytes32 root, bytes32[] providedRoots);

    function setStakersAddress(address payable addr) public {
        _stakers = Stakers(addr);
    }

    function setChallengeTime(uint256 newTime) public { _challengeTime = newTime; }

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
        firstVerseTimeStamp = nextVerseTimestamp;
    }
 
    function _incrementVerse() private {
        currentVerse += 1;
        nextVerseTimestamp += SECS_BETWEEN_VERSES;
    }
    
    function submitActionsRoot(bytes32 actionsRoot, bytes32 activeTeamsPerCountryRoot, bytes32 orgMapRoot, uint8 levelVerifiableByBC, string memory ipfsCid) public {
        require(now > nextVerseTimestamp, "too early to accept actions root");
        (uint8 newTZ, uint8 day, uint8 turnInDay) = nextTimeZoneToUpdate();
        (uint8 prevTz,,) = prevTimeZoneToUpdate();
        // make sure the last verse is settled
        if (prevTz != NULL_TIMEZONE) {
            ( , , bool isSettled) = getStatus(prevTz, true);
            require(isSettled, "last verse is still under challenge period");
        }
        if(newTZ != NULL_TIMEZONE) {
            uint8 idx = 1 - _newestRootsIdx[newTZ];
            _newestRootsIdx[newTZ] = idx;
            _actionsRoot[newTZ][idx] = actionsRoot;
            _activeTeamsPerCountryRoot[newTZ][idx] = activeTeamsPerCountryRoot;
            _orgMapRoot[newTZ][idx] = orgMapRoot;
            _levelVerifiableByBC[newTZ][idx] = levelVerifiableByBC;
            _lastActionsSubmissionTime[newTZ] = now;
        }
        _incrementVerse();
        _setCurrentVerseSeed(blockhash(block.number-1));
        emit ActionsSubmission(currentVerse, newTZ, day, turnInDay, blockhash(block.number-1), now, actionsRoot, ipfsCid);
    }

    // accepts an update about the root of the current state of a timezone. 
    // in order to accept it, either:
    //  - timezone is null,
    //  - timezone has not been updated yet (lastUpdate < _lastActionsSubmissionTime)
    function updateTZ(bytes32 root) public {
        // when actionRoots were submitted, nextTimeZone points to the future.
        // so the timezone waiting for updates & challenges is provided by prevTimeZoneToUpdate()
        (uint8 tz,,) = prevTimeZoneToUpdate();
        bool accept = (tz != NULL_TIMEZONE) && (getLastUpdateTime(tz) < getLastActionsSubmissionTime(tz));
        require(accept, "TZ has already been updated once");
        _setTZRoot(tz, root); // first time that we update this TZ
        emit TimeZoneUpdate(currentVerse, tz, root, now);
        _stakers.update(0, msg.sender);
    }

    // TODO: specify which leaf you challenge!!! And bring Merkle proof!
    function challengeTZ(bytes32 challLeaveVal, uint256 challLeavePos, bytes32[] memory proofChallLeave, bytes32[] memory providedRoots) public {
        // intData = [tz, level, levelVerifiable, idx]
        uint8[4] memory intData = _cleanTimeAcceptedChallenges();
        bytes32 root = _assertFormallyCorrectChallenge(
            intData,
            challLeaveVal, 
            challLeavePos, 
            proofChallLeave, 
            providedRoots
        );
        require(intData[1] < intData[2] - 1, "this function must only be called for non-verifiable-by-BC challenges");        
        // accept the challenge and store new root, or let the BC verify challenge and revert to level - 1
        uint8 level = intData[1] + 1;
        _roots[intData[0]][intData[3]][level] = root;
        _challengeLevel[intData[0]][intData[3]] = level;
        emit ChallengeAccepted(intData[0], level, root, providedRoots);
        _lastUpdateTime[intData[0]] = now;
        _stakers.update(level, msg.sender);
    }

    function _setTZRoot(uint8 tz, bytes32 root) internal {
        uint8 idx = _newestRootsIdx[tz];
        _roots[tz][idx][0] = root;
        for (uint8 level = 1; level < MAX_CHALLENGE_LEVELS; level++) _roots[tz][idx][level] = 0;
        _lastUpdateTime[tz] = now;
    }

    function _setCurrentVerseSeed(bytes32 seed) private {
        currentVerseSeed = seed;
    }

    // TODO: remove this test function
    function setLevelVerifiableByBC(uint8 newVal) public {
        for (uint8 tz = 1; tz < 25; tz++) {
            _levelVerifiableByBC[tz][0] = newVal;
            _levelVerifiableByBC[tz][1] = newVal;
        }
    }
    
    function setChallengeLevels(uint16 levelsInOneChallenge, uint16 leafsInLeague, uint16 levelsInLastChallenge) public {
        _levelsInOneChallenge   = levelsInOneChallenge;
        _leafsInLeague          = leafsInLeague;
        _levelsInLastChallenge  = levelsInLastChallenge;
    }
    
}
