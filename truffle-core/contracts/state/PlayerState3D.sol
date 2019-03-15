pragma solidity ^0.5.0;

import "./LeagueState.sol";

contract PlayerState3D is LeagueState {
    uint256 constant private DIMENSION_2_END = 1;

    function playerState3DCreate() public pure returns (uint256[] memory) {

    }

    function playerState3DAppend(
        uint256[] memory playerState3D, 
        uint256[] memory playerState2D
    ) 
        public 
        pure 
        returns (uint256[] memory state) 
    {
        require(isValidPlayerState3D(playerState3D), "invalid playerState3D");
        require(isValidTeamState(playerState2D), "invalid playerState2D");
        state = new uint256[](playerState2D.length + playerState2D.length + 1);
        for (uint256 i = 0 ; i < playerState2D.length ; i++)
            state[i] = playerState2D[i];
        for (uint256 i = 0 ; i < playerState2D.length ; i++) 
            state[playerState2D.length + i] = playerState2D[i];
        state[playerState2D.length + playerState2D.length] = DIMENSION_2_END;
    }

    function isValidPlayerState3D(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        if (state[state.length - 1] != DIMENSION_2_END)
            return false;
        return true;
    }
}