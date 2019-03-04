pragma solidity ^0.5.0;

import "./TeamState.sol";

contract DayState is TeamState {
    function dayStateCreate() public pure returns (uint256[] memory state) {
    }

    function dayStateAppend(uint256[] memory dayState, uint256[] memory teamState) public pure returns (uint256[] memory state) {
        require(isValidTeamState(teamState), "invalid team state");
        state = new uint256[](dayState.length + teamState.length + 1);
        for (uint256 i = 0 ; i < dayState.length ; i++)
            state[i] = dayState[i];
        for (uint256 i = 0 ; i < teamState.length ; i++) 
            state[dayState.length + i] = teamState[i];
        state[dayState.length + teamState.length] = TEAMSTATEEND;
    }

    function dayStateSize(uint256[] memory dayState) public pure returns (uint256 count) {
        require(isValidDayState(dayState), "invalid league state");
        for (uint256 i = 0 ; i < dayState.length ; i++)
            if (dayState[i] == TEAMSTATEEND)
                count++;
    }

    function dayStateAt(uint256[] memory dayState, uint256 idx) public pure returns (uint256[] memory) {
        require(isValidDayState(dayState), "invalid league state");
        require(idx < dayStateSize(dayState), "out of range");
        uint256 first = _getFirstPlayerOfTeam(dayState, idx);
        uint256 nPlayers;
        while (first+nPlayers < dayState.length && dayState[first+nPlayers] != TEAMSTATEEND)
            nPlayers++;
        uint256[] memory state = new uint256[](nPlayers);
        for (uint256 i = 0 ; i < nPlayers ; i++)
            state[i] = dayState[first+i];
        return state;
    } 
   
    function isValidDayState(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        if (state[state.length - 1] != TEAMSTATEEND)
            return false;
        return true;
    }

    function _getFirstPlayerOfTeam(uint256[] memory dayState, uint256 idx) private pure returns (uint256) {
        uint256 teamCounter;
        uint256 i;
        for (i = 0 ; i < dayState.length && teamCounter < idx; i++){
            if (dayState[i] == TEAMSTATEEND)
                teamCounter++;
        }
        return i;
    }
}