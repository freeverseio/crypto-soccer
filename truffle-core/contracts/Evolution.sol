pragma solidity ^0.5.0;


contract Evolution {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_CARD  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card

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
        returns (uint256[2] memory points)
    {
        
        // 1 point for winning at home, 2 points for winning
        // away, or in a cup match. 0 points for drawing.
        uint256 nGoals0 = (matchLog[0] & 15);
        uint256 nGoals1 = (matchLog[1] & 15);
        if ( nGoals0 > nGoals1) {
            points[0] = isHomeStadium ? 1 : 2;    
        } else if (nGoals0 < nGoals1) {
            points[1] = isHomeStadium ? 2 : 2;    
        }

        uint256[2] memory pointsNeg;

        // -1 for each opponent goal
        pointsNeg[0] = nGoals1;
        pointsNeg[1] = nGoals0;

        // -3 for redCards, -1 for yellows
        // ...note that offset for 1st half is 159, and for 2nd half is 179
        for (uint8 team = 0; team <2; team++) {
            for (uint8 offset = 159; offset < 180; offset += 20) {
                pointsNeg[team] += (((matchLog[0] >> offset) & 3) < RED_CARD) ? 3 : 0;
                pointsNeg[team] += (((matchLog[0] >> (offset + 2)) & 15) < NO_CARD) ? 1 : 0;
                pointsNeg[team] += (((matchLog[0] >> (offset + 6)) & 15) < NO_CARD) ? 1 : 0;
            }
        }

        

    }
    
}

