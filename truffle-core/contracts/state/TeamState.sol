pragma solidity >=0.4.21 <0.6.0;

import "./PlayerState.sol";

/// @title the state of a team
contract TeamState is PlayerState {
    function teamStateCreate() public pure returns (uint256[] memory state){
    }

    /// Append a player state to team state
    function teamStateAppend(uint256[] memory teamState, uint256 playerState) public pure returns (uint256[] memory state) {
        state = new uint256[](teamState.length + 1);
        for (uint256 i = 0 ; i < teamState.length ; i++)
            state[i] = teamState[i];
        state[state.length-1] = playerState;
    }

    /// @return how many player state are in team state
    function teamStateSize(uint256[] memory teamState) public pure returns (uint256 count) {
        return teamState.length;
    }

    /// @return player state at teamState[idx]
    function teamStateAt(uint256[] memory teamState, uint256 idx) public pure returns (uint256 playerState) {
        require(idx < teamState.length, "out of bound");
        playerState = teamState[idx];
    }

    /// Evolve the team of delta
    function teamStateEvolve(uint256[] memory teamState, uint8 delta) public pure returns (uint256[] memory) {
        for (uint256 i = 0 ; i < teamState.length ; i++)
            teamState[i] = playerStateEvolve(teamState[i], delta);
        return teamState;
    }

    function computeTeamRating(uint256[] memory teamState) public pure returns (uint256 rating) {
        for(uint256 i = 0 ; i < teamState.length ; i++){
            uint256 playerState = teamState[i];
            rating += getDefence(playerState);
            rating += getSpeed(playerState);
            rating += getPass(playerState);
            rating += getShoot(playerState);
            rating += getEndurance(playerState);
        }
    }
}