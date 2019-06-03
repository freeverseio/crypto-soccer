pragma solidity ^0.5.0;

import "./Leagues.sol";
import "../state/PlayerState.sol";

contract Engine is PlayerState {
    // @dev Max num of players allowed in a team
    uint8 constant kMaxPlayersInTeam = 11;
    uint256 constant kBitsPerRndNum = 14; 
    uint8 constant rndsPerUint256 = 18; // = 256 / kBitsPerRndNum;
    uint256 constant mask = (1 << kBitsPerRndNum)-1; // (2**bits)-1

    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param state0 a vector with the state of the players of team 0
     * @param state1 a vector with the state of the players of team 1
     * @param tactic0 a vector[3] with the tactic (ex. [4,4,2]) of team 0 
     * @param tactic1 a vector[3] with the tactic (ex. [4,4,2]) of team 1
     * @return the score of the match
     */
    function playMatch(
        bytes32 seed,
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
        bytes32 hash0 = keccak256(abi.encode(uint256(seed) + state0[0] + tactic0[0]));
        bytes32 hash1 = keccak256(abi.encode(uint256(seed) + state1[0] + tactic1[0]));
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

/*
    function decode(uint8 nElem, uint serialized, uint bits) internal pure returns(uint16[] decoded) {
        require (bits <= 16, "Not enough bits to encode each number, since they are read as uint16");
        uint mask = (1 << bits)-1; // (2**bits)-1
        decoded = new uint16[](nElem);
        for (uint8 i=0; i<nElem; i++) {
            decoded[i] = uint16(serialized & mask);
            serialized >>= bits;
        }
    }
*/


/*
    /// @dev Plays a game and, currently, returns the number of goals by each team.
    function playMatchOld(
        bytes32 seed,
        uint256[] memory state0,
        uint256[] memory state1, 
        uint8[3] memory tactic0, 
        uint8[3] memory tactic1
    )
        internal
        returns (uint16[2] memory teamGoals)
    {
        /// @dev We extract 18 randnumbers, each is 14 bit long, from a uint256
        ///  generated from a seed. We do that 4 times. Each of this 4 arrays
        ///  is used in a particular event of the 18 rounds. 
        uint16[] memory rndNum1 = getRndNumArrays(seed, kRoundsPerGame, kBitsPerRndNum);
        uint16[] memory rndNum2 = getRndNumArrays(seed+1, kRoundsPerGame, kBitsPerRndNum);
        uint16[] memory rndNum3 = getRndNumArrays(seed+2, kRoundsPerGame, kBitsPerRndNum);
        uint16[] memory rndNum4 = getRndNumArrays(seed+3, kRoundsPerGame, kBitsPerRndNum);

        uint[5][2] memory globSkills;
        uint[kMaxPlayersInTeam][2] memory attackersSpeed;
        uint[kMaxPlayersInTeam][2] memory attackersShoot;
        uint8[2] memory nAttackers;
        (globSkills[0], nAttackers[0], attackersSpeed[0], attackersShoot[0]) = getGameglobSkills(teamIdx1);
        (globSkills[1], nAttackers[1], attackersSpeed[1], attackersShoot[1]) = getGameglobSkills(teamIdx2);
        uint gameId = getGameId(teamIdx1, teamIdx2, seed);

        uint8 teamThatAttacks;
        /// @dev order of globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance, 5-startEndurance]
        for (uint8 round = 0; round < kRoundsPerGame; round++){
            if ( (round == 8) || (round == 13)) {
                teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice(globSkills[0][kMove2Attack], globSkills[1][kMove2Attack], rndNum1[round], kMaxRndNum);
            emit TeamAttacks(teamThatAttacks, round, gameId);
            if ( managesToShoot(teamThatAttacks, globSkills, rndNum2[round], kMaxRndNum)) {
                if ( managesToScore(
                    nAttackers[teamThatAttacks],
                    attackersSpeed[teamThatAttacks],
                    attackersShoot[teamThatAttacks],
                    globSkills[1-teamThatAttacks][kBlockShoot],
                    rndNum3[round],
                    rndNum4[round],
                    kMaxRndNum,
                    round,
                    gameId
                    )
                ) 
                {
                    teamGoals[teamThatAttacks]++;
                }
            }
        }
        return teamGoals;
    }
*/


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

