pragma solidity >= 0.6.3;

/**
 @title Library of pure functions to serialize match events
 @author Freeverse.io, www.freeverse.io
 @dev eventsLog contains all the info events
*/

/**
each events has 11b: (1b) teamThatAttacks, (1b): managesToShoot, (4b): shooter, (1b): isGoal, (4b): assister
so in total we may have 11*12 = 132b.
WARNING: we add explicit fromEvents suffix to these functions to avoid clashes with other similar functions for matchlog
*/

contract EncodingMatchEventsGetters  {

   function getTeamThatAttacksFromEvents(uint256 log, uint8 round) public pure returns (uint8) {
      return uint8((log >> (11*round)) & 1);
   }

   function getShooterFromEvents(uint256 log, uint8 round) public pure returns (uint8) {
      return uint8((log >> (11*round + 1)) & 15);
   }

   function getIsGoalFromEvents(uint256 log, uint8 round) public pure returns (bool) {
      return ((log >> (11*round + 5)) & 1) == 1;
   }

   function getAssisterFromEvents(uint256 log, uint8 round) public pure returns (uint256) {
      return uint8((log >> (11*round + 6)) & 15);
   }

   function getManagesToShootFromEvents(uint256 log, uint8 round) public pure returns (bool) {
      return ((log >> (11*round + 10)) & 1) == 1;
   }

}
