pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */
import "./EncodingSkills.sol";
import "./EncodingSkillsGetters.sol";
import "./EncodingSkillsSetters.sol";
import "./AssetsView.sol";

contract Privileged is AssetsView {
    
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

    // returns a value relative to 10000
    // Relative to 1, it would be = (age < 31) ? 1 - 0.02 * (age - 16) : 1 - 0.3 - 0.065 * (age - 31)
    function ageModifier(uint256 ageYears) public pure returns(uint256) {
        return (ageYears < 31) ? 10000 - 200 * (ageYears - 16) : 10000 - 3000 - 65 * (ageYears - 31);
    }

    // returns a value relative to 10000
    // relative to 1 it would be = 0.4 + potential/7.5 
    // relative to 1e4: 4000+10000*p/7.5 = (4000*7.5+10000* p)/7.5 = (4000*15+20000 * p)/15 
    function potentialModifier(uint256 potential) public pure returns(uint256) {
        return (4000 * 15 + 20000 * potential) / 15;
    }
    
    // birthTraits = [potential, forwardness, leftishness, aggressiveness]
    function createBuyNowPlayerId(uint256 playerValue, uint256 seed, uint8 forwardPos) public view returns(uint256) {
        (uint16[N_SKILLS] memory skillsVec, uint256 ageYears, uint8[4] memory birthTraits, uint256 internalPlayerId) 
            = createBuyNowPlayerIdPure(playerValue, seed, forwardPos);
        // 1 year = 31536000 sec
        return createSpecialPlayer(skillsVec, ageYears * 31536000, birthTraits, internalPlayerId);
    }
    
    function createBuyNowPlayerIdPure(
        uint256 playerValue, 
        uint256 seed, 
        uint8 forwardPos
    ) 
        public 
        pure 
        returns(uint16[N_SKILLS] memory skillsVec, uint256 ageYears, uint8[4] memory birthTraits, uint256 internalPlayerId) 
    {
        uint8 potential = uint8(seed % 10);
        seed /= 10;
        ageYears = 16 + (seed % 20);
        seed /= 20;
        uint256 avgSkills = (playerValue * 100000000)/(ageModifier(ageYears) * potentialModifier(potential));
        uint8 shirtNum;
        if (forwardPos == IDX_GK) {
            shirtNum = uint8(seed % 3);
        } else if (forwardPos == IDX_D) {
            shirtNum = 3 + uint8(seed % 5);
        } else if (forwardPos == IDX_M) {
            shirtNum = 8 + uint8(seed % 6);
        } else if (forwardPos == IDX_F) {
            shirtNum = 14 + uint8(seed % 4);
        }
        seed /= 8;
        (skillsVec, birthTraits, ) = computeSkills(seed, shirtNum);
        birthTraits[IDX_POT] = potential;
        for (uint8 sk = 0; sk < N_SKILLS; sk++) skillsVec[sk] = uint16((uint256(skillsVec[sk]) * avgSkills)/uint256(1000));
        internalPlayerId = seed % 8796093022207; // maxPlayerId (43b) = 2**43 - 1 = 8796093022207
    }
}
