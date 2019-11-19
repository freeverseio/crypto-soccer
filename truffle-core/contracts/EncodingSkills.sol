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
    uint8 constant public NO_SUBST = 11;

    // Birth Traits: potential, forwardness, leftishness, aggressiveness
    uint8 constant private IDX_POT = 0;
    uint8 constant private IDX_FWD = 1;
    uint8 constant private IDX_LEF = 2;
    uint8 constant private IDX_AGG = 3;
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
     * @dev Tactics serializes a total of 110 bits = 3 * 4 + 3 * 4 + 14*5 + 10 + 6:
     *      substitutions[3]    = 4 bit each = [3 different nums from 0 to 10], with 11 = no subs
     *      lineup[14]          = 5 bit each = [playerIdxInTeam1, ..., ]
     *      extraAttack[10]     = 1 bit each, 0: normal, 1: player has extra attack duties
     *      tacticsId           = 6 bit (0 = 442, 1 = 541, ...
    **/
    function encodeTactics(
        uint8[3] memory substitutions, 
        uint8[3] memory subsRounds, 
        uint8[14] memory lineup, 
        bool[10] memory extraAttack, 
        uint8 tacticsId
    ) 
        public 
        pure 
        returns (uint256) 
    {
        require(tacticsId < 64, "tacticsId should fit in 64 bit");
        uint256 encoded = uint256(tacticsId);
        for (uint8 p = 0; p < 10; p++) {
            encoded |= uint256(extraAttack[p] ? 1 : 0) << 6 + p;
        }          
        for (uint8 p = 0; p < 11; p++) {
            require(lineup[p] < PLAYERS_PER_TEAM_MAX, "incorrect lineup entry");
            encoded |= uint256(lineup[p]) << 16 + 5 * p;
        }          
        for (uint8 p = 0; p < 3; p++) {
            require(substitutions[p] < 12, "incorrect lineup entry");
            require(subsRounds[p] < 12, "incorrect round");
            // requirement: if there is no subst at "i", lineup[i + 11] = 25 + p (so that all lineups are different, and sortable)
            if (substitutions[p] == NO_SUBST) {
                require(lineup[p + 11] == 25 + p, "incorrect lineup entry for no substituted player");
            }
            encoded |= uint256(lineup[p + 11]) << 16 + 5 * (p + 11);
            encoded |= uint256(substitutions[p]) << 86 + 4 * p;
            encoded |= uint256(subsRounds[p]) << 98 + 4 * p;
        }          
        return encoded;
    }

    function decodeTactics(uint256 tactics) public pure returns (
        uint8[3] memory substitutions, 
        uint8[3] memory subsRounds, 
        uint8[14] memory lineup, 
        bool[10] memory extraAttack, 
        uint8 tacticsId
    ) {
        require(tactics < 2**110, "tacticsId should fit in 98 bit");
        tacticsId = uint8(tactics & 63);
        for (uint8 p = 0; p < 10; p++) {
            extraAttack[p] = (((tactics >> (6 + p)) & 1) == 1 ? true : false); // 2^1 - 1
        }          
        for (uint8 p = 0; p < 3; p++) {
            substitutions[p] = uint8((tactics >> (86 + 4 * p)) & 15); // 2^4 - 1
            // require(substitutions[p] < 12, "incorrect substitutions entry"); // 11 is used as "no substitution"
        }          
        for (uint8 p = 0; p < 14; p++) {
            lineup[p] = uint8((tactics >> (16 + 5 * p)) & 31); // 2^5 - 1
            // if ((p > 10) && (substitutions[p - 11] == NO_SUBST)) require(lineup[p] == 14 + p, "incorrect lineup entry for no substituted player");
            // else require(lineup[p] < PLAYERS_PER_TEAM_MAX, "incorrect lineup entry");
        }          
        for (uint8 p = 0; p < 3; p++) {
            subsRounds[p] = uint8(tactics >> (98 + 4 * p) & 15); // 2^4 - 1
            // require(subsRounds[p] < 12, "incorrect round entry");
        }          
    }
    
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
     *      substitutedDuringLastHalf = 1b (bool) 
     *      sumSkills                 = 19b (must equal sum(skills), of if each is 16b, this can be at most 5x16b => use 19b)
     *      isSpecialPlayer           = 1b (set at the left-most bit, 255)
    **/
    function encodePlayerSkills(
        uint16[N_SKILLS] memory skills, 
        uint256 dayOfBirth, 
        uint256 playerId, 
        uint8[4] memory birthTraits,
        bool alignedEndOfLastHalf, 
        bool redCardLastGame, 
        uint8 gamesNonStopping, 
        uint8 injuryWeeksLeft,
        bool substitutedLastHalf,
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
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
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
        encoded |= uint256(substitutedLastHalf ? 1 : 0) << 160;
        return (encoded | uint256(sumSkills) << 161);
    }
    
    function getShoot(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills & 65535); // 65535 = 2**16 - 1
    }
    
    function getSpeed(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 16 & 65535);
    }

    function getPass(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 32 & 65535);
    }

    function getDefence(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 48 & 65535);
    }

    function getEndurance(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 64 & 65535);
    }

    function getBirthDay(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 80 & 65535);
    }

    function getPlayerIdFromSkills(uint256 encodedSkills) public pure returns (uint256) {
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

    function getAlignedEndOfLastHalf(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 152 & 1) == 1;
    }

    function getRedCardLastGame(uint256 encodedSkills) public pure returns (bool) {
        return (encodedSkills >> 153 & 1) == 1;
    }

    function getGamesNonStopping(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 154 & 7);
    }

    function getInjuryWeeksLeft(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 157 & 7);
    }

    function getSubstitutedLastHalf(uint256 encodedSkills) public pure returns (bool) {
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

    function setTargetTeamId(uint256 encodedSkills, uint256 targetTeamId) public pure returns (uint256) {
        return (encodedSkills & ~(uint256(2**43-1) << 180)) | (targetTeamId << 180);
    }

    function getTargetTeamId(uint256 encodedSkills) public pure returns (uint256) {
        return (encodedSkills >> 180) & (2**43-1);
    }

    function getGeneration(uint256 encodedSkills) public pure returns (uint256) {
        return (encodedSkills >> 223) & 63;
    }
}