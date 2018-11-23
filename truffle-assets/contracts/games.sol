pragma solidity ^ 0.4.24;
import "./factories/teams.sol";

/*
    Main contract with the Game Engine
*/

contract GameEngine is TeamFactory {

    /// @dev gameId is needed to identify the game to which the events belong.
    ///  Currently, the hash of concat(teamIdx1, teamIdx2, seed)
    event TeamAttacks(uint8 homeOrAway, uint8 round, uint gameId);
    event ShootResult(bool isGoal, uint8 attackerIdx, uint8 round, uint gameId);

    /// @dev Plays a game and, currently, returns the number of goals by each team.
    function playGame(uint teamIdx1, uint teamIdx2, uint seed)
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

    /// @dev Rescales global skills of both teams according to their endurance
    function teamsGetTired(uint[5] skillsTeamA, uint[5] skillsTeamB )
        internal
        pure
    {
        /// @dev recall the endurance is a val for which 0 is greatest, 2000 is avg starting
        uint currentEnduranceA = skillsTeamA[kEndurance];
        uint currentEnduranceB = skillsTeamB[kEndurance];
        for (uint8 sk = kMove2Attack; sk < kEndurance; sk++) {
            skillsTeamA[sk] = (skillsTeamA[sk] * currentEnduranceA) / 100;
            skillsTeamB[sk] = (skillsTeamB[sk] * currentEnduranceB) / 100;
        }
    }

    /// @dev Wrapper for the function that actually computes if they manage to score.
    ///  This one emits the event, in case they do.
    function managesToScore(
        uint8 nAttackers,
        uint[kMaxPlayersInTeam] memory attackersSpeed,
        uint[kMaxPlayersInTeam] memory attackersShoot,
        uint blockShoot,
        uint rndNum1,
        uint rndNum2,
        uint factor,
        uint8 round,
        uint gameId
    )
        internal
        returns (bool)
    {
        (bool isGoal, uint8 shooter) = 
            managesToScorePure(
                nAttackers,
                attackersSpeed,
                attackersShoot,
                blockShoot,
                rndNum1,
                rndNum2,
                factor
            );
        emit ShootResult(isGoal, shooter, round, gameId);
        return isGoal; 
    }


    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScorePure(
        uint8 nAttackers,
        uint[kMaxPlayersInTeam] memory attackersSpeed,
        uint[kMaxPlayersInTeam] memory attackersShoot,
        uint blockShoot,
        uint rndNum1,
        uint rndNum2,
        uint factor
    )
        internal
        pure
        returns (bool, uint8)
    {
        /// @dev attacker who actually shoots is selected weighted by his speed
        uint[] memory weights = new uint[](nAttackers);
        for (uint8 p=0; p<nAttackers; p++) {
            weights[p] = attackersSpeed[p];
        }
        uint8 shooter = throwDiceArray(weights, rndNum1, factor);

        /// @dev a goal is scored by confronting his shoot skill to the goalkeeper block skill
        return (
            throwDice((attackersShoot[shooter]*7)/10, blockShoot, rndNum2, factor) == 0,
            shooter
        );
    }

    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    function managesToShoot(uint8 teamThatAttacks, uint[5][2] globSkills, uint rndNum, uint factor)
        internal
        pure
        returns (bool)
    {
        return throwDice(
            globSkills[1-teamThatAttacks][kDefendShoot],       // defendShoot of defending team against...
            (globSkills[teamThatAttacks][kCreateShoot]*6)/10,  // createShoot of attacking team.
            rndNum,
            factor
        ) == 1 ? true : false;
    }


    /// @dev Computes basic data, including globalSkills, needed during the game.
    /// @dev Basically implements the formulas:
    // move2attack =    defence(defenders + 2*midfields + attackers) +
    //                  speed(defenders + 2*midfields) +
    //                  pass(defenders + 3*midfields)
    // createShoot =    speed(attackers) + pass(attackers)
    // defendShoot =    speed(defenders) + defence(defenders);
    // blockShoot  =    shoot(keeper);
    function getGameglobSkills(uint _teamIdx)
        internal
        view
        returns (
            uint[5] globSkills,
            uint8 nAttackers,
            uint[kMaxPlayersInTeam] memory attackersSpeed,
            uint[kMaxPlayersInTeam] memory attackersShoot
        )
    {
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
    }
}

