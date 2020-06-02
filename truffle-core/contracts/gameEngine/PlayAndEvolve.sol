pragma solidity >= 0.6.3;

import "./TrainingPoints.sol";
import "./Evolution.sol";
import "./Engine.sol";
import "./Shop.sol";
import "../gameEngine/ErrorCodes.sol";

/**
 @title Main entry point for backend. Plays 1st and 2nd half and evolves players.
 @author Freeverse.io, www.freeverse.io
 @dev All functions are basically pure, but some had to be made view
 @dev because they use a storage pointer to other contracts.
*/

contract PlayAndEvolve is ErrorCodes {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 private constant IDX_IS_2ND_HALF = 0; 
    uint8 public constant ROUNDS_PER_MATCH = 12;   /// Number of relevant actions that happen during a game (12 equals one per 3.7 min)

    TrainingPoints public training;
    Evolution public evo;
    Engine public engine;
    Shop public shop;

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
        bool[3] memory matchBools, /// [is2ndHalf, isHomeStadium, isPlayoff]
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

        (skills[0], err) = training.applyTrainingPoints(skills[0], assignedTPs[0], tactics[0], matchStartTime, evo.getTrainingPoints(matchLogs[0]));
        (skills[1], err) = training.applyTrainingPoints(skills[1], assignedTPs[1], tactics[1], matchStartTime, evo.getTrainingPoints(matchLogs[1]));
        
        /// Note that the following call does not change de values of "skills" because it calls a separate contract.
        /// It would do so if playHalfMatch was part of this contract code.

        shop.validateItemsInTactics(tactics[0]);
        shop.validateItemsInTactics(tactics[1]);
        
        matchLogsAndEvents = engine.playHalfMatch(
            generateMatchSeed(verseSeed, teamIds), 
            matchStartTime, 
            skills, 
            tactics, 
            [uint256(0),uint256(0)], 
            matchBools
        );

        skills[0] = evo.updateSkillsAfterPlayHalf(skills[0], matchLogsAndEvents[0], tactics[0], false);
        skills[1] = evo.updateSkillsAfterPlayHalf(skills[1], matchLogsAndEvents[1], tactics[1], false);

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
        bool[3] memory matchBools /// [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public 
        view 
        returns(
            uint256[PLAYERS_PER_TEAM_MAX][2] memory, 
            uint256[2+5*ROUNDS_PER_MATCH] memory
        )
    {
        require(matchBools[IDX_IS_2ND_HALF], "play2ndHalfAndEvolve was called with the wrong is2ndHalf boolean!");

        shop.validateItemsInTactics(tactics[0]);
        shop.validateItemsInTactics(tactics[1]);

        /// Note that the following call does not change de values of "skills" because it calls a separate contract.
        /// It would do so if playHalfMatch was part of this contract code.
        uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents = 
            engine.playHalfMatch(generateMatchSeed(verseSeed, teamIds), matchStartTime, skills, tactics, matchLogs, matchBools);

        skills[0] = evo.updateSkillsAfterPlayHalf(skills[0], matchLogsAndEvents[0], tactics[0], true);
        skills[1] = evo.updateSkillsAfterPlayHalf(skills[1], matchLogsAndEvents[1], tactics[1], true);

        (matchLogsAndEvents[0], matchLogsAndEvents[1]) = training.computeTrainingPoints(matchLogsAndEvents[0], matchLogsAndEvents[1]);

        return (skills, matchLogsAndEvents);
    }
    

}

