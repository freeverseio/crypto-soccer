pragma solidity >=0.5.12 <0.6.2;

import "./TrainingPoints.sol";
import "./Evolution.sol";
import "./Engine.sol";

contract PlayAndEvolve {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 private constant IDX_IS_2ND_HALF      = 0; 
    uint8 public constant ROUNDS_PER_MATCH  = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)

    TrainingPoints private _evo;
    Engine private _engine;

    function setEvolutionAddress(address addr) public {
        _evo = TrainingPoints(addr);
    }
    function setEngine(address addr) public {
        _engine = Engine(addr);
    }

    function play2ndHalfAndEvolve(
        uint256 seed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public view returns(uint256[2] memory)
    {
        require(matchBools[IDX_IS_2ND_HALF], "play with evolution should only be called in 2nd half games");
        uint256[2+5*ROUNDS_PER_MATCH] memory matchLogsAndEvents = 
            _engine.playHalfMatch(seed, matchStartTime, states, tactics, matchLog, matchBools);

        uint256[2] memory matchLogs;
        matchLogs[0] = matchLogsAndEvents[0];
        matchLogs[1] = matchLogsAndEvents[1];
    
        // _evo.updateStatesAfterPlayHalf(states, tactics, matchLog, matchBools);
        return _evo.computeTrainingPoints(matchLogs);
    }
    
    // seedAndStartTimeAndEvents => logAnd... (in other code)



}

