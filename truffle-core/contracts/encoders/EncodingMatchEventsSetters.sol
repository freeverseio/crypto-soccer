pragma solidity >= 0.6.3;

/**
 @title Library of pure functions to serialize match events
 @author Freeverse.io, www.freeverse.io
 @dev eventsLog contains all the info events
*/

/**
each events has 11b: (1b) teamThatAttacks, (1b): managesToShoot, (4b): shooter, (1b): isGoal, (4b): assister
so in total we may have 11*12 = 132b.
*/

contract EncodingMatchEventsSetters  {

   function setTeamThatAttacks(uint256 log, uint8 round, uint8 teamThatAttacks) public pure returns (uint256) {
      return (log & ~(uint256(1) << (11*round))) | (uint256(teamThatAttacks) << (11*round));
   }

   function setShooter(uint256 log, uint8 round, uint8 player) public pure returns (uint256) {
      return (log & ~(uint256(15) << (11*round + 1))) | (uint256(player) << (11*round + 1));
   }

   function setIsGoal(uint256 log, uint8 round, bool isGoal) public pure returns (uint256) {
      return (log & ~(uint256(1) << (11*round + 5))) | (uint256(isGoal ? 1 : 0) << (11*round + 5));
   }

   function setAssister(uint256 log, uint8 round, uint8 player) public pure returns (uint256) {
      return (log & ~(uint256(1) << (11*round + 6))) | (uint256(player) << (11*round + 6));
   }

   function setManagesToShoot(uint256 log, uint8 round, bool managesToShoot) public pure returns (uint256) {
      return (log & ~(uint256(1) << (11*round + 10))) | (uint256(managesToShoot ? 1 : 0) << (11*round + 10));
   }
}
