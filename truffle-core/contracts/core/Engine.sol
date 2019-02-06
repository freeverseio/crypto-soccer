pragma solidity ^ 0.4.24;

import "./Leagues.sol";

contract Engine {
    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param stateTeam0 a vector with the state of the players of team 0
     * @param stateTeam1 a vector with the state of the players of team 1
     * @param tacticsTeam0 a vector[3] with the tactic (ex. [4,4,3]) of team 0 
     * @param tacticsTeam0 a vector[3] with the tactic (ex. [4,4,3]) of team 1
     * @return the score of the match
     */
    function playMatch(
        bytes32 seed,
        uint256[] memory stateTeam0,
        uint256[] memory stateTeam1, 
        uint256[3] memory tacticsTeam0, 
        uint256[3] memory tacticsTeam1
    ) 
        public 
        pure 
        returns (uint256, uint256) 
    {
        require(stateTeam0.length >= 11, "Team 0 needs at least 11 players");
        require(stateTeam1.length >= 11, "Team 1 needs at least 11 players");

        uint256 hash0 = uint256(seed) + stateTeam0[0] + tacticsTeam0[0];
        uint256 hash1 = uint256(seed) + stateTeam1[0] + tacticsTeam1[0];
        return (hash0 % 4, hash1 % 4);
    }
}