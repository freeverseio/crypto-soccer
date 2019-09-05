pragma solidity >=0.4.21 <0.6.0;

/// @title library of functions to encode-decode player states and skills, etc.
contract PlayerState {

    uint8 constant public MIN_PLAYER_AGE_AT_BIRTH = 16;
    uint8 constant public MAX_PLAYER_AGE_AT_BIRTH = 32;

    /**
     * @dev encoding a total of 43 bits:
     * timeZone                  = 5 bits
     * countryIdxInTZ            = 10 bits
     * val (playerId or teamIdx) = 28 bits
    **/
    function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) public pure returns (uint256)
    {
        require(timeZone < 2**5, "defence out of bound");
        require(countryIdxInTZ < 2**10, "defence out of bound");
        require(val < 2**28, "defence out of bound");
        uint256 encoded  = uint256(timeZone) << 251;        // 256 - 5
        encoded         |= uint256(countryIdxInTZ) << 241;  // 251 - 10
        return (encoded |= uint256(val) << 213);            // 241 - 28
    }

    function decodeTZCountryAndVal(uint256 encoded) public pure returns (uint8, uint256, uint256)
    {
        // 2**14 - 1 = 31;  2**10 - 1 = 1023; 2**28 - 1 = 268435455;
        return (uint8(encoded >> 251 & 31), uint256(encoded >> 241 & 1023), uint256(encoded >> 213 & 268435455));
    }

    /**
     * @dev encoding a total of 62 bits:
     * 5 skills                  = 5 x 14 bits
     * monthOfBirthInUnixTime    = 14 bits
     * playerId                  = 43 bits
    **/
    // TODO: avoid doing the uint256 for those variables that already are uint256
    function encodePlayerSkills(
        uint256 defence,
        uint256 speed,
        uint256 pass,
        uint256 shoot,
        uint256 endurance,
        uint256 monthOfBirthInUnixTime,
        uint256 playerId
    )
        public
        pure
        returns (uint256)
    {
        require(defence < 2**14, "defence out of bound");
        require(speed < 2**14, "defence out of bound");
        require(pass < 2**14, "defence out of bound");
        require(shoot < 2**14, "defence out of bound");
        require(endurance < 2**14, "defence out of bound");
        require(monthOfBirthInUnixTime < 2**14, "monthOfBirthInUnixTime out of bound");
        require(playerId > 0 && playerId < 2**43, "playerId out of bound");
        uint256 skills = uint256(defence) << 242;
        skills |= uint256(speed) << 228;
        skills |= uint256(pass) << 214;
        skills |= uint256(shoot) << 200;
        skills |= uint256(endurance) << 186;
        skills |= uint256(monthOfBirthInUnixTime) << 172;
        return (skills |= uint256(playerId) << 129);
    }
    
    function getDefence(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 242 & 0x3fff); // 0x3fff = 2**14 - 1
    }
    
    function getSpeed(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 228 & 0x3fff);
    }

    function getPass(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 214 & 0x3fff);
    }

    function getShoot(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 200 & 0x3fff);
    }

    function getEndurance(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 186 & 0x3fff);
    }

    function getMonthOfBirthInUnixTime(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 172 & 0x3fff);
    }

    function getPlayerIdFromSkills(uint256 encodedSkills) public pure returns (uint256) {
        return uint256(encodedSkills >> 129 & 8796093022207); // 2**43 - 1 = 8796093022207
    }

    function getSkills(uint256 encodedSkills) public pure returns (uint256) {
        return encodedSkills >> 186;
    }

    function getSkillsVec(uint256 encodedSkills) public pure returns (uint16[5] memory skills) {
        skills[0] = uint16(getDefence(encodedSkills));
        skills[1] = uint16(getSpeed(encodedSkills));
        skills[2] = uint16(getPass(encodedSkills));
        skills[3] = uint16(getShoot(encodedSkills));
        skills[4] = uint16(getEndurance(encodedSkills));
    }


    ////////////////////////////////
    /**
     * @dev encoding of playerState:
     * playerId                = 43 bits
     * currentTeamId           = 43 bits
     * currentShirtNum         =  5 bits
     * prevPlayerTeamId        = 43 bits
     * lastSaleBlocknum        = 35 bits
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
    function setNumAtIndex(uint value, uint serialized, uint8 index, uint bits)
        internal
        pure
        returns(uint)
    {
        uint maxnum = 1<<bits; // 2**bits
        require(value < maxnum, "Value too large to fit in available space");
        uint b = bits*index;
        uint mask = (1 << bits)-1; // (2**bits)-1
        serialized &= ~(mask << b); // clear all bits at index
        return serialized + (value << b);
    }

}