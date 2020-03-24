pragma solidity >=0.5.12 <=0.6.3;

import "./UpdatesView.sol";
import "./Merkle.sol";
 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Updates is UpdatesView, Merkle {
    event TeamTransfer(uint256 teamId, address to);
    event ActionsSubmission(uint256 verse, uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime, bytes32 root, string ipfsCid);
    event TimeZoneUpdate(uint8 timeZone, bytes32 root, uint256 submissionTime);
    event ChallengeAccepted(uint8 tz, uint8 newLevel, bytes32 root, bytes32[] providedRoots);
    event ChallengeResolved(uint8 tz, uint8 resolvedLevel, bool isSuccessful);

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
            ( , , bool isSettled) = getStatus(prevTz, true);
            require(isSettled, "last verse is still under challenge period");
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
        // when actionRoots were submitted, nextTimeZone points to the future.
        // so the timezone waiting for updates & challenges is provided by prevTimeZoneToUpdate()
        (uint8 tz,,) = prevTimeZoneToUpdate();
        bool accept = (tz == NULL_TIMEZONE) || (getLastUpdateTime(tz) < getLastActionsSubmissionTime(tz));
        require(accept, "TZ has already been updated once");
        _setTZRoot(tz, root); // first time that we update this TZ
        emit TimeZoneUpdate(tz, root, now);
    }

    // TODO: specify which leaf you challenge!!! And bring Merkle proof!
    function challengeTZ(bytes32 challLeaveVal, uint256 challLeavePos, bytes32[] memory proofChallLeave, bytes32[] memory providedRoots, bool forceSuccess) public {
        (uint8 tz,,) = prevTimeZoneToUpdate();
        require(tz != NULL_TIMEZONE, "cannot challenge the null timezone");
        require(now < getLastUpdateTime(tz) + CHALLENGE_TIME, "challenging time is over for the current timezone");
        bytes32 root = merkleRoot(providedRoots, LEVELS_IN_ONE_CHALLENGE);
        (uint8 newIdx, uint8 level, uint8 levelVerifiable) = getChallengeData(tz, true);
        level = _cleanTimeAcceptedChallenges(tz, level);
        // verify provided roots are an actual challenge (they lead to a root different from the one provided by previous challenge/update)
    
        if (level == 0) require(root != getRoot(tz, 0, true), "provided leafs lead to same root being challenged");
        else {
            require(root != challLeaveVal, "you are declaring that the provided leafs lead to same root being challenged");
            bytes32 prevRoot = getRoot(tz, level, true);
            require(verify(prevRoot, proofChallLeave, challLeaveVal, challLeavePos), "merkle proof not correct");
        }
        // accept the challenge and store new root, or let the BC verify challenge and revert to level - 1
        if (level < levelVerifiable - 1) {
           level += 1;
            _roots[tz][newIdx][level] = root;
            challengeLevel[tz][newIdx] = level;
            emit ChallengeAccepted(tz, level, root, providedRoots);
        } else {
            // bool success = computeChallenge(challLeaveVal, challLeavePos, providedRoots);
            bool success = forceSuccess;
            require(success, "challenge was not successful according to blockchain computation");
            _roots[tz][newIdx][level] = 0;
            challengeLevel[tz][newIdx] = level - 1;
            emit ChallengeResolved(tz, level + 1, true);
            emit ChallengeResolved(tz, level, false);
        }
        lastUpdateTime[tz] = now;
    }
    
    function _cleanTimeAcceptedChallenges(uint8 tz, uint8 writtenLevel) internal returns (uint8) {
        (uint8 finalLevel, uint8 nJumps, ) = getStatus(tz, true);
        // if there was 0 jumps, do nothing
        if (nJumps == 0) return writtenLevel;
        // otherwise clean all data except for the lowest level
        require(writtenLevel == finalLevel + 2 * nJumps, "challenge status: nJumps incompatible with writtenLevel and finalLevel");
        uint8 idx = newestRootsIdx[tz];
        for (uint8 j = 0; j < nJumps; j++) {
            uint8 levelAccepted = finalLevel + 2 * (j + 1);
            _roots[tz][idx][levelAccepted] = 0;
            _roots[tz][idx][levelAccepted-1] = 0;
            emit ChallengeResolved(tz, levelAccepted, true);
            emit ChallengeResolved(tz, levelAccepted - 1, false);
        }
        challengeLevel[tz][idx] = finalLevel;
        lastUpdateTime[tz] = now;
        return finalLevel;
    }
    
    function _setTZRoot(uint8 tz, bytes32 root) internal {
        uint8 newIdx = 1 - newestRootsIdx[tz];
        newestRootsIdx[tz] = newIdx;
        _roots[tz][newIdx][0] = root;
        for (uint8 level = 1; level < MAX_CHALLENGE_LEVELS; level++) _roots[tz][newIdx][level] = 0;
        levelVerifiableByBC[tz][newIdx] = 3;
        lastUpdateTime[tz] = now;
    }

    function _setCurrentVerseSeed(bytes32 seed) private {
        currentVerseSeed = seed;
    }


}
