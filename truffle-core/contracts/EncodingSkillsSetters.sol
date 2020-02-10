pragma solidity >=0.5.12 <0.6.2;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkillsSetters {

    uint8 constant private PLAYERS_PER_TEAM_INIT = 18;
    uint8 constant private PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant private MIN_PLAYER_AGE_AT_BIRTH = 16;
    uint8 constant private MAX_PLAYER_AGE_AT_BIRTH = 32;
    uint8 constant private N_SKILLS = 5;

    // Birth Traits: potential, forwardness, leftishness, aggressiveness
    uint8 constant private IDX_POT = 0;
    uint8 constant private IDX_FWD = 1;
    uint8 constant private IDX_LEF = 2;
    uint8 constant private IDX_AGG = 3;
    // prefPosition idxs: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    uint8 constant private IDX_GK = 0;
    uint8 constant private IDX_D  = 1;
    uint8 constant private IDX_M  = 2;
    uint8 constant private IDX_F  = 3;
    uint8 constant private IDX_MD = 4;
    uint8 constant private IDX_MF = 5;
    //  Leftishness:   0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    uint8 constant private IDX_R = 1;
    uint8 constant private IDX_C = 2;
    uint8 constant private IDX_CR = 3;
    uint8 constant private IDX_L = 4;
    uint8 constant private IDX_LR = 5;
    uint8 constant private IDX_LC = 6;
    uint8 constant private IDX_LCR = 7;

    
    function setShoot(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return ((encodedSkills & ~(uint256(65535))) | val);
    }
    
    function setSpeed(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(65535) << 16)) | (val << 16);
    }
    
    function setPass(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(65535) << 32)) | (val << 32);
    }
    
    function setDefence(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(65535) << 48)) | (val << 48);
    }

    function setEndurance(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(65535) << 64)) | (val << 64);
    }

    function setPotential(uint256 encodedSkills, uint256 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(15) << 139)) | (val << 139);
    }

    function setAlignedEndOfFirstHalf(uint256 encodedSkills, bool val) public pure returns (uint256) {
        if (val) return (encodedSkills & ~(uint256(1) << 152)) | (uint256(1) << 152);
        else return (encodedSkills & ~(uint256(1) << 152));
    }

    function setSubstitutedFirstHalf(uint256 encodedSkills, bool val) public pure returns (uint256) {
        if (val) return (encodedSkills & ~(uint256(1) << 160)) | (uint256(1) << 160);
        else return (encodedSkills & ~(uint256(1) << 160));
    }

    function setRedCardLastGame(uint256 encodedSkills, bool val) public pure returns (uint256) {
        if (val) return (encodedSkills & ~(uint256(1) << 153)) | (uint256(1) << 153);
        else return (encodedSkills & ~(uint256(1) << 153));
    }

    function setInjuryWeeksLeft(uint256 encodedSkills, uint8 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(7) << 157)) | (uint256(val) << 157);
    }

    function setSumOfSkills(uint256 encodedSkills, uint32 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(524287) << 161)) | (uint256(val) << 161);
    }

    function setGeneration(uint256 encodedSkills, uint32 val) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(255) << 223)) | (uint256(val) << 223);
    }

    function setGamesNonStopping(uint256 encodedSkills, uint8 val) public pure returns (uint256) {
        require(val < 8, "gamesNonStopping out of bound");
        return (encodedSkills & ~(uint256(7) << 154)) | (uint256(val) << 154);
    }
}
