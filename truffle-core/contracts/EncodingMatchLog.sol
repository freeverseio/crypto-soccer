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
        uint8[2] memory outOfGames,  // 4b each
        uint8[2] memory outOfGameRounds,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        uint8[4] memory yellowCards, // 4b each
        bool[6] memory substitutions // 4b each, the first 3 for half 1, the other for half 2
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
        log |= uint256(outOfGames[0]) << 151;
        log |= uint256(outOfGameRounds[0]) << 155;
        log |= uint256(typesOutOfGames[0]) << 159;
        log |= uint256(yellowCards[0]) << 161;
        log |= uint256(yellowCards[1]) << 165;
        // 2nd half
        log |= uint256(outOfGames[1]) << 169;
        log |= uint256(outOfGameRounds[1]) << 173;
        log |= uint256(typesOutOfGames[1]) << 177;
        log |= uint256(yellowCards[2]) << 179;
        log |= uint256(yellowCards[3]) << 183;
        // substitutions
        for (uint8 p = 0; p < 6; p++) {
            log |= uint256(substitutions[p] ? 1: 0) << 186 + p;
        }        
    }
    
    
    function decodeMatchLog(uint256 log) public pure returns(
        uint8 nGoals, // 4b
        uint8[14] memory assistersIdx, // 4b each
        uint8[14] memory shootersIdx, // 4b each
        uint8[14] memory shooterForwardPos, // 2b each
        bool[15] memory penalties, // 1b each
        uint8[2] memory outOfGames,  // 4b each
        uint8[2] memory outOfGameRounds,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        uint8[4] memory yellowCards, // 4b each
        bool[6] memory substitutions // 4b each
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
        outOfGames[0] = uint8((log >> 151) & 15);
        outOfGameRounds[0] = uint8((log >> 155) & 15);
        typesOutOfGames[0] = uint8((log >> 159) & 3);
        yellowCards[0] = uint8((log >> 161) & 15);
        yellowCards[1] = uint8((log >> 165) & 15);
        // 2nd half
        outOfGames[1] = uint8((log >> 169) & 15);
        outOfGameRounds[1] = uint8((log >> 173) & 15);
        typesOutOfGames[1] = uint8((log >> 177) & 3);
        yellowCards[2] = uint8((log >> 179) & 15);
        yellowCards[3] = uint8((log >> 183) & 15);
        // substitutions
        for (uint8 p = 0; p < 6; p++) {
            substitutions[p] = ((log >> 186 + p) & 1) == 1;
        }        
    }
}