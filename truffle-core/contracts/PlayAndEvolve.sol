pragma solidity ^0.5.0;

import "./Evolution.sol";

contract PlayAndEvolve {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 private constant IDX_IS_2ND_HALF      = 0; 

    Evolution private _evo;
    Engine private _engine;

    function setEvolutionAddress(address addr) public {
        _evo = Evolution(addr);
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
        return _evo.computeTrainingPoints(
            _engine.playHalfMatch(seed, matchStartTime, states, tactics, matchLog, matchBools)
        );
    }
    
}

