pragma solidity ^0.5.0;

import "./Leagues.sol";

contract Engine {
    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param state0 a vector with the state of the players of team 0
     * @param state1 a vector with the state of the players of team 1
     * @param tactic0 a vector[3] with the tactic (ex. [4,4,2]) of team 0 
     * @param tactic1 a vector[3] with the tactic (ex. [4,4,2]) of team 1
     * @return the score of the match
     */
    function playMatch(
        bytes32 seed,
        uint256[] memory state0,
        uint256[] memory state1, 
        uint8[3] memory tactic0, 
        uint8[3] memory tactic1
    ) 
        public 
        pure 
        returns (uint8 home, uint8 visitor) 
    {
        require(state0.length >= 11, "Team 0 needs at least 11 players");
        require(state1.length >= 11, "Team 1 needs at least 11 players");
        require(tactic0[0] + tactic0[1] + tactic0[2] == 10, "wrong tactic for team 0");
        require(tactic1[0] + tactic1[1] + tactic1[2] == 10, "wrong tactic for team 1");
        bytes32 hash0 = keccak256(abi.encode(uint256(seed) + state0[0] + tactic0[0]));
        bytes32 hash1 = keccak256(abi.encode(uint256(seed) + state1[0] + tactic1[0]));
        return (uint8(uint256(hash0) % 4), uint8(uint256(hash1) % 4));
    }
}