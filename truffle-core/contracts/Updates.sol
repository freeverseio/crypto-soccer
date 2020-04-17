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
        bool accept = (tz == NULL_TIMEZONE) || (getLastUpdateTime(tz) < getLastActionsSubmissionTime(tz));
        require(accept, "TZ has already been updated once");
        _setTZRoot(tz, root); // first time that we update this TZ
        emit TimeZoneUpdate(tz, root, now);
    }

    // TODO: specify which leaf you challenge!!! And bring Merkle proof!
    function challengeTZ(bytes32 challLeaveVal, uint256 challLeavePos, bytes32[] memory proofChallLeave, bytes32[] memory providedRoots) public {
        // intData = [tz, level, levelVerifiable, idx]
        (bytes32 root , uint8[4] memory intData) = _assertFormallyCorrectChallenge(
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
    }

    function BCVerifableChallengeFake(bytes32[] memory leagueLeafs, bool forceSuccess) public {
        // intData = [tz, level, levelVerifiable, idx]
        ( , uint8[4] memory intData) = _assertFormallyCorrectChallenge(0, 0, new bytes32[](0) , leagueLeafs);
        require(intData[1] == intData[2] - 1, "this function must only be called for non-verifiable-by-BC challenges"); 
        require(forceSuccess, "fake challenge failed because it was told to fail");
        _completeSuccessfulVerifiableChallenge(intData);
    }

    function BCVerifableChallengeZeros(bytes32[] memory leagueLeafs) public {
        // intData = [tz, level, levelVerifiable, idx]
        // PROBLEM: leagueleafs changes in the Merkle root!!
        ( , uint8[4] memory intData) = _assertFormallyCorrectChallenge(0, 0, new bytes32[](0), leagueLeafs);
        require(intData[1] == intData[2] - 1, "this function must only be called for non-verifiable-by-BC challenges"); 

        (, uint8 day, uint8 half) = prevTimeZoneToUpdate();
        require(day==0,"--");
        require(half==0,"--++");
        require(areThereUnexpectedZeros(leagueLeafs, day, half), "challenge to unexpected zeros failed");
        _completeSuccessfulVerifiableChallenge(intData);
    }
    
    // check that leagueLeafs.length == 640 has been done before calling this function (to preserve it being pure)
    function areThereUnexpectedZeros(bytes32[] memory leagueLeafs, uint8 day, uint8 half) public pure returns(bool) {
        if ((day == 0) && (half == 0)) {
            // at end of 1st half we still do not have league points
            for (uint16 i = 0; i < 8; i++) {
                if (leagueLeafs[i] != 0) return true;
            }
            // we do not have tactics, nor training, nor ML before
            for (uint16 team = 0; team < TEAMS_PER_LEAGUE; team++) {
                uint16 off = 128 + 64 * team;
                for (uint16 i = 25; i < 28; i++) {
                    if (leagueLeafs[off + i] != 0) return true;
                }
            }
        }
        // every element of team from 28 to 32 is 0
        for (uint16 team = 0; team < TEAMS_PER_LEAGUE; team++) {
            uint16 off = 128 + 64 * team;
            for (uint16 i = 28; i < 32; i++) {
                if (leagueLeafs[off + i] != 0) return true;
                if (leagueLeafs[off+ 32 + i] != 0) return true;
            }
        }
        // no goals after this day
        uint16 off = 8 + 8 * day;
        if (half == 1) off += 8;
        for (uint16 i = off; i < 128; i++) {
            if (leagueLeafs[i] != 0) return true;
        }
        return false;
    }
    

    function _completeSuccessfulVerifiableChallenge(uint8[4] memory intData) internal {
        // intData = [tz, level, levelVerifiable, idx]
        _roots[intData[0]][intData[3]][intData[1]] = 0;
        _challengeLevel[intData[0]][intData[3]] = intData[1] - 1;
        emit ChallengeResolved(intData[0], intData[1] + 1, true);
        emit ChallengeResolved(intData[0], intData[1], false);
    }

    function _assertFormallyCorrectChallenge(
        bytes32 challLeaveVal, 
        uint256 challLeavePos, 
        bytes32[] memory proofChallLeave, 
        bytes32[] memory providedRoots
    ) 
        private 
        returns (bytes32, uint8[4] memory intData)
    {
        // intData = [tz, level, levelVerifiable, idx]
        intData = _cleanTimeAcceptedChallenges();

        // build the root of the providedData
        bytes32 root;
        if (intData[1] + 2 >= intData[2]) {
            require(providedRoots.length == _leafsInLeague, "league leafs must have len 640");
            root = merkleRoot(providedRoots, _levelsInLastChallenge);
        } else {
            require(providedRoots.length == 2**uint256(_levelsInOneChallenge), "league leafs must have len 640");
            root = merkleRoot(providedRoots, _levelsInOneChallenge);
        }

        // We first check that the provided roots are an actual challenge,
        // hence leading to a root different from the one provided by previous challenge/update)
        if (intData[1] == 0) {
            // at level 0, the value one challenges is the one written in the BC, so we don't use challLeaveVal (could be anything)
            // and we don't need to verifiy that it belonged to a previous commit.
            require(root != getRoot(intData[0], 0, true), "provided leafs lead to same root being challenged");
        } else if ((intData[1] + 1) == intData[2]) {
            // at last level, we just provide the league leaves provided by the last challenger,
            // and we verify that they DO match with what is written.
            require(root == getRoot(intData[0], intData[1], true), "provided leafs lead to same root being challenged");
        } else {
            // otherwise we also check that the challVal belonged to a previous commit
            require(root != challLeaveVal, "you are declaring that the provided leafs lead to same root being challenged");
            bytes32 prevRoot = getRoot(intData[0], intData[1], true);
            require(verify(prevRoot, proofChallLeave, challLeaveVal, challLeavePos), "merkle proof not correct");
        }
        return (root, intData);        
    }
    
    function _cleanTimeAcceptedChallenges() internal returns (uint8[4] memory intData) {
        // intData = [tz, level, levelVerifiable, idx]
        (intData[0],,) = prevTimeZoneToUpdate();
        require(intData[0] != NULL_TIMEZONE, "cannot challenge the null timezone");
        (intData[3], intData[1], intData[2]) = getChallengeData(intData[0], true);

        (uint8 finalLevel, uint8 nJumps, bool isSettled) = getStatus(intData[0], true);
        require(!isSettled, "challenging time is over for the current timezone");
        // if there was 0 jumps, do nothing
        if (nJumps == 0) return intData;
        // otherwise clean all data except for the lowest level
        require(intData[1] == finalLevel + 2 * nJumps, "challenge status: nJumps incompatible with writtenLevel and finalLevel");
        uint8 idx = _newestRootsIdx[intData[0]];
        for (uint8 j = 0; j < nJumps; j++) {
            uint8 levelAccepted = finalLevel + 2 * (j + 1);
            _roots[intData[0]][idx][levelAccepted] = 0;
            _roots[intData[0]][idx][levelAccepted-1] = 0;
            emit ChallengeResolved(intData[0], levelAccepted, true);
            emit ChallengeResolved(intData[0], levelAccepted - 1, false);
        }
        _challengeLevel[intData[0]][idx] = finalLevel;
        intData[1] = finalLevel;
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
