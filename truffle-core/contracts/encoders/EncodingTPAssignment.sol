pragma solidity >= 0.6.3;

/**
 @title Library to serialize/deserialize assignment of Traning points by users
 @author Freeverse.io, www.freeverse.io
*/

/**
 Spec: 
 We have 5 buckets: GK, D, M, A, Special
 We need 5 TPperSkill per bucket 
      - 9 bit per each of the TPperSkill
      - such that sum(TPperSkill) < TP (except for special player)
 assignedTP encodes a total: 5 buckets * 5 TPperSKill * 9b + 1 totalTP * 9b + 5 for specialPlId = 239
 offsets:
      - TPperSkill: 0 --> 224
      - TP: 225 --> 233
      - specIf --> 234 -> 238
 9 bit for TP  => max val = 511
 5 bit for specialPlayer
 TP: all the available Training point earned in the previous match log
 specialPlayer: no specialPlayer if == 25
*/

contract EncodingTPAssignment {

    uint16 public constant MAX_PERCENT = 60; 
    uint8 private constant PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_PLAYER = PLAYERS_PER_TEAM_MAX; /// No player chosen
    
    function encodeTP(uint16 TP, uint16[25] memory TPperSkill, uint8 specialPlayer) public pure returns (uint256 encoded) {
        require(specialPlayer <= PLAYERS_PER_TEAM_MAX, "specialPlayer value too large");

        encoded |= uint256(TP) << 225;
        encoded |= uint256(specialPlayer) << 234;

        uint16 maxRHS = (TP < 4) ? 100 * TP : MAX_PERCENT * TP;
        uint8 lastBucket = (specialPlayer == NO_PLAYER ? 4 : 5);
        for (uint8 bucket = 0; bucket < lastBucket; bucket++) {
            if (bucket == 4) {
                TP = uint16((uint256(TP) * 11)/10);
                maxRHS = (TP < 4) ? 100 * TP : MAX_PERCENT * TP;
            }
            uint256 sum = 0;
            for (uint8 sk = 5 * bucket; sk < 5 * (bucket+1); sk++) {
                uint16 skill = TPperSkill[sk];
                require(100*skill <= maxRHS, "one of the assigned TPs is too large");
                sum += skill;
                encoded |= uint256(skill) << 9 * sk;
            }
            require(sum <= TP, "sum of Traning Points is too large");
        }
    } 

    function decodeTP(uint256 encoded) public pure returns(uint16[25] memory TPperSkill, uint8 specialPlayer, uint16 TP) {
        TP = uint16((encoded >> 225) & 511);
        uint16 TPtemp = TP;
        specialPlayer = uint8((encoded >> 234) & 31);
        require(specialPlayer <= PLAYERS_PER_TEAM_MAX, "specialPlayer value too large");
        uint16 maxRHS = (TP < 4) ? 100 * TPtemp : MAX_PERCENT * TPtemp;
        for (uint8 bucket = 0; bucket < 5; bucket++) {
            if (bucket == 4) {
                TPtemp = uint16((uint256(TPtemp) * 11000)/10000);
                maxRHS = (TP < 4) ? 100 * TPtemp : MAX_PERCENT * TPtemp;
            }
            uint256 sum = 0;
            for (uint8 sk = 5 * bucket; sk < 5* (bucket+1); sk++) {
                uint16 skill = uint16((encoded >> 9 * sk) & 511);
                require(100*skill <= maxRHS, "one of the assigned TPs is too large or too small");
                TPperSkill[sk] = skill;
                sum += skill;
            }
            require(sum <= TPtemp, "sum of Traning Points is too large");
        }
    } 
}
