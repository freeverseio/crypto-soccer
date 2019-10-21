pragma solidity ^0.5.0;

import "./EncodingSkills.sol";
import "./EngineLib.sol";
import "./Sort.sol";
import "./EncodingMatchLogPart1.sol";

contract EnginePreComp is EngineLib, EncodingMatchLogPart1, Sort {
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint256 private constant ONE256              = uint256(1); 
    uint8 private constant CHG_HAPPENED        = uint8(1); 
    uint256 private constant CHG_CANCELLED       = uint256(2); 
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    uint256 private constant TEN_TO_4       = uint256(10000); 
    uint256 private constant TEN_TO_10      = uint256(10000000000); 
    uint256 private constant TEN_TO_14      = uint256(100000000000000); 
    uint256 private constant SECS_IN_DAY    = 86400; // 24 * 3600 

    uint8 public constant ROUNDS_PER_MATCH  = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)
    uint8 public constant NO_SUBST  = 11;   // noone was subtituted
    uint8 public constant NO_CARD  = 14;   // noone saw a card
    uint8 public constant RED_CARD  = 3;   // type of event = redCard


    // Over a game, we would like:
    //      - injuryHard = 1 per 100 games => 0.01 per game per player => 0.02 per game
    //      - injuryLow = 0.7 per 100 games => 0.007 per game per player => 0.04 per game
    //      - redCard 1/10 = 0.1 per game
    //      - yellowCard 2.5 per game 
    // We encode this in uint16[3] events, which applies to 1 half of the game only.
    //  - 1 possible event that leaves a player out of the match, encoded in:
    //          events[0, 1] = [player (from 0 to 13), eventType (injuryHard, injuryLow, redCard)]
    //  - 2 possible events for yellow card:
    //          events[2] = player (from 0 to 13)
    //          events[3] = player (from 0 to 13)
    // The player value is set to NO_EVENT ( = 14) if no event took place
    // If we're on the 2nd half, the idx are events[4,5,6,7]
    // for out of game:
    //      it cannot return 0
    //      injuryHard:  1
    //      injuryLow:   2
    //      redCard:     3  (aka RED_CARD)
    function computeExceptionalEvents
    (
        uint256 matchLog,
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        uint8[3] memory substitutions,
        uint8[3] memory subsRounds,
        bool is2ndHalf,
        uint256 seed
    ) 
        public 
        pure 
        returns 
    (
        uint256
    ) 
    {
        uint8 offset = is2ndHalf ? 171 : 151;
        uint256[] memory weights = new uint256[](15);
        uint64[] memory rnds = getNRandsFromSeed(seed + 42, 4);

        // Start by logging that all substitutions are possible. It will be re-written 
        // by the red-card computation below, if needed.
        for (uint8 p = 0; p < 3; p++) {
            if (substitutions[p] != NO_SUBST) {
                matchLog = setInGameSubsHappened(matchLog, CHG_HAPPENED, p, is2ndHalf);
            } 
        }

        // Build weights for players, based on their aggressiveness.
        // Legit players have > 0, but those already out in first half (teams with 10 players) have 0.
        for (uint8 p = 0; p < NO_CARD; p++) {
            if (states[p] != 0) weights[p] = 1 + getAggressiveness(states[p]); // weights must be > 0 to ever be selected
        }
        
        // next: two events for yellow cards
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 2.5 per game = 1.25 per half => 0.75 per dice thrown
        // weight nothing happens = 9
        weights[NO_CARD] = 9;
        uint8[2] memory yellowCardeds;
        
        yellowCardeds[0] = throwDiceArray(weights, rnds[2]);
        yellowCardeds[1] = throwDiceArray(weights, rnds[3]);
    
        // Important: there is no need to check if a red card had been given in the 1st half, because that player, simply
        // would never be able to have been lined-up, and never made it up to here.
    
        // However, it is important to check if a player who saw a yellow card is not linedup anymore
    
        // if both yellows are for the same player => force red card, and leave. Enough punishment.
        if (yellowCardeds[0] == yellowCardeds[1]) {
            return logOutOfGame(is2ndHalf, true, yellowCardeds[0], matchLog, substitutions, subsRounds, rnds[0], rnds[1]);
        }
        // if we are in the 2nd half (and the 2 yellows are for different players):
        // - if any such player had alread received one in the 1st half => force red, record the other yellow card, leave. 
        if (is2ndHalf) {
            if (hadReceivedYellowIn1stHalf(matchLog, yellowCardeds[0])) {
                matchLog = logOutOfGame(is2ndHalf, true, yellowCardeds[0], matchLog, substitutions, subsRounds, rnds[0], rnds[1]);
                yellowCardeds[0] = NO_CARD;
            }
            if (hadReceivedYellowIn1stHalf(matchLog, yellowCardeds[1])) {
                if (!didOutOfGameHappenInThisHalf(matchLog, offset)) {
                    matchLog = logOutOfGame(is2ndHalf, true, yellowCardeds[1],  matchLog, substitutions, subsRounds, rnds[0], rnds[1]);
                }
                yellowCardeds[1] = NO_CARD;
            }
        }
        
        // if we get here: both yellows are to different players, who can continue playing. Record them.
        matchLog = addYellowCard(matchLog, yellowCardeds[0], 0, is2ndHalf);
        matchLog = addYellowCard(matchLog, yellowCardeds[1], 1, is2ndHalf);

        // Redcards & Injuries:
        // if a new red card is given to a previously yellow-carded player, no prob, such things happen.
        // events[0] => STUFF THAT REMOVES A PLAYER FROM FIELD: injuries and redCard 
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 0.07 per game = 0.035 per half => weight nothing happens = 758
        if (!didOutOfGameHappenInThisHalf(matchLog, offset)) {
            weights[NO_CARD] = 758;
            uint256 selectedPlayer = uint256(throwDiceArray(weights, rnds[0]));
            matchLog = logOutOfGame(is2ndHalf, false, selectedPlayer, matchLog, substitutions, subsRounds, rnds[0], rnds[1]);
        }
        // If 1st half, log the yellowCarded guys how managed to end linedup
        if (!is2ndHalf) {
            if (!didPlayerFinish1stHalf(matchLog, yellowCardeds[0], substitutions)) matchLog |= (ONE256 << 169);
            if (!didPlayerFinish1stHalf(matchLog, yellowCardeds[1], substitutions)) matchLog |= (ONE256 << 170);
        }
        return matchLog;
    }
    
    function didPlayerFinish1stHalf(uint256 matchLog, uint256 player, uint8[3] memory substitutions) private pure returns(bool) {
        // check if it was outOfGamed in 1st half: ((matchLog >> 151) & 15) = redCardeds in 1st Half
        // ...note: no need to check type of outOfGame, he cannot be linedup in 2nd half
        if (getOutOfGamePlayer(matchLog, false) == player) return false; 
        // check if it was substituted:
        // ...note: if substitution did not happen because he was redCarded, he'd have falled in previous check.
        for (uint p = 0; p < 3; p++) {
            if (player == substitutions[p]) return false;
        }
        return true;
    }
    
    function hadReceivedYellowIn1stHalf(uint256 matchLog, uint256 newYellowCarded) private pure returns(bool) {
        return   
            // ((matchLog >> 161) & 15) = 1st half yellow card [0]
            // ((log >> 169) & 1) = yellowCardedCouldNotFinish1stHalf[0]
            (newYellowCarded == ((matchLog >> 161) & 15)) && (((matchLog >> 169) & 1) == 0)  ||
            // ((matchLog >> 165) & 15) = 1st half yellow card [1]
            // ((log >> 170) & 1) = yellowCardedCouldNotFinish1stHalf[1]
            (newYellowCarded == ((matchLog >> 165) & 15)) && (((matchLog >> 170) & 1) == 0);
    }
    
    function didOutOfGameHappenInThisHalf(uint256 matchLog, uint8 offset) private pure returns (bool) {
        return ((matchLog >> offset + 8) & 3) != 0;
    }

    function logOutOfGame(
        bool is2ndHalf,
        bool forceRedCard,
        uint256 selectedPlayer, 
        uint256 matchLog,
        uint8[3] memory substitutions,
        uint8[3] memory subsRounds,
        uint64 rnd0,
        uint64 rnd1
    ) private pure returns(uint256) 
    {
        if (selectedPlayer == NO_CARD) return matchLog;
        uint8 offset = is2ndHalf ? 171 : 151;
        matchLog |= selectedPlayer << offset;
        uint8 minRound = 0;
        uint8 maxRound = ROUNDS_PER_MATCH;

        // first compute the type of event        
        uint256 typeOfEvent = forceRedCard ? uint256(RED_CARD) : uint256(computeTypeOfEvent(rnd1));
        matchLog |= typeOfEvent << (offset + 8);

        // if the selected player was one of the guys joining during this half (outGame = 11, 12, or 13),
        // make sure that the round selected for this event is after joining. 
        if (selectedPlayer < 14 && selectedPlayer > 10) {
            minRound = subsRounds[selectedPlayer - 11];
        }
        // if the selected player was one of the guys to be changed during this half (outGame = 0,...10),
        // make sure that the round selected for this event is before the change.
        // (note that substitution[p] == 11 => NO_SUBS, cannot happen 
        // in the next else-if (since selectedPlayer <= 10 in that branch)
        else {
            for (uint8 p = 0; p < 3; p++) {
                if (selectedPlayer == substitutions[p]) {
                    maxRound = subsRounds[p];
                    // log that this substitution was unable to take place
                    if (typeOfEvent == RED_CARD) {
                        matchLog = setInGameSubs(
                            matchLog,
                            is2ndHalf ? 195 + 2 * p : 189 + 2 * p
                        );
                    }
                } 
            }
        }
        matchLog |= uint256(computeRound(rnd0+1, minRound, maxRound)) << offset + 4;
        return matchLog;
    }

    function setInGameSubs(uint256 matchLog, uint8 pos) private pure returns (uint256) {
        return (matchLog & ~(uint256(3) << pos)) | (CHG_CANCELLED << pos);
    }

    
    function computeRound(uint256 seed, uint8 minRound, uint8 maxRound) private pure returns (uint8 round) {
        require(maxRound > minRound, "max and min rounds are not correct");
        return minRound + uint8(seed % (maxRound - minRound + 1));
    }

    // it cannot return 0.
    // injuryHard:  1
    // injuryLow:   2
    // redCard:     3
    function computeTypeOfEvent(uint256 rnd) private pure returns (uint8) {
        uint256[] memory weights = new uint256[](3);
        weights[0] = 1; // injuryHard   
        weights[1] = 2; // injuryLow
        weights[2] = 5; // redCard
        return 1 + throwDiceArray(weights, rnd);
    }

    // TODO: avoid redCarded or outOfGame players to shoot, include changed players.
    function computePenalties(
        uint256[2] memory matchLog, 
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states, 
        uint256 block0, 
        uint256 block1, 
        uint64 seed
    )
        public 
        pure 
        returns(uint256[2] memory) 
    {
        uint64[] memory rnds = getNRandsFromSeed(seed * 7, 14);
        uint8[2] memory totalGoals;
        for (uint8 round = 0; round < 6; round++) {
            if (throwDice(block1, 3 * getShoot(states[0][10-round]), rnds[2 *round]) == 1) {
                matchLog[0] |= (ONE256 << 144 + round);
                totalGoals[0] += 1;
            }
            if (throwDice(block0, 3 * getShoot(states[1][10-round]), rnds[2 *round + 1]) == 1) {
                matchLog[1] |= (ONE256 << 144 + round);
                totalGoals[1] += 1;
            }
            if ((round > 3) && (totalGoals[0] != totalGoals[1])) return matchLog;
        }
        if (throwDice(block0 + getShoot(states[0][4]), block0 + getShoot(states[0][4]), rnds[13]) == 1) {
            matchLog[0] |= (ONE256 << 144 + 6);
        } else {
            matchLog[1] |= (ONE256 << 144 + 6);
        }
        return matchLog;
    }

    /// @dev Computes basic data, including globalSkills, needed during the game.
    /// @dev Basically implements the formulas:
    // move2attack =    defence(defenders + 2*midfields + attackers) +
    //                  speed(defenders + 2*midfields) +
    //                  pass(defenders + 3*midfields)
    // globSkills[IDX_CREATE_SHOOT] =    speed(attackers) + pass(attackers)
    // globSkills[IDX_DEFEND_SHOOT] =    speed(defenders) + defence(defenders);
    // blockShoot  =    shoot(keeper);
    function getTeamGlobSkills(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint8[9] memory playersPerZone, 
        bool[10] memory extraAttack,
        uint256 matchStartTime
    )
        public
        pure
        returns 
    (
        uint256[5] memory globSkills
    )
    {
        // for a keeper, the 'shoot skill' is interpreted as block skill
        // if for whatever reason, user places a non-GK as GK, the block skill is a terrible minimum.
        uint256 penalty;
        uint256 playerSkills = states[0];
        globSkills[IDX_ENDURANCE] = getEndurance(playerSkills);
        if (computePenaltyBadPositionAndCondition(0, playersPerZone, playerSkills) == 0) {
            globSkills[IDX_BLOCK_SHOOT] = (10 * penaltyPerAge(playerSkills, matchStartTime))/1000000;
        }
        else globSkills[IDX_BLOCK_SHOOT] = (getShoot(playerSkills) * penaltyPerAge(playerSkills, matchStartTime))/1000000;
        
        uint256[3] memory fwdModFactors;

        for (uint8 p = 1; p < 11; p++){
            playerSkills = states[p];
            penalty = computePenaltyBadPositionAndCondition(p, playersPerZone, playerSkills) * penaltyPerAge(playerSkills, matchStartTime);
            fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
            if (p < 1 + getNDefenders(playersPerZone)) {computeDefenderGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}
            else if (p < 1 + getNDefenders(playersPerZone) + getNMidfielders(playersPerZone)) {computeMidfielderGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}
            else {computeForwardsGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}       
        }
        // endurance is converted to a percentage, 
        // used to multiply (and hence decrease) the start endurance.
        // 100 is super-endurant (1500), 70 is bad, for an avg starting team (550).
        if (globSkills[IDX_ENDURANCE] < 500) {
            globSkills[IDX_ENDURANCE] = 70;
        } else if (globSkills[IDX_ENDURANCE] < 1400) {
            globSkills[IDX_ENDURANCE] = 100 - (1400-globSkills[IDX_ENDURANCE])/30;
        } else {
            globSkills[IDX_ENDURANCE] = 100;
        }
    }

    // 0 penalty means no penalty
    // 1000 penalty means 10% penalty
    // etc... up to MAX_PENALTY
    function computePenaltyBadPositionAndCondition(
        uint8 lineupPos, 
        uint8[9] memory playersPerZone, 
        uint256 playerSkills
    ) 
        public
        pure
        returns (uint256 penalty) 
    {
        require(lineupPos < 11, "wrong arg in computePenaltyBadPositionAndCondition");
        uint256 forwardness = getForwardness(playerSkills);
        uint256 leftishness = getLeftishness(playerSkills);
        if (forwardness == IDX_GK && lineupPos > 0 || forwardness != IDX_GK && lineupPos == 0) return 0;
        uint8[9] memory playersBelow = playersBelowZones(playersPerZone);
        lineupPos--; // remove the offset due to the GK
        if (lineupPos < playersBelow[0]) { 
            // assigned to defense left
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[1]) { 
            // assigned to defense center
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else if (lineupPos < playersBelow[2]) { 
            // assigned to defense left
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForRights(leftishness);
        } else if (lineupPos < playersBelow[3]) { 
            // assigned to mid left
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[4]) { 
            // assigned to mid center
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else if (lineupPos < playersBelow[5]) { 
            // assigned to mid right
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForRights(leftishness);
        } else if (lineupPos < playersBelow[6]) { 
            // assigned to attack left
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[7]) { 
            // assigned to attack center
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else { 
            // assigned to attack right
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForRights(leftishness);
        }
        uint256 gamesNonStop = getGamesNonStopping(playerSkills);
        if (gamesNonStop > 5) {
            return 5000 - penalty;
        } else {
            return 10000 - gamesNonStop * 1000 - penalty;
        }
    }

    // recall order: [MOVE2ATTACK, CREATE_SHOOT, DEFEND_SHOOT, BLOCK_SHOOT, ENDURANCE]
    // the forward modifier factors only change the first 3.
    function getExtraAttackFactors(bool extraAttack) public pure returns (uint256[3] memory fwdModFactors) {
        if (extraAttack)    {fwdModFactors = [uint256(10500), uint256(10500), uint256(9500)];}
        else                {fwdModFactors = [TEN_TO_4, TEN_TO_4, TEN_TO_4];}
    }
  
    function computeDefenderGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills) + getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_DEFEND_SHOOT] += ((getDefence(playerSkills) + getSpeed(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 30;
            globSkills[IDX_DEFEND_SHOOT] += 20;
            globSkills[IDX_ENDURANCE]   += 10;
        }
    }


    function computeMidfielderGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((2*getDefence(playerSkills) + 2*getSpeed(playerSkills) + 3*getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 50;
            globSkills[IDX_ENDURANCE]   += 10;
        }
    }
    
    
    function computeForwardsGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_CREATE_SHOOT] += ((getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE] += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 10;
            globSkills[IDX_CREATE_SHOOT] += 20;
            globSkills[IDX_ENDURANCE] += 10;
        }
    }

    function playersBelowZones(uint8[9] memory playersPerZone) private pure returns(uint8[9] memory  playersBelow) {
        playersBelow[0] = playersPerZone[0];    
        for (uint8 z = 1; z < 9; z++) {
            playersBelow[z] = playersBelow[z-1] + playersPerZone[z];
        }
    }

    function penaltyForLefts(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_C || leftishness == IDX_CR) {return 1000;} 
            else if (leftishness == IDX_R) {return 2000;}
    }

    function penaltyForCenters(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_L || leftishness == IDX_R) {return 1000;} 
    }

    function penaltyForRights(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_C || leftishness == IDX_LC) {return 1000;} 
            else if (leftishness == IDX_L) {return 2000;}
    }
    
    function penaltyForDefenders(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_M || forwardness == IDX_MF) {return 1000;}
            else if (forwardness == IDX_F) {return 2000;}
    }

    function penaltyForMids(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_D || forwardness == IDX_F) {return 1000;}
    }

    function penaltyForAttackers(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_M || forwardness == IDX_MD) {return 1000;}
            else if (forwardness == IDX_D) {return 2000;}
    }

    function getLinedUpStates(
        uint256 matchLog, 
        uint256 tactics, 
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        bool is2ndHalf
    )
        public 
        pure 
        returns 
    (
        uint256[PLAYERS_PER_TEAM_MAX] memory outStates,
        uint8 tacticsId,
        bool[10] memory extraAttack,
        uint8[3] memory substitutions,
        uint8[3] memory subsRounds
    ) 
    {
        uint8[14] memory lineup;
        (substitutions, subsRounds, lineup, extraAttack, tacticsId) = decodeTactics(tactics);
        uint8 changes;
        uint8 emptyShirts; 
        
        // Count changes during half-time, as well as not-aligned players
        // ...note: substitutions = 11 means NO_SUBS
        for (uint8 p = 0; p < 11; p++) {
            outStates[p] = states[lineup[p]];
            if (outStates[p] == 0) {
                emptyShirts++;
            } else if (is2ndHalf && !getAlignedEndOfLastHalf(outStates[p])) {
                matchLog |= (uint256(p) << 201 + 4 * changes);
                changes++; 
            }
        }

        // if is2ndHalf: make sure we align 10 or 11 players depedning on possible 1st half redcards
        if (is2ndHalf && wasThereRedCardIn1stHalf(matchLog)) {
            require(emptyShirts == 1, "You cannot line up 11 players if there was a red card in 1st half");
        } else {
            require(emptyShirts == 0, "You must line up 11 players");
        }
        
        // Count changes ingame during 1st half
        // matchLog >> 189, 190, 191 contain ingameSubsCancelled
        if (is2ndHalf) {
            for (uint8 p = 0; p < 3; p++) {
                if(((matchLog >> 189 + 2*p) & 3) == CHG_HAPPENED) changes++;
            }        
        }

        if (substitutions[0] < 11) {
            changes++;
            outStates[11] = states[lineup[11]];
            assertCanPlay(outStates[11]);
            require(!getAlignedEndOfLastHalf(outStates[11]), "cannot align a player who already left the field once");
        }
        if (substitutions[1] < 11) { 
            changes++;
            require(substitutions[0] != substitutions[1], "changelist incorrect");
            outStates[12] = states[lineup[12]];
            assertCanPlay(outStates[12]);
            require(!getAlignedEndOfLastHalf(outStates[11]), "cannot align a player who already left the field once");
        }
        if (substitutions[2] < 11) {
            changes++;
            require((substitutions[0] != substitutions[2]) && (substitutions[1] != substitutions[2]), "changelist incorrect");
            outStates[13] = states[lineup[13]];
            assertCanPlay(outStates[13]);
            require(!getAlignedEndOfLastHalf(outStates[11]), "cannot align a player who already left the field once");
        }
        require(changes < 4, "max allowed changes in a game is 3");
        lineup = sort14(lineup);
        for (uint8 p = 1; p < 11; p++) require(lineup[p] > lineup[p-1], "player appears twice in lineup!");        
    }

    function assertCanPlay(uint256 playerSkills) public pure {
        require(getPlayerIdFromSkills(playerSkills) != FREE_PLAYER_ID, "free player shirt has been aligned");
        require(!getRedCardLastGame(playerSkills) && getInjuryWeeksLeft(playerSkills) == 0, "player injured or sanctioned");
        require(!getSubstitutedLastHalf(playerSkills), "cannot align player who was already substituted");
    }

    function wasThereRedCardIn1stHalf(uint256 matchLog) private pure returns(bool) {
        return ((matchLog >> 159) & 3) == RED_CARD;
    }

}

