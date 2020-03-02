pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkillsSetters {

    function setSkill(uint256 encodedSkills, uint256 val, uint8 skillIdx) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(65535) << (16 * uint256(skillIdx)))) | (val << (16 * uint256(skillIdx)));
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
