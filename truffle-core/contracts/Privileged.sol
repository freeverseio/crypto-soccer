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
        uint256 playerId,
        uint256 nowInSecs
    ) public pure returns (uint256) {
        uint256 dayOfBirth = (nowInSecs - ageInSecs/7)/86400; // 86400 = secsInDay
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
    
    // returns a value relative to 10000
    // Relative to 1, it would be = age < 31) ? 1.15 - 0.013 * (age - 16) : 1.15 - 0.013*15 - 0.05 * (age - 31)
    function ageModifier(uint256 ageYears) public pure returns(uint256) {
        return (ageYears < 31) ? 11500 - 130 * (ageYears - 16) : 9550 - 500 * (ageYears - 31);
    }

    // returns a value relative to 10000
    // relative to 1 it would be = 0.85 + potential/30
    // relative to 1e4: 8500+10000*p/30 = (8500*30+10000* p)/30
    function potentialModifier(uint256 potential) public pure returns(uint256) {
        return (8500 * 30 + 10000 * potential) / 30;
    }
    
    function computeAvgSkills(uint256 playerValue, uint256 ageYears, uint8 potential) public pure returns (uint256) {
        return (playerValue * 100000000)/(ageModifier(ageYears) * potentialModifier(potential));
    }
    
    function createBuyNowPlayerIdPure(
        uint256 playerValue, 
        uint256 seed, 
        uint8 forwardPos,
        uint8 tz,
        uint256 countryIdxInTZ
    ) 
        public 
        pure 
        returns(uint16[N_SKILLS] memory skillsVec, uint256 ageYears, uint8[4] memory birthTraits, uint256 internalPlayerId) 
    {
        uint8 potential = uint8(seed % 10);
        seed /= 10;
        ageYears = 16 + (seed % 20);
        seed /= 20;
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
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
            skillsVec[sk] = uint16(
                (uint256(skillsVec[sk]) * computeAvgSkills(playerValue, ageYears, potential))/uint256(1000)
            );
        }
        internalPlayerId = encodeTZCountryAndVal(tz, countryIdxInTZ, seed % 268435455); // maxPlayerIdxInCountry (28b) = 2**28 - 1 = 268435455
    }

    // birthTraits = [potential, forwardness, leftishness, aggressiveness]
    function createBuyNowPlayerId(
        uint256 playerValue, 
        uint256 seed, 
        uint8 forwardPos,
        uint256 epochInDays,
        uint8 tz,
        uint256 countryIdxInTZ
    ) 
        public 
        pure 
        returns
    (
        uint256 playerId,
        uint16[N_SKILLS] memory skillsVec, 
        uint16 dayOfBirth, 
        uint8[4] memory birthTraits, 
        uint256 internalPlayerId
    )
    {
        uint256 ageYears;
        (skillsVec, ageYears, birthTraits, internalPlayerId) = createBuyNowPlayerIdPure(playerValue, seed, forwardPos, tz, countryIdxInTZ);
        // 1 year = 31536000 sec
        playerId = createSpecialPlayer(skillsVec, ageYears * 31536000, birthTraits, internalPlayerId, epochInDays*24*3600);
        dayOfBirth = uint16(getBirthDay(playerId));
    }
    
    function createBuyNowPlayerIdBatch(
        uint256 playerValue, 
        uint256 seed, 
        uint8[4] memory nPlayersPerForwardPos,
        uint256 epochInDays,
        uint8 tz,
        uint256 countryIdxInTZ
    ) 
        public 
        pure 
        returns
    (
        uint256[] memory playerIdArray,
        uint16[N_SKILLS][] memory skillsVecArray, 
        uint16[] memory dayOfBirthArray, 
        uint8[4][] memory birthTraitsArray, 
        uint256[] memory internalPlayerIdArray
    )
    {
        uint16 nPlayers;
        for (uint8 pos = 0; pos < 4; pos++) { nPlayers += nPlayersPerForwardPos[pos]; }

        playerIdArray = new uint256[](nPlayers);
        skillsVecArray = new uint16[N_SKILLS][](nPlayers);
        dayOfBirthArray = new uint16[](nPlayers);
        birthTraitsArray = new uint8[4][](nPlayers);
        internalPlayerIdArray = new uint256[](nPlayers);

        uint16 counter;
        for (uint8 pos = 0; pos < 4; pos++) { 
            for (uint16 n = 0; n < nPlayersPerForwardPos[pos]; n++) {
                seed = uint256(keccak256(abi.encode(seed, n)));
                (playerIdArray[counter], skillsVecArray[counter], dayOfBirthArray[counter], birthTraitsArray[counter], internalPlayerIdArray[counter]) =
                    createBuyNowPlayerId(playerValue, seed, pos, epochInDays, tz, countryIdxInTZ);
                counter++;
            }
        }
    }
    
    function getTZandCountryIdxFromPlayerId(uint256 playerId) public pure returns (uint8 tz, uint256 countryIdxInTZ) {
        (tz, countryIdxInTZ, ) = decodeTZCountryAndVal(getInternalPlayerId(playerId));
    } 
}
