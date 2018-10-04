pragma solidity ^ 0.4.24;
import "./teams.sol";

/*
    Main contract with the Game Engine
*/

contract GameEngine is TeamFactory {

    /// @dev Plays a game and, currently, returns the number of goals by each team.
    function playGame(uint teamIdx1, uint teamIdx2, uint[] rndNum1, uint[] rndNum2, uint[] rndNum3, uint[] rndNum4)
        internal
        view
        returns (uint16[2] memory teamGoals)
    {
        /// @dev TODO: use an enum!!
        /// @dev order of globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance, 5-startEndurance]

        uint nRounds = rndNum1.length;
        require (nRounds == rndNum2.length, "We need more randoms for so many round");
        require (nRounds == rndNum3.length, "We need more randoms for so many round");
        require (nRounds == rndNum4.length, "We need more randoms for so many round");

        uint[5][2] memory globSkills;
        uint[kMaxPlayersInTeam][2] memory attackersSpeed;
        uint[kMaxPlayersInTeam][2] memory attackersShoot;
        uint8[2] memory nAttackers;
        (globSkills[0], nAttackers[0], attackersSpeed[0], attackersShoot[0]) = getGameglobSkills(teamIdx1);
        (globSkills[1], nAttackers[1], attackersSpeed[1], attackersShoot[1]) = getGameglobSkills(teamIdx2);

        uint8 teamThatAttacks;
        /// @dev order of globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance, 5-startEndurance]
        for (uint round = 0; round < nRounds; round++){
            if ( (round == 8) || (round == 13)) {
                teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice( globSkills[0][0], globSkills[1][0], rndNum1[round], 1000);
            if ( managesToShoot(teamThatAttacks, globSkills, rndNum2[round], 1000)) {
                if ( managesToScore(
                        nAttackers[teamThatAttacks],
                        attackersSpeed[teamThatAttacks],
                        attackersShoot[teamThatAttacks],
                        globSkills[1-teamThatAttacks][3],
                        rndNum3[round],
                        rndNum4[round],
                        1000
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
        uint currentEnduranceA = skillsTeamA[4];
        uint currentEnduranceB = skillsTeamB[4];
        for (uint8 sk=0; sk<4; sk++) {
            skillsTeamA[sk] = (skillsTeamA[sk] * currentEnduranceA) / 100;
            skillsTeamB[sk] = (skillsTeamB[sk] * currentEnduranceB) / 100;
        }
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
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
        returns (bool)
    {
        /// @dev attacker who actually shoots is selected weighted by his speed
        uint[] memory weights = new uint[](nAttackers);
        for (uint8 p=0; p<nAttackers; p++) {
            weights[p] = attackersSpeed[p];
        }
        uint8 shooter = throwDiceArray(weights, rndNum1, factor);

        /// @dev a goal is scored by confronting his shoot skill to the goalkeeper block skill
        return throwDice((attackersShoot[shooter]*7)/10, blockShoot, rndNum2, factor) == 0;
    }

    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    function managesToShoot(uint8 teamThatAttacks, uint[5][2] globSkills, uint rndNum, uint factor)
        internal
        pure
        returns (bool)
    {
        return throwDice(
            globSkills[1-teamThatAttacks][2],       // defendShoot of defending team against...
            (globSkills[teamThatAttacks][1]*6)/10,  // createShoot of attacking team.
            rndNum,
            factor
        ) == 1 ? true : false;
    }


    /// @dev Computes basic data, including globalSkills, needed during the game.
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
        // move2attack =    defence(defenders + 2*midfields + attackers) +
        //                  speed(defenders + 2*midfields) +
        //                  pass(defenders + 3*midfields)
        // createShoot =    speed(attackers) + pass(attackers)
        // defendShoot =    speed(defenders) + defence(defenders);
        // blockShoot  =    shoot(keeper);
        // skills:  0-age
        //          1-defense
        //          2-speed
        //          3-pass
        //          4-shoot (for a goalkeeper, this is interpreted as ability to block a shoot)
        //          5-endurance
        //          6-role (0=goalkeeper, 1=defence, 2=midfield, 3=attacker, 4=retired)
        uint move2attack;
        uint createShoot;
        uint defendShoot;
        uint blockShoot;
        uint endurance;

        nAttackers = 0;
        for (uint8 p = 0; p < kMaxPlayersInTeam; p++) {
            uint16[] memory skills = decode(7, getSkill(_teamIdx, p), 14);
            endurance += skills[5];
            if (skills[6] == 0) {
                blockShoot = skills[4];
            }
            else if (skills[6] == 1) {
                move2attack = move2attack + skills[1] + skills[2] + skills[3];
                defendShoot = defendShoot + skills[2] + skills[1];
            }
            else if (skills[6] == 2) {
                move2attack = move2attack + 2 * skills[1] + 2 * skills[2] + 3 * skills[3];
            }
            else if (skills[6] == 3) {
                move2attack = move2attack + skills[1];
                createShoot = createShoot + skills[2] + skills[3];
                attackersSpeed[nAttackers] = skills[2];
                attackersShoot[nAttackers] = skills[4];
                nAttackers++;
            }
        }
        // endurance is converted to a percentage that will be maintained:
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
            nAttackers,
            attackersSpeed,
            attackersShoot
        );
    }
}

