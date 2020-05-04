pragma solidity >=0.5.12 <=0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingState {

    ////////////////////////////////
    /**
     * @dev PlayerState serializes a total of 169 bits:
     *  currentTeamId           = 43 bits, offset = 0
     *  currentShirtNum         =  5 bits, offset = 43
     *  prevPlayerTeamId        = 43 bits, offset = 48
     *  lastSaleBlocknum        = 35 bits, offset = 91
     */
    function encodePlayerState(
        uint256 currentTeamId,
        uint8 currentShirtNum,
        uint256 prevPlayerTeamId,
        uint256 lastSaleBlock
    )
        public
        pure
        returns (uint256)
    {
        uint256 state = setCurrentTeamId(0, currentTeamId);
        state = setCurrentShirtNum(state, currentShirtNum);
        state = setPrevPlayerTeamId(state, prevPlayerTeamId);
        state = setLastSaleBlock(state, lastSaleBlock);
        return state;
    }

    function setCurrentTeamId(uint256 playerState, uint256 teamId) public pure returns (uint256) {
        require(teamId < 2**43, "currentTeamIdx out of bound");
        playerState &= ~uint256((2**43-1));
        playerState |= teamId;
        return playerState;
    }

    function getCurrentTeamIdFromPlayerState(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState & (2**43-1));
    }

    function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) public pure returns (uint256) {
        require(currentShirtNum < 2**5, "currentShirtNum out of bound");
        state &= ~uint256(2**5-1 << 43); 
        state |= uint256(currentShirtNum) << 43;
        return state;
    }

    function getCurrentShirtNum(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 43 & (2**5-1));
    }
    
    function setPrevPlayerTeamId(uint256 state, uint256 value) public pure returns (uint256) {
        require(value < 2**43, "prevLeagueIdx out of bound");
        state &= ~uint256(2**43-1 << 48); 
        state |= uint256(value) << 48;
        return state;
    }

    function getPrevPlayerTeamId(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 48 & (2**43-1));
    }

    function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) public pure returns (uint256) {
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        state &= ~uint256(2**35-1 << 91); // 256 - 43 - 43 - 5 - 43 - 35
        state |= uint256(lastSaleBlock) << 91;
        return state;
    }

    function getLastSaleBlock(uint256 playerState) public pure returns (uint256) {
        return uint256(playerState >> 91 & (2**35-1));
    }

}
