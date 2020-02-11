pragma solidity >=0.5.12 <0.6.2;

import "./EncodingTactics.sol";
import "./EncodingSkills.sol";
import "./EncodingSkillsSetters.sol";

contract EngineApplyBoosters is EncodingSkillsSetters, EncodingSkills, EncodingTactics  {

    // skills order: shoot, speed, pass, defence, endurance
    function applyItemBoost(uint256[PLAYERS_PER_TEAM_MAX] memory linedUpSkills, uint256 tactics) public pure returns(uint256[PLAYERS_PER_TEAM_MAX] memory) {
        ( , uint16 itemId, uint32 boost) = getItemsData(tactics);
        if (itemId == 0) return linedUpSkills;
        uint8[N_SKILLS+1] memory skillsBoost = decodeBoosts(boost);
        for (uint8 p = 0; p < 14; p++) {
            uint256 skills = linedUpSkills[p];
            skills = setShoot(
                skills, 
                (getShoot(skills) * (100 + skillsBoost[0])) / 100
            );
            skills = setSpeed(
                skills, 
                (getSpeed(skills) * (100 + skillsBoost[1])) / 100
            );
            skills = setPass(
                skills, 
                (getPass(skills) * (100 + skillsBoost[2])) / 100
            );
            skills = setDefence(
                skills, 
                (getDefence(skills) * (100 + skillsBoost[3])) / 100
            );
            linedUpSkills[p] = setEndurance(
                skills, 
                (getEndurance(skills) * (100 + skillsBoost[4])) / 100
            );
        }
        return linedUpSkills;
    }

}

