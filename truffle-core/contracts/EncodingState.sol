pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingState {

    ////////////////////////////////
    /**
     * @dev PlayerState serializes a total of 169 bits:
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

}