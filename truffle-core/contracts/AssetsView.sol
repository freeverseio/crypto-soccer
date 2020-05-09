pragma solidity >=0.5.12 <=0.6.3;

import "./EncodingSkills.sol";
import "./EncodingIDs.sol";
import "./EncodingState.sol";
import "./Storage.sol";
import "./AssetsLib.sol";

/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with tz = 0 being null.
 */

contract AssetsView is AssetsLib, EncodingSkills, EncodingState {
    
    function getPlayerSkillsAtBirth(uint256 playerId) public view returns (uint256) {
        if (getIsSpecial(playerId)) return playerId;
        if (!wasPlayerCreatedVirtually(playerId)) return 0;
        (uint256 teamId, uint256 playerCreationDay, uint8 shirtNum) = getTeamIdCreationDayAndShirtNum(playerId);
        (uint256 dayOfBirth, uint8 potential) = computeBirthDayAndPotential(teamId, playerCreationDay, shirtNum);
        (uint16[N_SKILLS] memory skills, uint8[4] memory birthTraits, uint32 sumSkills) = computeSkills(teamId, shirtNum, potential);
        return encodePlayerSkills(skills, dayOfBirth, 0, playerId, birthTraits, false, false, 0, 0, false, sumSkills);
        
    }

    function getTeamIdCreationDayAndShirtNum(uint256 playerId) public view returns(uint256 teamId, uint256 creationDay, uint8 shirtNum) {
        (uint8 tz, uint256 countryIdxInTZ, uint256 playerIdxInCountry) = decodeTZCountryAndVal(playerId);
        uint256 teamIdxInCountry = playerIdxInCountry / PLAYERS_PER_TEAM_INIT;
        uint256 divisionIdx = teamIdxInCountry / TEAMS_PER_DIVISION;
        uint256 divisionId = encodeTZCountryAndVal(tz, countryIdxInTZ, divisionIdx);
        teamId = encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
        creationDay = gameDeployDay + divisionIdToRound[divisionId] * DAYS_PER_ROUND;
        shirtNum = uint8(playerIdxInCountry % PLAYERS_PER_TEAM_INIT);
    }

    /// Compute a random age between 16 and 35.999, with random potential in [0,...7]
    /// @param playerCreationDay - days since unix epoch where the player was created as part of teams of a division
    /// @return dayOfBirth - days since unix epoch
    function computeBirthDayAndPotential(uint256 teamId, uint256 playerCreationDay, uint8 shirtNum) public pure returns (uint16 dayOfBirth, uint8 potential) {
        // generate a DNA that is unique to pairs of players in the universe.
        // each team has different DNAs, but within a same team, shirts = 0,1 have the same, shirts = 2,3 have the same...etc
        uint256 dna = uint256(keccak256(abi.encode(teamId, shirtNum/2)));
        // Generate pairs of potentials such that each is in [0,...,7] and the sum is 7, so average is 3.5
        potential = (shirtNum % 2 == 0) ? uint8(dna % (MAX_POTENTIAL_AT_BIRTH+1)) : MAX_POTENTIAL_AT_BIRTH - uint8(dna % (MAX_POTENTIAL_AT_BIRTH+1));
        // generate a different dna for each member of the pair by bit-shifting dna differently
        dna >>= (1 +(shirtNum % 2));
        // Increase potential average to 4.25 = 3.5 + 0.75 
        if ((potential < 7) && (dna % 4) != 0) potential += 1;
        // Compute days in range [16,36]
        uint256 ageInDays = 5840 + (dna % 7300);  // 5840 = 16y, 7300 = 20y
        // ensure that good potential players are not above 31,
        // by subtracting what is left to reach 31, plus a random between 0 and 2 years
        dna >>= 12;
        if (potential > 5 && ageInDays > 11315) ageInDays -= (ageInDays - 11315) + (dna % 730); // 11315 = 31y, 730 = 2y.
        dayOfBirth = uint16(playerCreationDay - ageInDays / INGAMETIME_VS_REALTIME); 
    }
    
    /// Compute the pseudorandom skills, sum of the skills is 5K (1K each skill on average)
    /// skills have currently, 16bits each, and there are 5 of them
    /// potential is a number between 0 and 9 => takes 4 bit
    /// 0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    /// @return uint16[N_SKILLS] skills, uint8 potential, uint8 forwardness, uint8 leftishness
    function computeSkills(uint256 teamId, uint8 shirtNum, uint8 potential) public pure returns (uint16[N_SKILLS] memory, uint8[4] memory, uint32) {
        uint16[5] memory skills;
        uint256[N_SKILLS] memory correctFactor;
        uint256 dna = uint256(keccak256(abi.encode(teamId, shirtNum)));
        uint8 forwardness;
        uint8 leftishness;
        uint8 aggressiveness = uint8(dna % 4);
        dna >>= 2; // log2(4) = 2
        // correctFactor/10 increases a particular skill depending on player's forwardness
        if (shirtNum < 2) {
            // 2 GoalKeepers:
            correctFactor[SK_SHO] = 14;
            correctFactor[SK_PAS] = 6;
            forwardness = IDX_GK;
            leftishness = 0;
        } else if (shirtNum < 7) {
            // 5 Defenders
            correctFactor[SK_SHO] = 4;
            correctFactor[SK_DEF] = 16;
            forwardness = IDX_D;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 10) {
            // 3 Pure Midfielders
            correctFactor[SK_PAS] = 16;
            forwardness = IDX_M;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 12) {
            // 2 Defensive Midfielders
            correctFactor[SK_PAS] = 13;
            correctFactor[SK_SHO] = 7;
            forwardness = IDX_MD;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 14) {
            // 2 Attachking Midfielders
            correctFactor[SK_PAS] = 13;
            correctFactor[SK_DEF] = 7;
            forwardness = IDX_MF;
            leftishness = uint8(1+ (dna % 7));
        } else if (shirtNum < 16) {
            // 2 Forwards that play center-left
            correctFactor[SK_SHO] = 16;
            correctFactor[SK_DEF] = 5;
            forwardness = IDX_F;
            leftishness = 6;
        } else {
            // 2 Forwards that play center-right
            correctFactor[SK_SHO] = 16;
            correctFactor[SK_DEF] = 5;
            forwardness = IDX_F;
            leftishness = 3;
        }
        dna >>= 3; // log2(7) = 2.9 => ceil = 3                      

        /// Compute initial skills, as a random with [0, 49] 
        /// ...apply correction factor depending on preferred pos,
        uint16 excess;
        for (uint8 i = 0; i < N_SKILLS; i++) {
            if (correctFactor[i] == 0) {
                skills[i] = uint16(dna % 800);
            } else {
                skills[i] = uint16(((dna % 800) * correctFactor[i])/10);
            }
            excess += skills[i];
            dna >>= 10; // los2(1000) -> ceil
        }
        // at this point, excess is at most, last two cases: (1.6+0.7+3)*800 = 4240, so 5000-excess is safe
        // and for GKS: (2+ 0.6 + 3)*800 = 4480, so 5000-excess is safe.
        uint16 delta;
        delta = (5000 - excess) / N_SKILLS;
        for (uint8 i = 0; i < N_SKILLS; i++) skills[i] = skills[i] + delta;
        // note: final sum of skills = excess + N_SKILLS * delta;
        return (skills, [potential, forwardness, leftishness, aggressiveness], uint32(excess + N_SKILLS * delta));
    }


    function secsToDays(uint256 secs) internal pure returns (uint256) {
        return secs / 86400;  // 86400 = 3600 * 24
    }

    function countCountries(uint8 tz) public view returns (uint256){
        return tzToNCountries[tz];
    }
    
    // TODO: remove from this contract, expose as interface for users
    function daysToSecs(uint256 dayz) internal pure returns (uint256) {
        return dayz * 86400; // 86400 = 3600 * 24 * 365
    }

    // TODO: remove from this contract, expose as interface for users
    function getPlayerAgeInDays(uint256 playerId) public view returns (uint256) {
        return secsToDays(INGAMETIME_VS_REALTIME * (now - daysToSecs(getBirthDay(getPlayerSkillsAtBirth(playerId)))));
    }
    
}
