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
        bool[7] memory penalties, // 1b each
        uint8[2] memory outOfGames,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        uint8[4] memory yellowCards // 4b each
    )
        public
        pure
        returns (uint256 log) 
    {
        log = nGoals;
        for (uint8 p = 0; p < 14; p++) {
            log |= uint256(assistersIdx[p]) << 4 + 4 * p;
            log |= uint256(shootersIdx[p]) << 60 + 4 * p;
        }            
        for (uint8 p = 0; p < 7; p++) {
            log |= uint256(penalties[p] ? 1: 0) << 116 + p;
        }            
        // 1st half
        log |= uint256(outOfGames[0]) << 123;
        log |= uint256(typesOutOfGames[0]) << 127;
        log |= uint256(yellowCards[0]) << 129;
        log |= uint256(yellowCards[1]) << 133;
        // 2nd half
        log |= uint256(outOfGames[1]) << 137;
        log |= uint256(typesOutOfGames[1]) << 141;
        log |= uint256(yellowCards[2]) << 143;
        log |= uint256(yellowCards[3]) << 147;
    }
    
    
    function decodeMatchLog(uint256 log) public pure returns(
        uint8 nGoals, // 4b
        uint8[14] memory assistersIdx, // 4b each
        uint8[14] memory shootersIdx, // 4b each
        bool[15] memory penalties, // 1b each
        uint8[2] memory outOfGames,  // 4b each
        uint8[2] memory typesOutOfGames, // 2b each
        uint8[4] memory yellowCards // 4b each
    ) 
    {
        nGoals = uint8(log & 15);
        for (uint8 p = 0; p < 14; p++) {
            assistersIdx[p] = uint8((log >> 4 + 4 * p) & 15);
            shootersIdx[p] = uint8((log >> 60 + 4 * p) & 15);
        }    
        for (uint8 p = 0; p < 7; p++) {
            penalties[p] = ((log >> 116 + p) & 1) == 1;
        }            
        // 1st half
        outOfGames[0] = uint8((log >> 123) & 15);
        typesOutOfGames[0] = uint8((log >> 127) & 3);
        yellowCards[0] = uint8((log >> 129) & 15);
        yellowCards[1] = uint8((log >> 133) & 15);
        // 2nd half
        outOfGames[1] = uint8((log >> 137) & 15);
        typesOutOfGames[1] = uint8((log >> 141) & 3);
        yellowCards[2] = uint8((log >> 143) & 15);
        yellowCards[3] = uint8((log >> 147) & 15);
    }



}