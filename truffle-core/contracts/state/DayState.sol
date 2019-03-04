pragma solidity ^0.5.0;

import "./TeamState.sol";

contract DayState is TeamState {
    function dayStateCreate() public pure returns (uint256[] memory state) {
    }

    function dayStateAppend(uint256[] memory dayState, uint256[] memory teamState) public pure returns (uint256[] memory state) {
        require(isValidTeamState(teamState), "invalid team state");
        require(teamState.length != 0, "empty team not allowed");

        if (dayState.length == 0)
            return teamState;

        state = new uint256[](dayState.length + teamState.length + 1);
        for (uint256 i = 0 ; i < dayState.length ; i++)
            state[i] = dayState[i];
        state[dayState.length] = TEAMSTATEDIVIDER;
        for (uint256 i = 0 ; i < teamState.length ; i++) 
            state[dayState.length + 1 + i] = teamState[i];
    }

    function countTeamsInState(uint256[] memory dayState) public pure returns (uint256) {
        require(isValidLeagueState(dayState), "invalid league state");
        if (dayState.length == 0)
            return 0;

        uint256 count = 1;
        for (uint256 i = 0 ; i < dayState.length ; i++) {
            if (dayState[i] == TEAMSTATEDIVIDER)
                count++; 
        }
        return count;
    }

    function countTeamPlayers(uint256[] memory dayState, uint256 idx) public pure returns (uint256) {
        require(isValidLeagueState(dayState), "invalid league state");
        require(idx < countTeamsInState(dayState), "out of range");
        uint256 first = _getFirstPlayerOfTeam(dayState, idx);
        uint256 counter;
        while (first+counter < dayState.length && dayState[first+counter] != TEAMSTATEDIVIDER)
            counter++;
        return counter;
    }

    function getTeam(uint256[] memory dayState, uint256 idx) public pure returns (uint256[] memory) {
        require(isValidLeagueState(dayState), "invalid league state");
        require(idx < countTeamsInState(dayState), "out of range");
        uint256 nPlayers = countTeamPlayers(dayState, idx);
        uint256[] memory state = new uint256[](nPlayers);
        uint256 first = _getFirstPlayerOfTeam(dayState, idx);
        for (uint256 i = 0 ; i < nPlayers ; i++)
            state[i] = dayState[first+i];
        return state;
    } 
   
    function isValidLeagueState(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        // first element has to be a valid player state
        if (!isValidPlayerState(state[0]))
            return false;
        // last element has to be a valid player state
        if (!isValidPlayerState(state[state.length-1]))
            return false;
        // consecutive element can't be invalid player state
        for (uint256 i = 0 ; i < state.length - 1 ; i++)
            if (!isValidPlayerState(state[i]) && !isValidPlayerState(state[i+1]))
                return false;
        return true;
    }

    function _getFirstPlayerOfTeam(uint256[] memory dayState, uint256 idx) private pure returns (uint256) {
        uint256 teamCounter;
        uint256 i;
        for (i = 0 ; i < dayState.length && teamCounter < idx; i++){
            if (dayState[i] == TEAMSTATEDIVIDER)
                teamCounter++;
        }
        return i;
    }
}