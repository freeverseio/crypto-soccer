pragma solidity >=0.5.12 <=0.6.3;

import "./Constants.sol";

/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkills is Constants {

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
    function encodePlayerSkills(
        uint16[N_SKILLS] memory skills, 
        uint256 dayOfBirth, 
        uint8 generation,
        uint256 playerId, 
        uint8[4] memory birthTraits,
        bool alignedEndOfLastHalf, 
        bool redCardLastGame, 
        uint8 gamesNonStopping, 
        uint8 injuryWeeksLeft,
        bool substitutedFirstHalf,
        uint32 sumSkills
    )
        public
        pure
        returns (uint256 encoded)
    {
        // checks:
        require(birthTraits[IDX_POT] < 10, "potential out of bound");
        require(birthTraits[IDX_FWD] < 6, "forwardness out of bound");
        require(birthTraits[IDX_LEF] < 8, "lefitshness out of bound");
        require(birthTraits[IDX_AGG] < 8, "aggressiveness out of bound");
        if (birthTraits[IDX_LEF] == 0) require(birthTraits[IDX_FWD] == 0, "leftishnes can only be zero for goalkeepers");
        require(gamesNonStopping < 8, "gamesNonStopping out of bound");
        require(dayOfBirth < 2**16, "dayOfBirthInUnixTime out of bound");
        require(playerId > 0 && playerId < 2**43, "playerId out of bound");

        // start encoding:
        for (uint32 sk = 0; sk < N_SKILLS; sk++) {
            encoded |= uint256(skills[sk]) << 16 * sk;
        }
        encoded |= dayOfBirth << 80;
        encoded |= playerId << 96;
        encoded |= uint256(birthTraits[IDX_POT]) << 139;
        encoded |= uint256(birthTraits[IDX_FWD]) << 143;
        encoded |= uint256(birthTraits[IDX_LEF]) << 146;
        encoded |= uint256(birthTraits[IDX_AGG]) << 149;
        encoded |= uint256(alignedEndOfLastHalf ? 1 : 0) << 152;
        encoded |= uint256(redCardLastGame ? 1 : 0) << 153;
        encoded |= uint256(gamesNonStopping) << 154;
        encoded |= uint256(injuryWeeksLeft) << 157;
        encoded |= uint256(substitutedFirstHalf ? 1 : 0) << 160;
        encoded |= uint256(sumSkills) << 161;
        return (encoded | uint256(generation) << 223);
    }

}
