pragma solidity ^0.5.0;

import "./DayState.sol";

contract LeagueState is DayState {
    // function isValidDayStatePerDay(uint256[] memory state) public pure returns (bool) {
    //     return true;
    // }

    // function leagueStatePerDayCreate() public pure returns (uint256[] memory state) {
    // }

    // /// @return number of days in leagueStatePerDay
    // function leagueStatePerDayCount(uint256[] memory leagueStatePerDay) public pure returns (uint256 count) {
    //     for (uint256 i = 0 ; i < leagueStatePerDay.length ; i++)
    //         count++;
    // }

    // function leagueStatePerDayAppend(
    //     uint256[] memory leagueStatePerDay, 
    //     uint256[] memory leagueState
    // ) 
    //     public 
    //     pure 
    //     returns (uint256[] memory state) 
    // {
    //     require(isValidDayStatePerDay(leagueStatePerDay), "invalid league state per day");
    //     require(isValidDayState(leagueState), "invalid league state");

    //     state = new uint256[](leagueStatePerDay.length + 1 + leagueState.length);
    //     for (uint256 i = 0 ; i < leagueStatePerDay.length ; i++)
    //         state[i] = leagueStatePerDay[i];
    //     state[leagueState.length] = LEAGUESTATEDIVIDER;
    //     for (uint256 i = 0 ; i < leagueState.length ; i++) 
    //         state[leagueStatePerDay.length + 1 + i] = leagueState[i];
    // }
}