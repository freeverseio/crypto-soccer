pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLog {

    uint256 private constant ONE256       = 1; 
  
    function encodeMatchLog(
        uint8 nGoals, // 4b, offset 0
        uint8[42] memory assistersShootersForwardsPos, 
            // [assistersIdx[14], shootersIdx[14], shooterFwdPos[14]]
            // [ each 4b, each 4b, each 2b]
            // [ offset 4, offset 60, offset 116]
        bool[7] memory penalties, // 1b each, offset 144
        uint8[6] memory outOfGamesAndYellowCards,  // 2 outOfGames, 4 yellowCards, 4b each, offset 151
        uint8[2] memory outOfGameRounds,  
        uint8[2] memory typesOutOfGames, 
        bool[3] memory yellowCardedDidNotFinish1stHalfAndIsHomeStadium, // 1b each, last one has offset 227
        uint8[3] memory halfTimeSubstitutions, // 4b each, offset 201, the first 3 for half 1, the other for half 2
        uint8[6] memory ingameSubs, // 2b each, offset 189
                                    //  ...the first 3 for half 1, the other for half 2.
                                    // ...0: no change required, 1: change happened, 2: change could not happen  
        uint8[4] memory numDefTotWinner
                                    // [4b, 4b, 4b, 2b], offset 213
                                    // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
                                    // winner: 0 = home, 1 = away, 2 = draw
    )
        public
        pure
        returns (uint256 log) 
    {
        log = nGoals;
        for (uint8 p = 0; p < 14; p++) {
            log |= uint256(assistersShootersForwardsPos[p]) << 4 + 4 * p;
            log |= uint256(assistersShootersForwardsPos[p + 14]) << 60 + 4 * p;
            log |= uint256(assistersShootersForwardsPos[p + 28]) << 116 + 2 * p;
        }            
        for (uint8 p = 0; p < 7; p++) {
            log |= uint256(penalties[p] ? 1: 0) << 144 + p;
        }            
        // 1st half
        log |= uint256(outOfGamesAndYellowCards[0]) << 151; // redCard
        log |= uint256(outOfGameRounds[0]) << 155;
        log |= uint256(typesOutOfGames[0]) << 159;
        log |= uint256(outOfGamesAndYellowCards[2]) << 161; // yellowCard
        log |= uint256(outOfGamesAndYellowCards[3]) << 165; // yellowCard
        log |= uint256(yellowCardedDidNotFinish1stHalfAndIsHomeStadium[0] ? 1: 0) << 169;
        log |= uint256(yellowCardedDidNotFinish1stHalfAndIsHomeStadium[1] ? 1: 0) << 170;
        // 2nd half
        log |= uint256(outOfGamesAndYellowCards[1]) << 171; // redCard
        log |= uint256(outOfGameRounds[1]) << 175;
        log |= uint256(typesOutOfGames[1]) << 179;
        log |= uint256(outOfGamesAndYellowCards[4]) << 181; // yellowCard
        log |= uint256(outOfGamesAndYellowCards[5]) << 185; // yellowCard
        // ingameSubs
        for (uint8 p = 0; p < 6; p++) {
            log |= uint256(ingameSubs[p]) << 189 + 2*p;
        }        
        for (uint8 p = 0; p < 3; p++) {
            log |= uint256(halfTimeSubstitutions[p]) << 201 + 4 * p;
            log |= uint256(numDefTotWinner[p]) << 213 + 4 * p; // nDefs1, nDefs2, nTot2
        }            
        log |= uint256(numDefTotWinner[3]) << 225; // winner
        log |= uint256(yellowCardedDidNotFinish1stHalfAndIsHomeStadium[2] ? 1: 0) << 227; // isHomeStadium
        
    }
    
    
    function decodeMatchLog(uint256 log) public pure returns(
        uint8 nGoals, // 4b
        uint8[42] memory assistersShootersForwardsPos, // 4b each
        bool[15] memory penalties, // 1b each
        uint8[6] memory outOfGamesAndYellowCards,  // 2 outOfGames, 4 yellowCards, 4b each
        uint8[2] memory outOfGameRounds,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        bool[3] memory yellowCardedDidNotFinish1stHalfAndIsHomeStadium, // 1b each
        uint8[3] memory halfTimeSubstitutions, // 4b each, the first 3 for half 1, the other for half 2
        uint8[6] memory ingameSubs, // 1b each, the first 3 for half 1, the other for half 2
        uint8[4] memory numDefTotWinner
                                    // [4b, 4b, 4b, 2b], offset 213
                                    // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
                                    // winner: 0 = home, 1 = away, 2 = draw
    ) 
    {
        nGoals = uint8(log & 15);
        for (uint8 p = 0; p < 14; p++) {
            assistersShootersForwardsPos[p]     = uint8((log >> 4 + 4 * p) & 15);
            assistersShootersForwardsPos[p+14]  = uint8((log >> 60 + 4 * p) & 15);
            assistersShootersForwardsPos[p+28]  = uint8((log >> 116 + 2 * p) & 3);
        }    
        for (uint8 p = 0; p < 7; p++) {
            penalties[p] = ((log >> 144 + p) & 1) == 1;
        }            
        // 1st half
        outOfGamesAndYellowCards[0] = uint8((log >> 151) & 15);
        outOfGameRounds[0] = uint8((log >> 155) & 15);
        typesOutOfGames[0] = uint8((log >> 159) & 3);
        outOfGamesAndYellowCards[2] = uint8((log >> 161) & 15);
        outOfGamesAndYellowCards[3] = uint8((log >> 165) & 15);
        yellowCardedDidNotFinish1stHalfAndIsHomeStadium[0] = ((log >> 169) & 1) == 1;
        yellowCardedDidNotFinish1stHalfAndIsHomeStadium[1] = ((log >> 170) & 1) == 1;
        // 2nd half
        outOfGamesAndYellowCards[1] = uint8((log >> 171) & 15);
        outOfGameRounds[1] = uint8((log >> 175) & 15);
        typesOutOfGames[1] = uint8((log >> 179) & 3);
        outOfGamesAndYellowCards[4] = uint8((log >> 181) & 15);
        outOfGamesAndYellowCards[5] = uint8((log >> 185) & 15);
        // ingameSubs
        for (uint8 p = 0; p < 6; p++) {
            ingameSubs[p] = uint8((log >> 189 + 2 * p) & 3);
        }        
        for (uint8 p = 0; p < 3; p++) {
            halfTimeSubstitutions[p]  = uint8((log >> 201 + 4 * p) & 15);
            numDefTotWinner[p] = uint8((log >> 213 + 4 * p) & 15);
        }            
        numDefTotWinner[3] = uint8((log >> 225) & 3);
        yellowCardedDidNotFinish1stHalfAndIsHomeStadium[2] = ((log >> 227) & 1) == 1;
    }
}