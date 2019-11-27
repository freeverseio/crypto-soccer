pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */
import "./EncodingSkills.sol";
import "./EncodingSkillsSetters.sol";

contract Privileged is EncodingSkills, EncodingSkillsSetters {
    
    // order of idxs:
    // skills: shoot, speed, pass, defence, endurance
    // birthTraits: potential, forwardness, leftishness, aggressiveness
    // prefPosition: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    // leftishness:   0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111

    function createSpecialPlayer(

        uint16[N_SKILLS] memory skillsVec,
        uint256 ageInSecs,
        uint8[4] memory birthTraits,
        uint256 playerId
    ) public view returns (uint256) {
        uint256 dayOfBirth = (now - ageInSecs/7)/86400; // 86400 = secsInDay
        uint32 sumSkills;
        for (uint8 s = 0; s < N_SKILLS; s++) sumSkills += skillsVec[s];
        uint256 skills = encodePlayerSkills(
            skillsVec, 
            dayOfBirth, 
            0,
            playerId, 
            birthTraits, 
            false, 
            false, 
            0, 
            0, 
            false, 
            sumSkills
        );
        return addIsSpecial(skills);
    }
    
    function createPromoPlayer(
        uint16[N_SKILLS] memory skillsVec,
        uint256 ageInSecs,
        uint8[4] memory birthTraits,
        uint256 playerId,
        uint256 targetTeamId
    ) public view returns (uint256) {
        return setTargetTeamId(
            createSpecialPlayer(skillsVec, ageInSecs, birthTraits, playerId),
            targetTeamId
        );
    }
}