pragma solidity ^0.5.0;


contract Evolution {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    
    function computeTrainingPoints( 
        uint256[2] memory matchLog,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf1,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf2,
        uint256[2] memory tacticsHalf1,
        uint256[2] memory tacticsHalf2
    )
        public
        pure
        returns (uint16[2] memory points)
    {
    }
    
}

