pragma solidity >=0.5.12 <0.6.2;

import "./TrainingPoints.sol";
import "./Evolution.sol";
import "./Engine.sol";

contract PlayAndEvolve {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 private constant IDX_IS_2ND_HALF      = 0; 
    uint8 public constant ROUNDS_PER_MATCH  = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)

    TrainingPoints private _training;
    Evolution private _evo;
    Engine private _engine;

    function setTrainingAddress(address addr) public {
        _training = TrainingPoints(addr);
    }
 
    function setEvolutionAddress(address addr) public {
        _evo = Evolution(addr);
    }
 
    function setEngineAddress(address addr) public {
        _engine = Engine(addr);
    }

    function generateMatchSeed(bytes32 seed, uint256[2] memory teamIds) public pure returns (uint256) {
        return uint256(keccak256(abi.encode(seed, teamIds[0], teamIds[1])));
    }

    // In a 1st half we need to:
    //      1. applyTrainingPoints: (oldStates, assignedTPs) => (newStates)
    //      2. playHalfMatch: (newStates) => (matchLogs, events)
    //      3. updateStatesAfterPlayHalf: (newStates, matchLogs) => finalStates
    // Output: (finalStates, matchLogsAndEvents)
    function play1stHalfAndEvolve(
        bytes32 verseSeed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory teamIds,
        uint256[2] memory tactics,
        uint256[2] memory matchLogs,
        bool[3] memory matchBools, // [is2ndHalf, isHomeStadium, isPlayoff]
        uint256[2] memory assignedTPs
    )
        public view returns(uint256[PLAYERS_PER_TEAM_MAX][2] memory, uint256[2+5*ROUNDS_PER_MATCH] memory)
    {
        require(!matchBools[IDX_IS_2ND_HALF], "play1stHalfAndEvolve was called with the wrong is2ndHalf boolean!");

        states[0] = _training.applyTrainingPoints(states[0], assignedTPs[0], matchStartTime, _evo.getTrainingPoints(matchLogs[0]));
        states[1] = _training.applyTrainingPoints(states[1], assignedTPs[1], matchStartTime, _evo.getTrainingPoints(matchLogs[1]));
        
        uint256[2] memory nullLogs;
        uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents = 
            _engine.playHalfMatch(generateMatchSeed(verseSeed, teamIds), matchStartTime, states, tactics, nullLogs, matchBools);

        states[0] = _evo.updateStatesAfterPlayHalf(states[0], matchLogsAndEvents[0], tactics[0], false);
        states[1] = _evo.updateStatesAfterPlayHalf(states[1], matchLogsAndEvents[1], tactics[1], false);

        return (states, matchLogsAndEvents);
    }
    
    
    // In a 2nd half we need to:
    //      1. playHalfMatch: (oldStates, matchLogsHalf1) => (matchLogsHalf2, events)
    //      2. updateStatesAfterPlayHalf: (oldStates, matchLogsHalf2) => newStates
    //      3. computeTrainingPoints: (matchLogsHalf2) => (matchLogsHalf2 with TPs)
    // Output: (newStates, matchLogsAndEvents with TPs)
    function play2ndHalfAndEvolve(
        bytes32 verseSeed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory teamIds,
        uint256[2] memory tactics,
        uint256[2] memory matchLogs,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public view returns(uint256[PLAYERS_PER_TEAM_MAX][2] memory, uint256[2+5*ROUNDS_PER_MATCH] memory)
    {
        require(matchBools[IDX_IS_2ND_HALF], "play2ndHalfAndEvolve was called with the wrong is2ndHalf boolean!");

        uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents = 
            _engine.playHalfMatch(generateMatchSeed(verseSeed, teamIds), matchStartTime, states, tactics, matchLogs, matchBools);

        states[0] = _evo.updateStatesAfterPlayHalf(states[0], matchLogsAndEvents[0], tactics[0], true);
        states[1] = _evo.updateStatesAfterPlayHalf(states[1], matchLogsAndEvents[1], tactics[1], true);

        (matchLogsAndEvents[0], matchLogsAndEvents[1]) = _training.computeTrainingPoints(matchLogsAndEvents[0], matchLogsAndEvents[1]);

        return (states, matchLogsAndEvents);
    }
    

}

