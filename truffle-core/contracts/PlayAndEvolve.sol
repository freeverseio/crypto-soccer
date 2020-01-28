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


    // In a 2nd half we need to:
    //      1. playHalfMatch: (oldStates, matchLogsHalf1) => (matchLogsHalf2, events)
    //      2. updateStatesAfterPlayHalf: (oldStates, matchLogsHalf2) => newStates
    //      3. computeTrainingPoints: (matchLogsHalf2) => (matchLogsHalf2 with TPs)
    // Output: (newStates, matchLogsAndEvents with TPs)
    function play2ndHalfAndEvolve(
        uint256 seed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public view returns(uint256[PLAYERS_PER_TEAM_MAX][2] memory, uint256[2+5*ROUNDS_PER_MATCH] memory)
    {
        require(matchBools[IDX_IS_2ND_HALF], "play2ndHalfAndEvolve was called with the wrong is2ndHalf boolean!");

        uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents = 
            _engine.playHalfMatch(seed, matchStartTime, states, tactics, matchLog, matchBools);

        states[0] = _evo.updateStatesAfterPlayHalf(states[0], matchLog[0], tactics[0], true);
        states[1] = _evo.updateStatesAfterPlayHalf(states[1], matchLog[1], tactics[1], true);

        (matchLogsAndEvents[0], matchLogsAndEvents[1]) = _training.computeTrainingPoints(matchLogsAndEvents[0], matchLogsAndEvents[1]);

        return (states, matchLogsAndEvents);
    }
    

}

