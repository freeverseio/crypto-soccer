pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkillsSetters {

    uint8 constant public PLAYERS_PER_TEAM_INIT = 18;
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public MIN_PLAYER_AGE_AT_BIRTH = 16;
    uint8 constant public MAX_PLAYER_AGE_AT_BIRTH = 32;
    uint8 constant public N_SKILLS = 5;

    // Birth Traits: potential, forwardness, leftishness, aggressiveness
    uint8 constant private IDX_POT = 0;
    uint8 constant private IDX_FWD = 1;
    uint8 constant private IDX_LEF = 2;
    uint8 constant private IDX_AGG = 3;
    // prefPosition idxs: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    uint8 constant public IDX_GK = 0;
    uint8 constant public IDX_D  = 1;
    uint8 constant public IDX_M  = 2;
    uint8 constant public IDX_F  = 3;
    uint8 constant public IDX_MD = 4;
    uint8 constant public IDX_MF = 5;
    //  Leftishness:   0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    uint8 constant public IDX_R = 1;
    uint8 constant public IDX_C = 2;
    uint8 constant public IDX_CR = 3;
    uint8 constant public IDX_L = 4;
    uint8 constant public IDX_LR = 5;
    uint8 constant public IDX_LC = 6;
    uint8 constant public IDX_LCR = 7;

    
    function setShoot(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(0x3fff) << 242)) | (val << 242);
    }
    
    function setSpeed(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(0x3fff) << 228)) | (val << 228);
    }
    
    function setPass(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(0x3fff) << 214)) | (val << 214);
    }
    
    function setDefence(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(0x3fff) << 200)) | (val << 200);
    }

    function setEndurance(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(0x3fff) << 186)) | (val << 186);
    }
}