pragma solidity >=0.4.21 <0.6.0;

contract EncodingTPAssignment {

    uint16 public constant MAX_PERCENT = 80; 
    uint16 public constant MIN_PERCENT = 5; 
    uint8  public constant PLAYERS_PER_TEAM_MAX  = 25;
    // We have 5 buckets: GK, D, M, A, Special
    // We need 5 TPperSkill per bucket
    //      that sum(TPperSkill) < TP
    // 9 bit per TP
    // 9 bit per each of the TPperSkill
    // 5 bit for specialPlayer
    
    
    function encodeTP(uint16 TP, uint16[25] memory TPperSkill, uint8 specialPlayer) public pure returns (uint256 encoded) {
        uint16 minRHS = MIN_PERCENT * TP;
        uint16 maxRHS = MAX_PERCENT * TP;
        for (uint8 bucket = 0; bucket < 5; bucket++) {
            uint256 sum = 0;
            for (uint8 sk = 5 * bucket; sk < 5 * (bucket+1); sk++) {
                uint16 skill = TPperSkill[sk];
                require((100*skill >= minRHS) && (100*skill <= maxRHS), "one of the assigned TPs is too large or too small");
                sum += skill;
                encoded |= uint256(skill) << 9 * sk;
            }
            require(sum <= TP, "sum of Traning Points is too large");
        }
        require(specialPlayer < PLAYERS_PER_TEAM_MAX, "specialPlayer value too large");
        encoded |= uint256(TP) << 225;
        encoded |= uint256(specialPlayer) << 234;
    } 

    function decodeTP(uint256 encoded) public pure returns(uint16[25] memory TPperSkill, uint8 specialPlayer, uint16 TP) {
        for (uint8 bucket = 0; bucket < 5; bucket++) {
            for (uint8 sk = 5 * bucket; sk < 5* (bucket+1); sk++) {
                TPperSkill[sk] = uint16((encoded >> 9 * sk) & 511);
            }
        }
        return (TPperSkill, uint8((encoded >> 234) & 31), uint16((encoded >> 225) & 511));
    } 
        
}