pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkillsGetters {

    /**
     * @dev PlayerSkills serializes a total of 148 bits:  6*14 + 4 + 3+ 3 + 43 + 1 + 1 + 3 + 3 + 3
     *      5 skills                  = 5 x 16 bits
     *                                = shoot, speed, pass, defence, endurance
     *      dayOfBirth                = 16 bits  (since Unix time, max 180 years)
     *      birthTraits               = variable num of bits: [potential, forwardness, leftishness, aggressiveness]
     *      potential                 = 4 bits (number is limited to [0,...,9])
     *      forwardness               = 3 bits
     *                                  GK: 0, D: 1, M: 2, F: 3, MD: 4, MF: 5
     *      leftishness               = 3 bits, in boolean triads: (L, C, R):
     *                                  0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
     *      aggressiveness            = 3 bits
     *      playerId                  = 43 bits
     *      
     *      alignedEndOfLastHalf      = 1b (bool)
     *      redCardLastGame           = 1b (bool)
     *      gamesNonStopping          = 3b (0, 1, ..., 6). Finally, 7 means more than 6.
     *      injuryWeeksLeft           = 3b 
     *      substitutedFirstHalf      = 1b (bool) 
     *      sumSkills                 = 19b (must equal sum(skills), of if each is 16b, this can be at most 5x16b => use 19b)
     *      isSpecialPlayer           = 1b (set at the left-most bit, 255)
     *      targetTeamId              = 43b
     *      generation                = 8b. From [0,...,31] => not-a-child, from [32,..63] => a child
    **/
    function getSkill(uint256 encodedSkills, uint8 skillIdx) public pure returns (uint256) {
        return (encodedSkills >> (uint256(skillIdx) * 16)) & 65535; // 65535 = 2**16 - 1
    }

    function getBirthDay(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 80 & 65535);
    }

    function getPlayerIdFromSkills(uint256 encodedSkills) public pure returns (uint256) {
        if (getIsSpecial(encodedSkills)) return encodedSkills;
        return uint256(encodedSkills >> 96 & 8796093022207); // 2**43 - 1 = 8796093022207
    }

    function getPotential(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 139 & 15);
    }

    function getForwardness(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 143 & 7);
    }

    function getLeftishness(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 146 & 7);
    }

    function getAggressiveness(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 149 & 7);
    }

    function getAlignedEndOfFirstHalf(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 152 & 1) == 1;
    }

    function getRedCardLastGame(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 153 & 1) == 1;
    }

    function getGamesNonStopping(uint256 encodedSkills) public pure returns (uint8) {
        return uint8(encodedSkills >> 154 & 7);
    }

    function getInjuryWeeksLeft(uint256 encodedSkills) public pure returns (uint8) {
        return uint8(encodedSkills >> 157 & 7);
    }

    function getSubstitutedFirstHalf(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 160 & 1) == 1;
    }

    function getSumOfSkills(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 161 & 524287); // 2**19-1
    }
    
    function getIsSpecial(uint256 encodedSkills) public pure returns (bool) {
        return uint256(encodedSkills >> 255 & 1) == 1; 
    }
     
    function addIsSpecial(uint256 encodedSkills) public pure returns (uint256) {
        return (encodedSkills | (uint256(1) << 255));
    }

    function getGeneration(uint256 encodedSkills) public pure returns (uint256) {
        return (encodedSkills >> 223) & 255;
    }
}
