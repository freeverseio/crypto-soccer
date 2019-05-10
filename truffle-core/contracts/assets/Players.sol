pragma solidity ^0.5.0;

import "./Storage.sol";

contract Players is Storage {
    uint8 constant BITS_PER_SKILL = 14;
    uint16 constant SKILL_MASK = 0x3fff; 

    /// this function uses the inverse of the following formula
    /// playerId = playersPerTeam * (teamId -1) + 1 + posInTeam;
    function getPlayerTeam(uint256 playerId) external view returns (uint256) {
        require(playerId != 0, "invalid player id");
        uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM; 
        require(teamId <= countTeams(), "playerId not created");
        return teamId;
    }

     /**
     * @dev Compute the pseudorandom skills, sum of the skills is 250
     * @param rnd is a random number
     * @return 5 skills
     */
    function _computeSkills(uint256 rnd) internal pure returns (uint16[5] memory) {
        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(rnd & SKILL_MASK);
            rnd >>= BITS_PER_SKILL;
        }

        /// The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < 5; i++) {
            skills[i] = skills[i] % 50;
            excess += skills[i];
        }

        /// At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        uint16 delta = (250 - excess) / 5;
        for (uint8 i = 0; i < 5; i++) 
            skills[i] = skills[i] + delta;

        uint16 remainder = (250 - excess) % 5;
        for (uint8 i = 0 ; i < remainder ; i++)
            skills[i]++;

        return skills;
    }

    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(arg)));
    }
}