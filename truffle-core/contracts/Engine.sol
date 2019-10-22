pragma solidity ^0.5.0;

import "./EnginePreComp.sol";
import "./EngineLib.sol";
import "./EncodingMatchLogPart3.sol";

contract Engine is EngineLib, EncodingMatchLogPart3 {
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint8 public constant ROUNDS_PER_MATCH  = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    //
    uint8 private constant IDX_IS_2ND_HALF      = 0; 
    uint8 private constant IDX_IS_HOME_STADIUM  = 1; 
    uint8 private constant IDX_IS_PLAYOFF       = 2; 
    //
    uint8 private constant IDX_SEED         = 0; 
    uint8 private constant IDX_ST_TIME      = 1; 
    // 
    uint256 private constant CHG_HAPPENED   = uint256(1); 
    uint8 public constant RED_CARD  = 3;   // type of event = redCard
    uint8 private constant WINNER_AWAY = 1;
    uint8 private constant WINNER_DRAW = 2;

    bool dummyBoolToEstimateCost;

    EnginePreComp private _precomp;

    function setCardsAndInjuries(address addr) public {
        _precomp = EnginePreComp(addr);
    }

    // mock up to estimate cost of a match.
    // to be removed before deployment
    function playMatchWithCost(
        uint256 seed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public
    {
        playMatch(seed, matchStartTime, states, tactics, matchLog, matchBools);
        dummyBoolToEstimateCost = !dummyBoolToEstimateCost; 
    }


    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param states a 2-vector, each of the 2 being vector with the state of the players of team 0
     * @param tactics a 2-vector with the tacticId (ex. 0 for [4,4,2]) for each team
     * @return the score of the match
     */
    function playMatch(
        uint256 seed,
        uint256 matchStartTime, //actionsSubmissionTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public
        view
        returns (uint256[2] memory)
    {
        uint256 block0;
        uint256 block1;
        (matchLog, block0, block1) = playMatchWithoutPenalties(
            [seed, matchStartTime], 
            states,
            tactics,
            matchLog,
            matchBools
        );
        if (matchBools[IDX_IS_PLAYOFF] && ( getNGoals(matchLog[0]) == getNGoals(matchLog[1]))) {
            matchLog = _precomp.computePenalties(matchLog, states, block0, block1, uint64(seed));  // TODO seed
        } else {
            // note that WINNER_HOME = 0, so no need to write anything if home wins.
            if (getNGoals(matchLog[0]) == getNGoals(matchLog[1])) addWinnerToBothLogs(matchLog, WINNER_DRAW);
            else if (getNGoals(matchLog[0]) < getNGoals(matchLog[1])) addWinnerToBothLogs(matchLog, WINNER_AWAY);
        }
        return matchLog;
    }
    
    /**
     * @dev playMatch returns the result of a match
     * @param states a 2-vector, each of the 2 being vector with the state of the players of team 0
     * @param tactics a 2-vector with the tacticId (ex. 0 for [4,4,2]) for each team
     * @return the score of the match
     */
    function playMatchWithoutPenalties(
        uint256[2] memory seedAndStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        private
        view
        returns (uint256[2] memory, uint256, uint256)
    {
        uint256[5][2] memory globSkills;
        uint8[9][2] memory playersPerZone;
        bool[10][2] memory extraAttack;

        (matchLog[0], states[0], playersPerZone[0]) = getLineUpAndPlayerPerZone(states[0], tactics[0], matchBools[IDX_IS_2ND_HALF], matchLog[0], seedAndStartTime[IDX_SEED]);
        (matchLog[1], states[1], playersPerZone[1]) = getLineUpAndPlayerPerZone(states[1], tactics[1], matchBools[IDX_IS_2ND_HALF], matchLog[1], seedAndStartTime[IDX_SEED]);

        matchLog[0] = writeNDefs(matchLog[0], states[0], getNDefenders(playersPerZone[0]), matchBools[IDX_IS_2ND_HALF]);
        matchLog[1] = writeNDefs(matchLog[1], states[1], getNDefenders(playersPerZone[1]), matchBools[IDX_IS_2ND_HALF]);


        globSkills[0] = _precomp.getTeamGlobSkills(states[0], playersPerZone[0], extraAttack[0], seedAndStartTime[IDX_ST_TIME]);
        globSkills[1] = _precomp.getTeamGlobSkills(states[1], playersPerZone[1], extraAttack[1], seedAndStartTime[IDX_ST_TIME]);

        if (matchBools[IDX_IS_HOME_STADIUM]) {
            globSkills[0][IDX_ENDURANCE] = (globSkills[0][IDX_ENDURANCE] * 11500)/10000;
        }
        computeRounds(matchLog, seedAndStartTime[IDX_SEED], seedAndStartTime[IDX_ST_TIME], states, playersPerZone, extraAttack, globSkills, matchBools[IDX_IS_2ND_HALF]);
        return (matchLog, globSkills[0][IDX_BLOCK_SHOOT], globSkills[1][IDX_BLOCK_SHOOT]);
    }
    
    function computeRounds(
        uint256[2] memory matchLog,
        uint256 seed, 
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states, 
        uint8[9][2] memory playersPerZone, 
        bool[10][2] memory extraAttack, 
        uint256[5][2] memory globSkills, 
        bool is2ndHalf
    ) 
        private
        pure
    {
        uint64[] memory rnds = getNRandsFromSeed(seed, ROUNDS_PER_MATCH*5);
        uint8 teamThatAttacks;
        for (uint8 round = 0; round < ROUNDS_PER_MATCH; round++){
            if (is2ndHalf && ((round == 0) || (round == 5))) {
                teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice(globSkills[0][IDX_MOVE2ATTACK], globSkills[1][IDX_MOVE2ATTACK], rnds[5*round]);
            if ( managesToShoot(teamThatAttacks, globSkills, rnds[5*round+1])) {
                managesToScore(
                    matchStartTime,
                    matchLog,
                    teamThatAttacks,
                    states[teamThatAttacks],
                    playersPerZone[teamThatAttacks],
                    extraAttack[teamThatAttacks],
                    globSkills[1-teamThatAttacks][IDX_BLOCK_SHOOT],
                    [rnds[5*round+2], rnds[5*round+3], rnds[5*round+4]]
                );
            }
        }
    }
    
    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    // this function returns an array of (still) size 25, but only with the correctly linedup first 11 entries filled.
    function getLineUpAndPlayerPerZone(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint256 tactics,
        bool is2ndHalf,
        uint256 matchLog,
        uint256 seed
    ) 
        public 
        view 
        returns (uint256, uint256[PLAYERS_PER_TEAM_MAX] memory outStates, uint8[9] memory) 
    {
        uint8 tacticsId;
        (matchLog, outStates, tacticsId) = _precomp.getLinedUpStates(matchLog, tactics, states, is2ndHalf);
        matchLog = _precomp.computeExceptionalEvents(matchLog, outStates, tactics, is2ndHalf, seed); 
        return (matchLog, outStates, getPlayersPerZone(tacticsId));
    }

    function writeNDefs(
        uint256 matchLog, 
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint8 nDefsInTactics, 
        bool is2ndHalf
    ) private pure returns (uint256) {
        if (is2ndHalf) {
            for (uint8 p = 1; p < 1 + nDefsInTactics; p++) {
                if (states[p] == 0) return addNDefs(matchLog, nDefsInTactics - 1, true);
            }
        }
        return addNDefs(matchLog, nDefsInTactics, is2ndHalf);
    }

    /// @dev Rescales global skills of both teams according to their endurance
    function teamsGetTired(uint256[5] memory skillsTeamA, uint256[5]  memory skillsTeamB )
        public
        pure
         returns (uint256[5] memory , uint256[5] memory ) 
    {
        uint256 currentEnduranceA = skillsTeamA[IDX_ENDURANCE];
        uint256 currentEnduranceB = skillsTeamB[IDX_ENDURANCE];
        for (uint8 sk = IDX_MOVE2ATTACK; sk < IDX_ENDURANCE; sk++) {
            skillsTeamA[sk] = (skillsTeamA[sk] * currentEnduranceA) / 100;
            skillsTeamB[sk] = (skillsTeamB[sk] * currentEnduranceB) / 100;
        }
        return (skillsTeamA, skillsTeamB);
    }

    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    function managesToShoot(uint8 teamThatAttacks, uint256[5][2] memory globSkills, uint256 rndNum)
        public
        pure
        returns (bool)
    {
        return throwDice(
            globSkills[1-teamThatAttacks][IDX_DEFEND_SHOOT],       // globSkills[IDX_DEFEND_SHOOT] of defending team against...
            (globSkills[teamThatAttacks][IDX_CREATE_SHOOT]*6)/10,  // globSkills[IDX_CREATE_SHOOT] of attacking team.
            rndNum
        ) == 1 ? true : false;
    }

    function selectAssister(
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint8 shooter,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[] memory weights = new uint256[](11);
        // if selected assister == selected shooter =>  
        //  there was no assist => individual play by shoorter
        weights[0] = penaltyPerAge(states[0], matchStartTime);
        uint256 teamPassCapacity = weights[0];
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 90 : 20 ) * getPass(states[p]) * penaltyPerAge(states[p], matchStartTime);
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 150 : 100 ) * getPass(states[p]) * penaltyPerAge(states[p], matchStartTime);
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 200 * getPass(states[p]) * penaltyPerAge(states[p], matchStartTime);
            teamPassCapacity += weights[p];
            p++;
        }
        // on average: teamPassCapacity442 = (1 + 4 * 20 + 4 * 100 + 2 * 200) < getPass > = 881 <pass>_team
        // on average: shooterSumOfSkills = 5 * <skills>_shooter
        // so a good ratio is shooterSumOfSkills/teamPassCapacity442 = 5/881 * <skills_shooter>/<pass>_team
        // or better, to have an avg of 1: (shooterSumOfSkills*271)/(teamPassCapacity * 5) = <skills_shooter>/<pass>_team
        // or to have a 50% change, multiply by 10, and to have say, 1/3, multiply by 10/3
        // this is to be compensated by an overall factor of about.
        weights[shooter] = (weights[shooter] * getSumOfSkills(states[shooter]) * 8810 * penaltyPerAge(states[shooter], matchStartTime))/ (N_SKILLS * (teamPassCapacity - weights[shooter]) * 3);
        return throwDiceArray(weights, rnd);
    }


    function selectShooter(
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[] memory weights = new uint256[](11);
        // GK has minimum weight, all others are relative to this.
        weights[0] = penaltyPerAge(states[0], matchStartTime);
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 15000 : 5000 ) * getSpeed(states[p]) * penaltyPerAge(states[p], matchStartTime);
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 50000 : 25000 ) * getSpeed(states[p]) * penaltyPerAge(states[p], matchStartTime);
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 75000 * getSpeed(states[p]) * penaltyPerAge(states[p], matchStartTime);
            p++;
        }
        return throwDiceArray(weights, rnd);
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
        uint256 matchStartTime,
        uint256[2] memory matchLog,
        uint8 teamThatAttacks,
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint256 blockShoot,
        uint64[3] memory rnds
    )
        public
        pure
        returns (uint256[2] memory)
    {
        uint8 currentGoals = getNGoals(matchLog[teamThatAttacks]);
        if (currentGoals > 13) return matchLog;
        uint8 shooter = selectShooter(matchStartTime, states, playersPerZone, extraAttack, rnds[0]);
        /// a goal is scored by confronting his shoot skill to the goalkeeper block skill
        uint256 shootPenalty = ( getForwardness(states[shooter]) == IDX_GK ? 10 : 1) * penaltyPerAge(states[shooter], matchStartTime)/1000000;
        bool isGoal = throwDice((getShoot(states[shooter])*7)/(shootPenalty*10), blockShoot, rnds[1]) == 0;
        if (isGoal) {
            uint8 assister = selectAssister(matchStartTime, states, playersPerZone, extraAttack, shooter, rnds[2]);
            matchLog[teamThatAttacks] = addAssister(matchLog[teamThatAttacks], assister, currentGoals);
            matchLog[teamThatAttacks] = addShooter(matchLog[teamThatAttacks], shooter, currentGoals);
            matchLog[teamThatAttacks] = addForwardPos(matchLog[teamThatAttacks], getForwardPos(shooter, playersPerZone), currentGoals);
            matchLog[teamThatAttacks]++; // adds 1 goal because nGoals is the right-most number serialized
        }
        return matchLog;
    }
    
    function getForwardPos(uint8 posInLineUp, uint8[9] memory playersPerZone) private pure returns (uint8) {
        if (posInLineUp == 0) return 0;
        else if (posInLineUp < 1 + getNDefenders(playersPerZone)) return 1;
        else if (posInLineUp < 1 + getNDefenders(playersPerZone)+ getNMidfielders(playersPerZone)) return 2;
        else return 3;
    }
    
    
}

