pragma solidity ^0.5.0;

import "./PlayerState.sol";

contract TeamState is PlayerState {
    function isValidTeamState(uint256[] memory state) public pure returns (bool) {
        for (uint256 i = 0 ; i < state.length ; i++)
            if (!isValidPlayerState(state[i]))
                return false;
        return true;
    }

    function getPlayerState(uint256[] memory teamState, uint256 idx) public pure returns (uint256) {
        require(idx < teamState.length, "out of bound");
        require(isValidTeamState(teamState), "invalid team state");
        return teamState[idx];
    }

    function teamStateEvolve(uint256[] memory teamState, uint8 delta) public pure returns (uint256[] memory) {
        require(isValidTeamState(teamState), "invalid team state");
        uint256[] memory state = new uint256[](teamState.length);
        for (uint256 i = 0 ; i < state.length ; i++)
            state[i] = playerStateEvolve(teamState[i], delta);
        return state;
    }

    function teamStateCreate() public pure returns (uint256[] memory state){
    }

    function teamStateCountPlayers(uint256[] memory teamState) public pure returns (uint256 count) {
        require(isValidTeamState(teamState), "invalid team state");
        return teamState.length;
    }

    /// @dev append a player state to team state
    function teamStateAppend(uint256[] memory teamState, uint256 playerState) public pure returns (uint256[] memory state) {
        state = new uint256[](teamState.length + 1);
        for (uint256 i = 0 ; i < teamState.length ; i++)
            state[i] = teamState[i];
        state[state.length-1] = playerState;
    }

    function computeTeamRating(uint256[] memory teamState) public pure returns (uint128 rating) {
        require(isValidTeamState(teamState), "invalid team state");
        for(uint256 i = 0 ; i < teamState.length ; i++){
            rating += uint8(teamState[i] >> 8 * 4 & 0xff);
            rating += uint8(teamState[i] >> 8 * 3 & 0xff);
            rating += uint8(teamState[i] >> 8 * 2 & 0xff);
            rating += uint8(teamState[i] >> 8 & 0xff);
            rating += uint8(teamState[i] & 0xff);
        }
    }
}