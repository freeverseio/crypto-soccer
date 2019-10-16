pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLog {

    uint256 private constant ONE256       = 1; 
 
    function encodeMatchLog(
        uint8 nGoals, // 4b
        uint8[14] memory assistersIdx, // 4b each
        uint8[14] memory shootersIdx, // 4b each
        uint8[14] memory shooterForwardPos, // 2b each
        bool[7] memory penalties, // 1b each
        uint8[6] memory outOfGamesAndYellowCards,  // 2 outOfGames, 4 yellowCards, 4b each
        uint8[2] memory outOfGameRounds,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        bool[2] memory yellowCardedFinished1stHalf, // 1b each
        uint8[3] memory halfTimeSubstitutions, // 4b each, the first 3 for half 1, the other for half 2
        bool[6] memory ingameSubstitutions // 1b each, the first 3 for half 1, the other for half 2
    )
        public
        pure
        returns (uint256 log) 
    {
        log = nGoals;
        for (uint8 p = 0; p < 14; p++) {
            log |= uint256(assistersIdx[p]) << 4 + 4 * p;
            log |= uint256(shootersIdx[p]) << 60 + 4 * p;
            log |= uint256(shooterForwardPos[p]) << 116 + 2 * p;
        }            
        for (uint8 p = 0; p < 7; p++) {
            log |= uint256(penalties[p] ? 1: 0) << 144 + p;
        }            
        // 1st half
        log |= uint256(outOfGamesAndYellowCards[0]) << 151;
        log |= uint256(outOfGameRounds[0]) << 155;
        log |= uint256(typesOutOfGames[0]) << 159;
        log |= uint256(outOfGamesAndYellowCards[2]) << 161;
        log |= uint256(outOfGamesAndYellowCards[3]) << 165;
        log |= uint256(yellowCardedFinished1stHalf[0] ? 1: 0) << 169;
        log |= uint256(yellowCardedFinished1stHalf[1] ? 1: 0) << 170;
        // 2nd half
        log |= uint256(outOfGamesAndYellowCards[1]) << 171;
        log |= uint256(outOfGameRounds[1]) << 175;
        log |= uint256(typesOutOfGames[1]) << 179;
        log |= uint256(outOfGamesAndYellowCards[4]) << 181;
        log |= uint256(outOfGamesAndYellowCards[5]) << 185;
        // ingameSubstitutions
        for (uint8 p = 0; p < 6; p++) {
            log |= uint256(ingameSubstitutions[p] ? 1: 0) << 189 + p;
        }        
        for (uint8 p = 0; p < 3; p++) {
            log |= uint256(halfTimeSubstitutions[p]) << 195 + 4 * p;
        }            
    }
    
    
    function decodeMatchLog(uint256 log) public pure returns(
        uint8 nGoals, // 4b
        uint8[14] memory assistersIdx, // 4b each
        uint8[14] memory shootersIdx, // 4b each
        uint8[14] memory shooterForwardPos, // 2b each
        bool[15] memory penalties, // 1b each
        uint8[6] memory outOfGamesAndYellowCards,  // 2 outOfGames, 4 yellowCards, 4b each
        uint8[2] memory outOfGameRounds,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        bool[2] memory yellowCardedFinished1stHalf, // 1b each
        uint8[3] memory halfTimeSubstitutions, // 4b each, the first 3 for half 1, the other for half 2
        bool[6] memory ingameSubstitutions // 1b each, the first 3 for half 1, the other for half 2
    ) 
    {
        nGoals = uint8(log & 15);
        for (uint8 p = 0; p < 14; p++) {
            assistersIdx[p] = uint8((log >> 4 + 4 * p) & 15);
            shootersIdx[p] = uint8((log >> 60 + 4 * p) & 15);
            shooterForwardPos[p] = uint8((log >> 116 + 2 * p) & 3);
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
        yellowCardedFinished1stHalf[0] = ((log >> 169) & 1) == 1;
        yellowCardedFinished1stHalf[1] = ((log >> 170) & 1) == 1;
        // 2nd half
        outOfGamesAndYellowCards[1] = uint8((log >> 171) & 15);
        outOfGameRounds[1] = uint8((log >> 175) & 15);
        typesOutOfGames[1] = uint8((log >> 179) & 3);
        outOfGamesAndYellowCards[4] = uint8((log >> 181) & 15);
        outOfGamesAndYellowCards[5] = uint8((log >> 185) & 15);
        // ingameSubstitutions
        for (uint8 p = 0; p < 6; p++) {
            ingameSubstitutions[p] = ((log >> 189 + p) & 1) == 1;
        }        
        for (uint8 p = 0; p < 3; p++) {
            halfTimeSubstitutions[p]  = uint8((log >> 195 + 4 * p) & 15);
        }            
    }
}