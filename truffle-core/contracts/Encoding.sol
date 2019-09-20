pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract Encoding {

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
     * @dev encoding of a total of 71 bits:
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
     * @dev encoding of a total of 43 bits:
     *      timeZone                  = 5 bits
     *      countryIdxInTZ            = 10 bits
     *      val (playerId or teamIdx) = 28 bits
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
     * @dev encoding of a total of 137 bits:  6*14 + 4 + 3+ 3 + 43
     *      5 skills                  = 5 x 14 bits
     *                                = shoot, speed, pass, defence, endurance
     *      potential                 = 4 bits (number is limited to [0,...,9])
     *      monthOfBirth              = 14 bits  (since Unix time)
     *      forwardness               = 3 bits
     *                                  GK: 0, D: 1, M: 2, F: 3, MD: 4, MF: 5
     *      leftishness               = 3 bits, in boolean triads: (L, C, R):
     *                                  0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
     *      playerId                  = 43 bits
    **/
    function encodePlayerSkills(
        uint16[N_SKILLS] memory skills, 
        uint256 monthOfBirth, 
        uint256 playerId, 
        uint8 potential, 
        uint8 forwardness, 
        uint8 leftishness
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
        return (encoded | uint256(leftishness) << 119);
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


    ////////////////////////////////
    /**
     * @dev encoding of playerState:
     *  playerId                = 43 bits
     *  currentTeamId           = 43 bits
     *  currentShirtNum         =  5 bits
     *  prevPlayerTeamId        = 43 bits
     *  lastSaleBlocknum        = 35 bits
     */
    function encodePlayerState(
        uint256 playerId,
        uint256 currentTeamId,
        uint8 currentShirtNum,
        uint256 prevPlayerTeamId,
        uint256 lastSaleBlock
    )
        public
        pure
        returns (uint256)
    {
        require(playerId > 0 && playerId < 2**43, "playerId out of bound");
        uint256 state = uint256(playerId) << 213;  // 256 - 43
        state = setCurrentTeamId(state, currentTeamId);
        state = setCurrentShirtNum(state, currentShirtNum);
        state = setPrevPlayerTeamId(state, prevPlayerTeamId);
        state = setLastSaleBlock(state, lastSaleBlock);
        return state;
    }

    function getPlayerIdFromState(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 213 & (2**43-1)); 
    }
    
    function setCurrentTeamId(uint256 playerState, uint256 teamId) public pure returns (uint256) {
        require(teamId < 2**43, "currentTeamIdx out of bound");
        playerState &= ~uint256((2**43-1) << 170); // 256 - 43 - 43
        playerState |= uint256(teamId) << 170;
        return playerState;
    }

    function getCurrentTeamId(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 170 & (2**43-1));
    }

    function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) public pure returns (uint256) {
        require(currentShirtNum < 2**5, "currentShirtNum out of bound");
        state &= ~uint256(2**5-1 << 165); // 256 - 43 - 43 - 5
        state |= uint256(currentShirtNum) << 165;
        return state;
    }

    function getCurrentShirtNum(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 165 & (2**5-1));
    }
    
    function setPrevPlayerTeamId(uint256 state, uint256 value) public pure returns (uint256) {
        require(value < 2**43, "prevLeagueIdx out of bound");
        state &= ~uint256(2**43-1 << 122); // 256 - 43 - 43 - 5 - 43
        state |= uint256(value) << 122;
        return state;
    }

    function getPrevPlayerTeamId(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 122 & (2**43-1));
    }

    function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) public pure returns (uint256) {
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        state &= ~uint256(2**35-1 << 87); // 256 - 43 - 43 - 5 - 43 - 35
        state |= uint256(lastSaleBlock) << 87;
        return state;
    }

    function getLastSaleBlock(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 87 & (2**35-1));
    }

    /// @dev Sets the number at a given index in a serialized uint256
    // function setNumAtIndex(uint value, uint serialized, uint8 index, uint bits)
    //     internal
    //     pure
    //     returns(uint)
    // {
    //     uint maxnum = 1<<bits; // 2**bits
    //     require(value < maxnum, "Value too large to fit in available space");
    //     uint b = bits*index;
    //     uint mask = (1 << bits)-1; // (2**bits)-1
    //     serialized &= ~(mask << b); // clear all bits at index
    //     return serialized + (value << b);
    // }

}