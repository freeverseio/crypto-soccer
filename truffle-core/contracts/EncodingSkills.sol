pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingSkills {

    uint8 constant public PLAYERS_PER_TEAM_INIT = 18;
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public MIN_PLAYER_AGE_AT_BIRTH = 16;
    uint8 constant public MAX_PLAYER_AGE_AT_BIRTH = 32;
    uint8 constant public N_SKILLS = 5;

    // prefPosition idxs: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    uint8 constant public IDX_GK = 0;
    uint8 constant public IDX_D  = 1;
    uint8 constant public IDX_M  = 2;
    uint8 constant public IDX_F  = 3;
    uint8 constant public IDX_MD = 4;
    uint8 constant public IDX_MF = 5;
    //  Leftishness:   0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    uint8 constant public IDX_R = 1;
    uint8 constant public IDX_C = 2;
    uint8 constant public IDX_CR = 3;
    uint8 constant public IDX_L = 4;
    uint8 constant public IDX_LR = 5;
    uint8 constant public IDX_LC = 6;
    uint8 constant public IDX_LCR = 7;


    /**
     * @dev Tactics serializes a total of 71 bits = 55 + 10 + 6:
     *      lineup[11]          = 5 bit each = [playerIdxInTeam1, ..., ]
     *      extraAttack[10]     = 1 bit each, 0: normal, 1: player has extra attack duties
     *      tacticsId           = 6 bit (0 = 442, 1 = 541, ...
    **/
    function encodeTactics(uint8[11] memory lineup, bool[10] memory extraAttack, uint8 tacticsId) public pure returns (uint256) {
        require(tacticsId < 64, "tacticsId should fit in 6 bit");
        uint256 encoded = uint256(tacticsId);
        for (uint8 p = 0; p < 10; p++) {
            encoded |= uint256(extraAttack[p] ? 1 : 0) << 6 + 1 * p;
        }          
        for (uint8 p = 0; p < 11; p++) {
            require(lineup[p] < PLAYERS_PER_TEAM_MAX, "incorrect lineup entry");
            encoded |= uint256(lineup[p]) << 16 + 5 * p;
        }          
        return encoded;
    }

    function decodeTactics(uint256 tactics) public pure returns (uint8[11] memory lineup, bool[10] memory extraAttack, uint8 tacticsId) {
        require(tactics < 2**81, "tacticsId should fit in 61 bit");
        tacticsId = uint8(tactics & 63);
        tactics >>= 6;
        for (uint8 p = 0; p < 10; p++) {
            extraAttack[p] = ((tactics & 1) == 1 ? true : false); // 2^1 - 1
            tactics >>= 1;
        }          
        for (uint8 p = 0; p < 11; p++) {
            lineup[p] = uint8(tactics & 31); // 2^5 - 1
            require(lineup[p] < PLAYERS_PER_TEAM_MAX, "incorrect lineup entry");
            tactics >>= 5;
        }          
    }
    
    /**
     * @dev PlayerId and TeamId both serialize a total of 43 bits:
     *      timeZone        = 5 bits
     *      countryIdxInTZ  = 10 bits
     *      val             = 28 bits (either  (playerIdxInCountry or teamIdxInCountry)
    **/
    function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) public pure returns (uint256)
    {
        require(timeZone < 2**5, "defence out of bound");
        require(countryIdxInTZ < 2**10, "defence out of bound");
        require(val < 2**28, "defence out of bound");
        uint256 encoded  = uint256(timeZone) << 38;        // 43 - 5
        encoded         |= countryIdxInTZ << 28;  // 38 - 10
        return (encoded | val);            // 28 - 28
    }

    function decodeTZCountryAndVal(uint256 encoded) public pure returns (uint8, uint256, uint256)
    {
        // 2**14 - 1 = 31;  2**10 - 1 = 1023; 2**28 - 1 = 268435455;
        return (uint8(encoded >> 38 & 31), uint256(encoded >> 28 & 1023), uint256(encoded & 268435455));
    }

    /**
     * @dev PlayerSkills serializes a total of 145 bits:  6*14 + 4 + 3+ 3 + 43 + 1 + 1 + 3 +3
     *      5 skills                  = 5 x 14 bits
     *                                = shoot, speed, pass, defence, endurance
     *      potential                 = 4 bits (number is limited to [0,...,9])
     *      monthOfBirth              = 14 bits  (since Unix time)
     *      forwardness               = 3 bits
     *                                  GK: 0, D: 1, M: 2, F: 3, MD: 4, MF: 5
     *      leftishness               = 3 bits, in boolean triads: (L, C, R):
     *                                  0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
     *      playerId                  = 43 bits
     *      
     *      alignedLastHalf           = 1 (bool)
     *      redCardLastGame           = 1 (bool)
     *      gamesNonStopping          = 3 (0, 1, ..., 6). Finally, 7 means more than 6.
     *      injuryWeeksLeft           = 3 
    **/
    function encodePlayerSkills(
        uint16[N_SKILLS] memory skills, 
        uint256 monthOfBirth, 
        uint256 playerId, 
        uint8 potential, 
        uint8 forwardness, 
        uint8 leftishness,
        bool alignedLastHalf, 
        bool redCardLastGame, 
        uint8 gamesNonStopping, 
        uint8 injuryWeeksLeft
    )
        public
        pure
        returns (uint256 encoded)
    {
        // checks:
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
            require(skills[sk] < 2**14, "skill out of bound");
        }
        require(potential < 10, "potential out of bound");
        require(forwardness < 6, "prefPos out of bound");
        require(leftishness < 8, "prefPos out of bound");
        if (leftishness == 0) require(forwardness == 0, "leftishnes can only be zero for goalkeepers");
        require(gamesNonStopping < 8, "gamesNonStopping out of bound");
        require(monthOfBirth < 2**14, "monthOfBirthInUnixTime out of bound");
        require(playerId > 0 && playerId < 2**43, "playerId out of bound");

        // start encoding:
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
            encoded |= uint256(skills[sk]) << 256 - (sk + 1) * 14;
        }
        encoded |= monthOfBirth << 172;
        encoded |= playerId << 129;
        encoded |= uint256(potential) << 125;
        encoded |= uint256(forwardness) << 122;
        encoded |= uint256(leftishness) << 119;
        encoded |= uint256(alignedLastHalf ? 1 : 0) << 118;
        encoded |= uint256(redCardLastGame ? 1: 0 ) << 117;
        encoded |= uint256(gamesNonStopping) << 114;
        return (encoded | uint256(injuryWeeksLeft) << 111);
    }
    
    function getShoot(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 242 & 0x3fff); // 0x3fff = 2**14 - 1
    }
    
    function getSpeed(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 228 & 0x3fff);
    }

    function getPass(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 214 & 0x3fff);
    }

    function getDefence(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 200 & 0x3fff);
    }

    function getEndurance(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 186 & 0x3fff);
    }

    function getMonthOfBirth(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 172 & 0x3fff);
    }

    function getPotential(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 125 & 15);
    }

    function getForwardness(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 122 & 7);
    }

    function getLeftishness(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 119 & 7);
    }

    function getAlignedLastHalf(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 118 & 1) == 1;
    }

    function getRedCardLastGame(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 117 & 1) == 1;
    }

    function getGamesNonStopping(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 114 & 7);
    }

    function getInjuryWeeksLeft(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 111 & 7);
    }

    function getPlayerIdFromSkills(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 129 & 8796093022207); // 2**43 - 1 = 8796093022207
    }

    function getSkills(uint256 encodedSkills) public pure returns (uint256) {
        return encodedSkills >> 186;
    }

    function getSkillsVec(uint256 encodedSkills) public pure returns (uint16[5] memory skills) {
        skills[0] = uint16(getShoot(encodedSkills));
        skills[1] = uint16(getSpeed(encodedSkills));
        skills[2] = uint16(getPass(encodedSkills));
        skills[3] = uint16(getDefence(encodedSkills));
        skills[4] = uint16(getEndurance(encodedSkills));
    }

    function getSumOfSkills(uint256 encodedSkills) public pure returns (uint256) {
        return      getShoot(encodedSkills) 
                  + getSpeed(encodedSkills) 
                  + getPass(encodedSkills)
                  + getDefence(encodedSkills)
                  + getEndurance(encodedSkills);
    }



}