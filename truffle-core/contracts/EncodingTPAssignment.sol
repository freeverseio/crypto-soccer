pragma solidity >=0.4.21 <0.6.0;

contract EncodingTPAssignment {

    // 5 bit for specialPlayer
    uint256 public constant MAX_WEIGHT = 75; 
    uint256 public constant MIN_WEIGHT = 5; 
    uint8  public constant PLAYERS_PER_TEAM_MAX  = 25;
    // We have 5 buckets: GK, D, M, A, Special
    // We need 5 weights per bucket
    //      that sum(weights) = MAX_WEIGHT
    // 7 bit per weight
    // 5 bit for specialPlayer
    
    
    function encodeTP(uint8[25] memory skillWeights, uint8 specialPlayer) public pure returns (uint256 encoded) {
        for (uint8 bucket = 0; bucket < 5; bucket++) {
            uint256 sum = 0;
            for (uint8 sk = 5 * bucket; sk < 5 * (bucket+1); sk++) {
                uint8 weight = skillWeights[sk];
                sum += weight;
                encoded |= uint256(MIN_WEIGHT + weight) << 7 * sk;
            }
            require(sum == MAX_WEIGHT, "weights too large");
        }
        require(specialPlayer < PLAYERS_PER_TEAM_MAX, "specialPlayer value too large");
        encoded |= uint256(specialPlayer) << 175;
    } 

    function decodeTP(uint256 encoded) public pure returns(uint8[25] memory skillWeights, uint8 specialPlayer) {
        for (uint8 bucket = 0; bucket < 5; bucket++) {
            for (uint8 sk = 5 * bucket; sk < 5* (bucket+1); sk++) {
                skillWeights[sk] = uint8((encoded >> 7 * sk) & 127);
            }
        }
        specialPlayer = uint8((encoded >> 175) & 31);
    } 
        
}