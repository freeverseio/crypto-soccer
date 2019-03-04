pragma solidity ^0.5.0;

import "./LeaguesState.sol";

contract LeaguesStatePerDay is LeaguesState {
    function isValidLeagueStatePerDay(uint256[] memory state) public pure returns (bool) {
        return true;
    }

    function leagueStatePerDayCreate() public pure returns (uint256[] memory state) {
    }

    function leagueStatePerDayAppend(
        uint256[] memory leagueStatePerDay, 
        uint256[] memory leagueState
    ) 
        public 
        pure 
        returns (uint256[] memory state) 
    {
        require(isValidLeagueStatePerDay(leagueStatePerDay), "invalid league state per day");
        require(isValidLeagueState(leagueState), "invalid league state");

        if (leagueStatePerDay.length == 0)   
            return leagueState;
        
        state = new uint256[](leagueStatePerDay.length + 1 + leagueState.length);
        for (uint256 i = 0 ; i < leagueStatePerDay.length ; i++)
            state[i] = leagueStatePerDay[i];
        state[leagueState.length] = LEAGUESTATEDIVIDER;
        for (uint256 i = 0 ; i < leagueState.length ; i++) 
            state[leagueStatePerDay.length + 1 + i] = leagueState[i];
    }
}