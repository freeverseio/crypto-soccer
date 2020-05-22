pragma solidity >= 0.6.3;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingState {

    uint256 constant internal TWO_TO_43_MINUS_ONE = 8796093022207;
    uint256 constant internal TWO_TO_35_MINUS_ONE = 34359738367;
    uint256 constant internal TWO_TO_5_MINUS_ONE = 31;

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
        require(teamId <= TWO_TO_43_MINUS_ONE, "currentTeamIdx out of bound");
        playerState &= ~ TWO_TO_43_MINUS_ONE;
        playerState |= teamId;
        return playerState;
    }

    function getCurrentTeamIdFromPlayerState(uint256 playerState) public pure returns (uint256) {
        return playerState & TWO_TO_43_MINUS_ONE;
    }

    function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) public pure returns (uint256) {
        require(currentShirtNum <= TWO_TO_5_MINUS_ONE, "currentShirtNum out of bound");
        state &= ~(TWO_TO_5_MINUS_ONE << 43); 
        state |= uint256(currentShirtNum) << 43;
        return state;
    }

    function getCurrentShirtNum(uint256 playerState) public pure returns (uint256) {
        return (playerState >> 43) & TWO_TO_5_MINUS_ONE;
    }
    
    function setPrevPlayerTeamId(uint256 state, uint256 value) public pure returns (uint256) {
        require(value <= TWO_TO_43_MINUS_ONE, "prevLeagueIdx out of bound");
        state &= ~(TWO_TO_43_MINUS_ONE << 48); 
        state |= (value << 48);
        return state;
    }

    function getPrevPlayerTeamId(uint256 playerState) public pure returns (uint256) {
        return (playerState >> 48) & TWO_TO_43_MINUS_ONE;
    }

    function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) public pure returns (uint256) {
        require(lastSaleBlock <= TWO_TO_35_MINUS_ONE, "lastSaleBlock out of bound");
        state &= ~(TWO_TO_35_MINUS_ONE << 91); // 256 - 43 - 43 - 5 - 43 - 35
        state |= (lastSaleBlock << 91);
        return state;
    }

    function getLastSaleBlock(uint256 playerState) public pure returns (uint256) {
        return (playerState >> 91) & TWO_TO_35_MINUS_ONE;
    }

}
