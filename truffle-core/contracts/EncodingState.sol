pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingState {

    ////////////////////////////////
    /**
     * @dev PlayerState serializes a total of 169 bits:
     *  playerId                = 43 bits, offset = 0
     *  currentTeamId           = 43 bits, offset = 43
     *  currentShirtNum         =  5 bits, offset = 86
     *  prevPlayerTeamId        = 43 bits, offset = 91
     *  lastSaleBlocknum        = 35 bits, offset = 134
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
        uint256 state = playerId;
        
        
        state = setCurrentTeamId(state, currentTeamId);
        state = setCurrentShirtNum(state, currentShirtNum);
        state = setPrevPlayerTeamId(state, prevPlayerTeamId);
        state = setLastSaleBlock(state, lastSaleBlock);
        return state;
    }

    function getPlayerIdFromState(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState & (2**43-1)); 
    }
    
    function setCurrentTeamId(uint256 playerState, uint256 teamId) public pure returns (uint256) {
        require(teamId < 2**43, "currentTeamIdx out of bound");
        playerState &= ~uint256((2**43-1) << 43);
        playerState |= uint256(teamId) << 43;
        return playerState;
    }

    function getCurrentTeamIdFromPlayerState(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 43 & (2**43-1));
    }

    function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) public pure returns (uint256) {
        require(currentShirtNum < 2**5, "currentShirtNum out of bound");
        state &= ~uint256(2**5-1 << 86); 
        state |= uint256(currentShirtNum) << 86;
        return state;
    }

    function getCurrentShirtNum(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 86 & (2**5-1));
    }
    
    function setPrevPlayerTeamId(uint256 state, uint256 value) public pure returns (uint256) {
        require(value < 2**43, "prevLeagueIdx out of bound");
        state &= ~uint256(2**43-1 << 91); 
        state |= uint256(value) << 91;
        return state;
    }

    function getPrevPlayerTeamId(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 91 & (2**43-1));
    }

    function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) public pure returns (uint256) {
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        state &= ~uint256(2**35-1 << 134); // 256 - 43 - 43 - 5 - 43 - 35
        state |= uint256(lastSaleBlock) << 134;
        return state;
    }

    function getLastSaleBlock(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 134 & (2**35-1));
    }

}
