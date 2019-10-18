pragma solidity ^0.5.0;


contract Evolution {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    
    function computeTrainingPoints( 
        uint256[2] memory matchLog,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf1,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf2,
        uint256[2] memory tacticsHalf1,
        uint256[2] memory tacticsHalf2,
        bool isHomeStadium
    )
        public
        pure
        returns (uint16[2] memory points)
    {
        // 1 point for winning at home, 2 points for winning
        // away, or in a cup match. 0 points for drawing.
        if ((matchLog[0] & 15) > (matchLog[1] & 15)) {
            points[0] = isHomeStadium ? 1 : 2;    
        } else if ((matchLog[0] & 15) < (matchLog[1] & 15)) {
            points[1] = isHomeStadium ? 2 : 2;    
        }

    }
    
}

