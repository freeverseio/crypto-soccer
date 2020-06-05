pragma solidity >= 0.6.3;

import "./TrainingPoints.sol";
import "./Evolution.sol";
import "./Engine.sol";
import "./Shop.sol";
import "../gameEngine/ErrorCodes.sol";
import "../encoders/EncodingTacticsBase1.sol";


/**
 @title Main entry point for backend. Plays 1st and 2nd half and evolves players.
 @author Freeverse.io, www.freeverse.io
 @dev All functions are basically pure, but some had to be made view
 @dev because they use a storage pointer to other contracts.
*/

contract PlayAndEvolve is ErrorCodes, EncodingTacticsBase1 {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 private constant IDX_IS_2ND_HALF = 0; 
    uint8 public constant ROUNDS_PER_MATCH = 12;   /// Number of relevant actions that happen during a game (12 equals one per 3.7 min)
    uint8 private constant IDX_IS_BOT_HOME      = 3; 
    uint8 private constant IDX_IS_BOT_AWAY      = 4; 
    
    TrainingPoints private training;
    Evolution private evo;
    Engine private engine;
    Shop private shop;

    constructor(address trainingAddr, address evolutionAddr, address engineAddr, address shopAddr) public {
        training = TrainingPoints(trainingAddr);
        evo = Evolution(evolutionAddr);
        engine = Engine(engineAddr);
        shop = Shop(shopAddr);
    }

    function generateMatchSeed(bytes32 seed, uint256[2] memory teamIds) public pure returns (uint256) {
        return uint256(keccak256(abi.encode(seed, teamIds[0], teamIds[1])));
    }

    /// In a 1st half we need to:
    ///      1. applyTrainingPoints: (oldSkills, assignedTPs) => (newSkills)
    ///      2. playHalfMatch: (newSkills) => (matchLogs, events)
    ///      3. updateSkillsAfterPlayHalf: (newSkills, matchLogs) => finalSkills
    /// Output: (finalSkills, matchLogsAndEvents)
    function play1stHalfAndEvolve(
        bytes32 verseSeed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory skills,
        uint256[2] memory teamIds,
        uint256[2] memory tactics,
        uint256[2] memory matchLogs,
        bool[5] memory matchBools, /// [is2ndHalf, isHomeStadium, isPlayoff, isBotHome, isBotAway]
        uint256[2] memory assignedTPs
    )
        public 
        view 
        returns (
            uint256[PLAYERS_PER_TEAM_MAX][2] memory, 
            uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents,
            uint8 err
        )
    {
        if (matchBools[IDX_IS_2ND_HALF]) { return (skills, matchLogsAndEvents, ERR_IS2NDHALF); }

        for (uint8 team = 0; team < 2; team++) {
            if (matchBools[IDX_IS_BOT_HOME + team]) {
                tactics[team] = getBotTactics();
            } else {
                (skills[team], err) = training.applyTrainingPoints(skills[team], assignedTPs[team], tactics[team], matchStartTime, evo.getTrainingPoints(matchLogs[team]));
                if (err > 0) return (skills, matchLogsAndEvents, err);
                err = shop.validateItemsInTactics(tactics[team]);
                if (err > 0) return (skills, matchLogsAndEvents, err);
            }
        }
            
        /// Note that the following call does not change de values of "skills" because it calls a separate contract.
        /// It would do so if playHalfMatch was part of this contract code.
        (matchLogsAndEvents, err) = engine.playHalfMatch(
            generateMatchSeed(verseSeed, teamIds), 
            matchStartTime, 
            skills, 
            tactics, 
            [uint256(0),uint256(0)], 
            matchBools
        );

        for (uint8 team = 0; team < 2; team++) {
            (skills[team], err) = evo.updateSkillsAfterPlayHalf(skills[team], matchLogsAndEvents[team], tactics[team], false, matchBools[IDX_IS_BOT_HOME + team]);
            if (err > 0) return (skills, matchLogsAndEvents, err);
        }

        return (skills, matchLogsAndEvents, 0);
    }
    
    
    /// In a 2nd half we need to:
    ///      1. playHalfMatch: (oldSkills, matchLogsHalf1) => (matchLogsHalf2, events)
    ///      2. updateSkillsAfterPlayHalf: (oldSkills, matchLogsHalf2) => newSkills
    ///      3. computeTrainingPoints: (matchLogsHalf2) => (matchLogsHalf2 with TPs)
    /// Output: (newSkills, matchLogsAndEvents with TPs)
    function play2ndHalfAndEvolve(
        bytes32 verseSeed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory skills,
        uint256[2] memory teamIds,
        uint256[2] memory tactics,
        uint256[2] memory matchLogs,
        bool[5] memory matchBools /// [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public 
        view 
        returns(
            uint256[PLAYERS_PER_TEAM_MAX][2] memory, 
            uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents,
            uint8 err
        )
    {
        if (!matchBools[IDX_IS_2ND_HALF]) { return (skills, matchLogsAndEvents, ERR_IS2NDHALF); }

        for (uint8 team = 0; team < 2; team++) {
            if (matchBools[IDX_IS_BOT_HOME + team]) {
                tactics[team] = getBotTactics();
            } else {
                err = shop.validateItemsInTactics(tactics[team]);
                if (err > 0) return (skills, matchLogsAndEvents, err);
            }
        }

        /// Note that the following call does not change de values of "skills" because it calls a separate contract.
        /// It would do so if playHalfMatch was part of this contract code.
        (matchLogsAndEvents, err) = engine.playHalfMatch(
            generateMatchSeed(verseSeed, teamIds), 
            matchStartTime, 
            skills, 
            tactics, 
            matchLogs, 
            matchBools
        );
        if (err > 0) return (skills, matchLogsAndEvents, err);

        for (uint8 team = 0; team < 2; team++) {
            (skills[team], err) = evo.updateSkillsAfterPlayHalf(skills[team], matchLogsAndEvents[team], tactics[team], true, matchBools[IDX_IS_BOT_HOME + team]);
            if (err > 0) return (skills, matchLogsAndEvents, err);
        }

        (matchLogsAndEvents[0], matchLogsAndEvents[1]) = training.computeTrainingPoints(matchLogsAndEvents[0], matchLogsAndEvents[1]);

        return (skills, matchLogsAndEvents, 0);
    }
    

    function getBotTactics() public pure returns(uint256) { 
        return encodeTactics(
            [NO_SUBST, NO_SUBST, NO_SUBST], // no substitutions
            [0, 0, 0], // subRounds don't matter
            [0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25], // consecutive lineup, with one single GK 
            [false, false, false, false, false, false, false, false, false, false], // no extra attack
            1 // tacticsId = 1 = 5-4-1
        );
    }
}

