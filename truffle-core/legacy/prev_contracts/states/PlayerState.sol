pragma solidity >=0.4.21 <0.6.0;

/// @title the state of a player
contract PlayerState {
    function isValidPlayerState(uint256 state) public pure returns (bool) {
        return getPlayerId(state) != 0;
    }

    /**
     * @dev encoding:
     * 5x14bits
     * skills                  = 5x14 bits
     * monthOfBirthInUnixTime  = 14 bits
     * playerIdx               = 28 bits
     * currentTeamIdx          = 28 bits
     * currentShirtNum         =  5 bits
     * prevLeagueIdx           = 25 bits
     * prevTeamPosInLeague     =  8 bits
     * prevShirtNumInLeague    =  5 bits // TODO: remove: unused
     * lastSaleBlocknum        = 35 bits
     * available               = 40 bits
     */
    function playerStateCreate(
        uint256 defence,
        uint256 speed,
        uint256 pass,
        uint256 shoot,
        uint256 endurance,
        uint256 monthOfBirthInUnixTime,
        uint256 playerId,
        uint256 currentTeamId,
        uint256 currentShirtNum,
        uint256 prevLeagueId,
        uint256 prevTeamPosInLeague,
        uint256 prevShirtNumInLeague,
        uint256 lastSaleBlock
    )
        public
        pure
        returns (uint256 state)
    {
        require(defence < 2**14, "defence out of bound");
        require(speed < 2**14, "defence out of bound");
        require(pass < 2**14, "defence out of bound");
        require(shoot < 2**14, "defence out of bound");
        require(endurance < 2**14, "defence out of bound");
        require(monthOfBirthInUnixTime < 2**14, "monthOfBirthInUnixTime out of bound");
        require(playerId > 0 && playerId < 2**28, "playerId out of bound");
        require(prevShirtNumInLeague < 2**5, "prevShirtNumInLeague out of bound");
        state = uint256(defence) << 242;
        state |= uint256(speed) << 228;
        state |= uint256(pass) << 214;
        state |= uint256(shoot) << 200;
        state |= uint256(endurance) << 186;
        state |= uint256(monthOfBirthInUnixTime) << 172;
        state |= uint256(playerId) << 144;
        state = setCurrentTeamId(state, currentTeamId);
        state = setCurrentShirtNum(state, currentShirtNum);
        state = setPrevLeagueId(state, prevLeagueId);
        state = setPrevTeamPosInLeague(state, prevTeamPosInLeague);
        state |= uint256(prevShirtNumInLeague) << 73;
        state = setLastSaleBlock(state, lastSaleBlock);
    }

    function setPrevLeagueId(uint256 state, uint256 value) public pure returns (uint256) {
        require(value < 2**25, "prevLeagueIdx out of bound");
        state &= ~uint256(2**25-1 << 86);
        state |= uint256(value) << 86;
        return state;
    }

    function setPrevTeamPosInLeague(uint256 state, uint256 value) public pure returns (uint256) {
        require(value < 2**8, "prevTeamPosInLeague out of bound");
        state &= ~uint256(2**8-1 << 78);
        state |= uint256(value) << 78;
        return state;
    }
    
    /// increase the skills of delta
    function playerStateEvolve(uint256 playerState, uint16 delta) public pure returns (uint256 evolvedState) {
        require(isValidPlayerState(playerState), "invalid player playerState");
        return playerState;
        uint256 defence = getDefence(playerState) + delta;
        uint256 speed = getSpeed(playerState) + delta;
        uint256 pass = getPass(playerState) + delta;
        uint256 shoot = getShoot(playerState) + delta;
        uint256 endurance = getEndurance(playerState) + delta;
        require(defence < 2**14, "defence out of bound");
        require(speed < 2**14, "speed out of bound");
        require(pass < 2**14, "pass out of bound");
        require(shoot < 2**14, "shoot out of bound");
        require(endurance < 2**14, "endurance out of bound");
        evolvedState = playerState;
        evolvedState = evolvedState & (uint256(-1) ^ (0x3fff << 242)) | uint256(defence) << 242;
        evolvedState = evolvedState & (uint256(-1) ^ (0x3fff << 228)) | uint256(speed) << 228;
        evolvedState = evolvedState & (uint256(-1) ^ (0x3fff << 214)) | uint256(pass) << 214;
        evolvedState = evolvedState & (uint256(-1) ^ (0x3fff << 200)) | uint256(shoot) << 200;
        evolvedState = evolvedState & (uint256(-1) ^ (0x3fff << 186)) | uint256(endurance) << 186;
    }

    function setCurrentTeamId(uint256 playerState, uint256 teamId) public pure returns (uint256) {
        require(teamId < 2**28, "currentTeamIdx out of bound");
        playerState &= ~uint256(2**28-1 << 116);
        playerState |= uint256(teamId) << 116;
        return playerState;
    }

    function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) public pure returns (uint256) {
        require(currentShirtNum < 2**5, "currentShirtNum out of bound");
        state &= ~uint256(2**5-1 << 111);
        state |= uint256(currentShirtNum) << 111;
        return state;
    }

    function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) public pure returns (uint256) {
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        state &= ~uint256(2*35-1 << 38);
        state |= uint256(lastSaleBlock) << 38;
        return state;
    }

    function getLastSaleBlock(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 38 & 0x7ffffffff);
    }

    function getPrevShirtNumInLeague(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 73 & 0x1f);
    }

    function getPrevTeamPosInLeague(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 78 & 0xff);
    }

    function getPrevLeagueId(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 86 & 0x1ffffff);
    }

    function getCurrentShirtNum(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 111 & 0x1f);
    }

    function getCurrentTeamId(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 116 & 0xfffffff);
    }

    function getPlayerId(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 144 & 0xfffffff);
    }

    function getMonthOfBirthInUnixTime(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 172 & 0x3fff);
    }

    function getDefence(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 242 & 0x3fff);
    }
    
    function getSpeed(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 228 & 0x3fff);
    }

    function getPass(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 214 & 0x3fff);
    }

    function getShoot(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 200 & 0x3fff);
    }

    function getEndurance(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint256(playerState >> 186 & 0x3fff);
    }

    function getSkills(uint256 playerState) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player state");
        return playerState >> 186;
    }

    function getSkillsVec(uint256 playerState) public pure returns (uint16[5] memory skills) {
        require(isValidPlayerState(playerState), "invalid player state");
        skills[0] = uint16(getDefence(playerState));
        skills[1] = uint16(getSpeed(playerState));
        skills[2] = uint16(getPass(playerState));
        skills[3] = uint16(getShoot(playerState));
        skills[4] = uint16(getEndurance(playerState));
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