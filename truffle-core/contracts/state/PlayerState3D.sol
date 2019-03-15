pragma solidity ^0.5.0;

import "./LeagueState.sol";

contract PlayerState3D is LeagueState {
    uint256 constant private DIMENSION_2_END = 1;

    function PlayerState3DCreate() public pure returns (uint256[] memory) {

    }

    function PlayerState3DAppend(
        uint256[] memory playerState2D, 
        uint256[] memory teamState
    ) 
        public 
        pure 
        returns (uint256[] memory state) 
    {
        // state = new u
    }
}