pragma solidity ^0.5.0;

import "./LeaguesStorage.sol";

contract LeaguesState is LeaguesStorage {
    uint256 constant public DIVIDER = 0;

    function append(uint256[] memory leagueState, uint256[] memory state) public pure returns (uint256[] memory) {
        require(isValid(leagueState), "invalid league result");
        require(isValid(state), "invalid team result");

        if(leagueState.length == 0)
            return state;
        if(state.length == 0)
            return leagueState;

        uint256[] memory result = new uint256[](leagueState.length + state.length + 1);
        uint256 i;
        for (i = 0; i < leagueState.length ; i++)
            result[i] = leagueState[i];
        result[leagueState.length] = DIVIDER;
        for (i = 0 ; i < state.length ; i++)
            result[leagueState.length + 1 + i] = state[i];

        return result;        
    }

    function countTeamsInState(uint256[] memory leagueState) public pure returns (uint256) {
        require(isValid(leagueState), "invalid league state");
        if (leagueState.length == 0)
            return 0;

        uint256 count = 1;
        for (uint256 i = 0 ; i < leagueState.length ; i++) {
            if (leagueState[i] == DIVIDER)
                count++; 
        }
        return count;
    }

    function countTeamPlayers(uint256[] memory leagueState, uint256 idx) public pure returns (uint256) {
        require(isValid(leagueState), "invalid league state");
        require(idx < countTeamsInState(leagueState), "out of range");
        uint256 first = _getFirstPlayerOfTeam(leagueState, idx);
        uint256 counter;
        while (first+counter < leagueState.length && leagueState[first+counter] != DIVIDER)
            counter++;
        return counter;
    }

    function getTeam(uint256[] memory leagueState, uint256 idx) public pure returns (uint256[] memory) {
        require(isValid(leagueState), "invalid league state");
        require(idx < countTeamsInState(leagueState), "out of range");
        uint256 nPlayers = countTeamPlayers(leagueState, idx);
        uint256[] memory state = new uint256[](nPlayers);
        uint256 first = _getFirstPlayerOfTeam(leagueState, idx);
        for (uint256 i = 0 ; i < nPlayers ; i++)
            state[i] = leagueState[first+i];
        return state;
    } 
   
    function isValid(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        if (state[0] == DIVIDER)
            return false;
        if (state[state.length-1] == DIVIDER)
            return false;
        for (uint256 i = 0 ; i < state.length - 1 ; i++)
            if (state[i] == DIVIDER && state[i+1] == DIVIDER)
                return false;
        return true;
    }

    function _getFirstPlayerOfTeam(uint256[] memory leagueState, uint256 idx) private pure returns (uint256) {
        uint256 teamCounter;
        uint256 i;
        for (i = 0 ; i < leagueState.length && teamCounter < idx; i++){
            if (leagueState[i] == DIVIDER)
                teamCounter++;
        }
        return i;
    }
}