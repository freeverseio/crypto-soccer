pragma solidity ^0.5.0;

import "./Leagues.sol";
import "../state/PlayerState.sol";

contract Engine is PlayerState {
    // @dev Max num of players allowed in a team
    uint8 constant kMaxPlayersInTeam = 11;
    uint256 constant kBitsPerRndNum = 14; 
    uint8 constant rndsPerUint256 = 18; // = 256 / kBitsPerRndNum;
    uint256 constant mask = (1 << kBitsPerRndNum)-1; // (2**bits)-1
    uint8 constant kRoundsPerGame = 9; // 18; 
    uint16 constant kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 
    /// @dev Ennum for globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance]
    uint8 constant kMove2Attack = 0; 
    uint8 constant kCreateShoot = 1; 
    uint8 constant kDefendShoot = 2; 
    uint8 constant kBlockShoot = 3; 
    uint8 constant kEndurance = 4; 


    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param state0 a vector with the state of the players of team 0
     * @param state1 a vector with the state of the players of team 1
     * @param tactic0 a vector[3] with the tactic (ex. [4,4,2]) of team 0 
     * @param tactic1 a vector[3] with the tactic (ex. [4,4,2]) of team 1
     * @return the score of the match
     */
    function playMatchOld(
        uint256 seed,
        uint256[] memory state0,
        uint256[] memory state1, 
        uint8[3] memory tactic0, 
        uint8[3] memory tactic1
    ) 
        public 
        pure 
        returns (uint8 home, uint8 visitor) 
    {
        require(state0.length == 11, "Team 0 needs 11 players");
        require(state1.length == 11, "Team 1 needs 11 players");
        require(tactic0[0] + tactic0[1] + tactic0[2] == 10, "wrong tactic for team 0");
        require(tactic1[0] + tactic1[1] + tactic1[2] == 10, "wrong tactic for team 1");
        bytes32 hash0 = keccak256(abi.encode(seed + state0[0] + tactic0[0]));
        bytes32 hash1 = keccak256(abi.encode(seed + state1[0] + tactic1[0]));
        return (uint8(uint256(hash0) % 4), uint8(uint256(hash1) % 4));
    }

    function getNRandsFromSeed(uint16 nRands, uint256 seed) public pure returns (uint16[] memory rnds) {
        rnds = new uint16[](nRands);
        uint256 currentBigRnd = uint(keccak256(abi.encodePacked(seed)));
        uint8 rndsFromSameBigRnd = 0;
        for (uint8 n = 0; n < nRands; n++) {
            if (rndsFromSameBigRnd == rndsPerUint256) {
                currentBigRnd = uint(keccak256(abi.encodePacked(seed+1)));
                rndsFromSameBigRnd = 0;
            }
            rnds[n] = uint16(currentBigRnd & mask);
            currentBigRnd >>= kBitsPerRndNum;
            rndsFromSameBigRnd ++;
        }
        return rnds;
    }


    /// @dev Throws a dice that returns 0 with probability weight1/(weight1+weight2), and 1 otherwise.
    /// @dev So, returning 0 has semantics: "the responsible for weight1 is selected".
    /// @dev We return a uint8, not bool, to allow the return to be used as an idx in an array by the callee.
    /// @dev The formula is derived as follows. Throw a random number R in the range [0,maxR].
    /// @dev Then, w1 wins if (w1+w2)*(R/maxR) < w1, and w2 wins otherise. 
    /// @dev maxRndNum controls the resolution or fine-graining of the algorithm.
    function throwDice(uint weight1, uint weight2, uint rndNum, uint maxRndNum)
        public
        pure
        returns(uint8)
    {
        if( ( (weight1 + weight2) * rndNum ) < ( weight1 * (maxRndNum-1) ) ) {
            return 0;
        } else {
            return 1;
        }
    }

    /// @dev Generalization of the previous to any number of input weights
    /// @dev It therefore throws any number of dice and returns the winner's idx.
    function throwDiceArray(uint[] memory weights, uint rndNum, uint maxRndNum)
        public
        pure
        returns(uint8 w)
    {
        uint uniformRndInSumOfWeights;
        for (w = 0; w<weights.length; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;
        uint cumSum = 0;
        for (w = 0; w<weights.length-1; w++) {
            cumSum += weights[w];
            if( uniformRndInSumOfWeights < ( cumSum * (maxRndNum-1) )) {
                return w;
            }
        }
        return w;
    }


    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    // TODO: remove maxRndNum for the constant
    function managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum, uint maxRndNum)
        public
        pure
        returns (bool)
    {
        return throwDice(
            globSkills[1-teamThatAttacks][kDefendShoot],       // defendShoot of defending team against...
            (globSkills[teamThatAttacks][kCreateShoot]*6)/10,  // createShoot of attacking team.
            rndNum,
            maxRndNum
        ) == 1 ? true : false;
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
        uint8 nAttackers,
        uint[] memory attackersSpeed,
        uint[] memory attackersShoot,
        uint blockShoot,
        uint rndNum1,
        uint rndNum2,
        uint maxRndNum
    )
        public
        pure
        returns (bool)
    {
        /// attacker who actually shoots is selected weighted by his speed
        uint[] memory weights = new uint[](nAttackers);
        for (uint8 p = 0; p < nAttackers; p++) {
            weights[p] = attackersSpeed[p];
        }
        uint8 shooter = throwDiceArray(weights, rndNum1, maxRndNum);

        /// a goal is scored by confronting his shoot skill to the goalkeeper block skill
        return throwDice((attackersShoot[shooter]*7)/10, blockShoot, rndNum2, maxRndNum) == 0;
    }

    /// @dev Plays a game and, currently, returns the number of goals by each team.
    function playMatch(
        uint256 seed,
        uint256[] memory state0,
        uint256[] memory state1, 
        uint8[3] memory tactic0, 
        uint8[3] memory tactic1
    )
        public
        pure
        returns (uint8 home, uint8 visitor) 
    {
        require(state0.length == 11, "Team 0 needs 11 players");
        require(state1.length == 11, "Team 1 needs 11 players");
        require(tactic0[0] + tactic0[1] + tactic0[2] == 10, "wrong tactic for team 0");
        require(tactic1[0] + tactic1[1] + tactic1[2] == 10, "wrong tactic for team 1");
        uint16[] memory rnds = getNRandsFromSeed(kRoundsPerGame*4, seed);
        uint[5][2] memory globSkills;
        uint[][2] memory attackersSpeed;
        uint[][2] memory attackersShoot;
        uint8[2] memory nAttackers;
        // TODO: ugly
        nAttackers[0] = tactic0[2];
        nAttackers[1] = tactic1[2];
        (globSkills[0], attackersSpeed[0], attackersShoot[0]) = getTeamGlobSkills(state0, tactic0);
        (globSkills[1], attackersSpeed[1], attackersShoot[1]) = getTeamGlobSkills(state1, tactic1);
        uint8 teamThatAttacks;
        uint8[2] memory teamGoals;

        for (uint8 round = 0; round < kRoundsPerGame; round++){
            // TODO: team gets tired
            teamThatAttacks = throwDice(globSkills[0][kMove2Attack], globSkills[1][kMove2Attack], rnds[4*round], kMaxRndNum);
            if ( managesToShoot(teamThatAttacks, globSkills, rnds[4*round+1], kMaxRndNum)) {
                if ( managesToScore(
                    nAttackers[teamThatAttacks],
                    attackersSpeed[teamThatAttacks],
                    attackersShoot[teamThatAttacks],
                    globSkills[1-teamThatAttacks][kBlockShoot],
                    rnds[4*round+2],
                    rnds[4*round+3],
                    kMaxRndNum
                    )
                ) 
                {
                    teamGoals[teamThatAttacks]++;
                }
            }
        }
        return (teamGoals[0], teamGoals[1]);
    }


    /// @dev Computes basic data, including globalSkills, needed during the game.
    /// @dev Basically implements the formulas:
    // move2attack =    defence(defenders + 2*midfields + attackers) +
    //                  speed(defenders + 2*midfields) +
    //                  pass(defenders + 3*midfields)
    // createShoot =    speed(attackers) + pass(attackers)
    // defendShoot =    speed(defenders) + defence(defenders);
    // blockShoot  =    shoot(keeper);
    function getTeamGlobSkills(uint256[] memory teamState, uint8[3] memory tactic)
        public
        pure
        returns (
            uint[5] memory globSkills,
            uint[] memory attackersSpeed, 
            uint[] memory attackersShoot
        )
    {
        attackersSpeed = new uint[](tactic[2]); 
        attackersShoot = new uint[](tactic[2]); 

        uint move2attack;
        uint createShoot;
        uint defendShoot;
        uint blockShoot;
        uint endurance;

        uint8 p = 0;

        // for a keeper, the 'shoot skill' is interpreted as block skill
        blockShoot  += getShoot(teamState[p]); 
        endurance   += getEndurance(teamState[p]);
        p++;

        // loop over defenders
        for (uint8 i = 0; i < tactic[0]; i++) {
            move2attack += getDefence(teamState[p]) + getSpeed(teamState[p]) + getPass(teamState[p]);
            defendShoot += getDefence(teamState[p]) + getSpeed(teamState[p]);
            endurance   += getEndurance(teamState[p]);
            p++;
        }
        // loop over midfielders
        for (uint8 i = 0; i < tactic[1]; i++) {
            move2attack += 2*getDefence(teamState[p]) + 2*getSpeed(teamState[p]) + 3*getPass(teamState[p]);
            endurance   += getEndurance(teamState[p]);
            p++;
        }
        // loop over strikers
        for (uint8 i = 0; i < tactic[2]; i++) {
            move2attack += getDefence(teamState[p]) ;
            createShoot += getSpeed(teamState[p]) + getPass(teamState[p]);
            endurance   += getEndurance(teamState[p]);
            attackersSpeed[i] = getSpeed(teamState[p]); 
            attackersShoot[i] = getShoot(teamState[p]); 
            p++;
        }

        // endurance is converted to a percentage, 
        // used to multiply (and hence decrease) the start endurance.
        // 100 is super-endurant (1500), 70 is bad, for an avg starting team (550).
        if (endurance < 500) {
            endurance = 70;
        } else if (endurance < 1400) {
            endurance = 100 - (1400-endurance)/30;
        } else {
            endurance = 100;
        }


        return (
            [move2attack, createShoot, defendShoot, blockShoot, endurance],
            attackersSpeed,
            attackersShoot
        );

/*
        uint move2attack;
        uint createShoot;
        uint defendShoot;
        uint blockShoot;
        uint endurance;

        nAttackers = 0;

        for (uint8 p = 0; p < kMaxPlayersInTeam; p++) {
            uint16[] memory skills = decode(kNumStates, getStatePlayerInTeam(p, _teamIdx), kBitsPerState);
            endurance += skills[kStatEndur];
            if (skills[kStatRole] == kRoleKeeper) {
                blockShoot = skills[kStatShoot];
            }
            else if (skills[kStatRole] == kRoleDef) {
                move2attack = move2attack + skills[kStatDef] + skills[kStatSpeed] + skills[kStatPass];
                defendShoot = defendShoot + skills[kStatSpeed] + skills[kStatDef];
            }
            else if (skills[kStatRole] == kRoleMid) {
                move2attack = move2attack + 2 * skills[kStatDef] + 2 * skills[kStatSpeed] + 3 * skills[kStatPass];
            }
            else if (skills[kStatRole] == kRoleAtt) {
                move2attack = move2attack + skills[kStatDef];
                createShoot = createShoot + skills[kStatSpeed] + skills[kStatPass];
                attackersSpeed[nAttackers] = skills[kStatSpeed];
                attackersShoot[nAttackers] = skills[kStatShoot];
                nAttackers++;
            }
        }
        /// @dev endurance is converted to a percentage, 
        /// @dev used to multiply (and hence decrease) the start endurance.
        /// @dev 100 is super-endurant (1500), 70 is bad, for an avg starting team (550).
        if (endurance < 500) {
            endurance = 70;
        } else if (endurance < 1400) {
            endurance = 100 - (1400-endurance)/30;
        } else {
            endurance = 100;
        }

        return (
            [move2attack, createShoot, defendShoot, blockShoot, endurance],
            nAttackers,
            attackersSpeed,
            attackersShoot
        );
*/

    }
}

