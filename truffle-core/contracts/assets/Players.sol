pragma solidity ^0.5.0;

import "./Storage.sol";

contract Players is Storage {
    uint8 constant BITS_PER_SKILL = 14;
    uint16 constant SKILL_MASK = 0x3fff; 
    uint8 constant public NUM_SKILLS = 5;

    /// this function uses the inverse of the following formula
    /// playerId = playersPerTeam * (teamId -1) + 1 + posInTeam;
    function getPlayerTeam(uint256 playerId) public view returns (uint256) {
        require(playerId != 0, "invalid player id");
        uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM; 
        require(teamId <= countTeams(), "playerId not created");
        return teamId;
    }

    function getPlayerPosInTeam(uint256 playerId) public view returns (uint256) {
        uint256 teamId = getPlayerTeam(playerId);
        return playerId - PLAYERS_PER_TEAM * (teamId - 1) - 1;
    }

    function getPlayerSkills(uint256 playerId) external view returns (uint16[NUM_SKILLS] memory) {
        uint256 teamId = getPlayerTeam(playerId);
        uint256 posInTeam = getPlayerPosInTeam(playerId);
        string memory teamName = getTeamName(teamId);
        uint256 seed = uint256(keccak256(abi.encodePacked(teamName, posInTeam)));
        return _computeSkills(seed);
    }

    /// Compute a random age between 16 and 35 
    /// @return monthOfBirth in monthUnixTime
    function _computeBirth(uint256 rnd, uint256 currentTime) internal pure returns (uint16) {
        rnd >>= BITS_PER_SKILL*NUM_SKILLS;
        uint16 seed = uint16(rnd & SKILL_MASK);
        /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
        uint16 age = 16 + (seed % 20);

        /// @dev Convert age to monthOfBirthAfterUnixEpoch.
        /// @dev I leave it this way for clarity, for the time being.
        uint years2secs = 365 * 24 * 3600;
        uint month2secs = 30 * 24 * 3600;

        return uint16((currentTime - age * years2secs) / month2secs);
    }

     /**
     * @dev Compute the pseudorandom skills, sum of the skills is 250
     * @param rnd is a random number
     * @return 5 skills
     */
    function _computeSkills(uint256 rnd) internal pure returns (uint16[NUM_SKILLS] memory) {
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