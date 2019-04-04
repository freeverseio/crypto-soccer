pragma solidity ^0.5.0;

import "./TeamState.sol";

/// @title The state of a League
contract LeagueState is TeamState {
    uint256 constant private DIMENSION_1_END = 0;

    function leagueStateCreate() public pure returns (uint256[] memory state) {
    }

    function leagueStateAppend(uint256[] memory leagueState, uint256[] memory teamState) public pure returns (uint256[] memory state) {
        require(isValidLeagueState(leagueState), "invalid league state");
        require(isValidTeamState(teamState), "invalid team state");
        state = new uint256[](leagueState.length + teamState.length + 1);
        for (uint256 i = 0 ; i < leagueState.length ; i++)
            state[i] = leagueState[i];
        for (uint256 i = 0 ; i < teamState.length ; i++) 
            state[leagueState.length + i] = teamState[i];
        state[leagueState.length + teamState.length] = DIMENSION_1_END;
    }

    function leagueStateSize(uint256[] memory leagueState) public pure returns (uint256 count) {
        require(isValidLeagueState(leagueState), "invalid league state");
        for (uint256 i = 0 ; i < leagueState.length ; i++)
            if (leagueState[i] == DIMENSION_1_END)
                count++;
    }

    function leagueStateAt(uint256[] memory leagueState, uint256 idx) public pure returns (uint256[] memory) {
        require(isValidLeagueState(leagueState), "invalid league state");
        require(idx < leagueStateSize(leagueState), "out of range");
        uint256 first = _getFirstPlayerOfTeam(leagueState, idx);
        uint256 nPlayers;
        while (first+nPlayers < leagueState.length && leagueState[first+nPlayers] != DIMENSION_1_END)
            nPlayers++;
        uint256[] memory state = new uint256[](nPlayers);
        for (uint256 i = 0 ; i < nPlayers ; i++)
            state[i] = leagueState[first+i];
        return state;
    } 

    function leagueStateUpdate(
        uint256[] memory leagueState, 
        uint256 teamIdx, 
        uint256[] memory teamState
    ) 
        public 
        pure 
        returns (uint256[] memory) 
    {
        uint256 nPlayers = teamStateSize(leagueStateAt(leagueState, teamIdx));
        require(nPlayers == teamStateSize(teamState), "mismatch in teams size");
        uint256 firstPlayerIdx = _getFirstPlayerOfTeam(leagueState, teamIdx);
        for (uint256 i = 0; i < teamState.length ; i++)
            leagueState[firstPlayerIdx + i] = teamState[i];
        return leagueState;
    }
   
    function isValidLeagueState(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        if (state[state.length - 1] != DIMENSION_1_END)
            return false;
        return true;
    }

    function _getFirstPlayerOfTeam(uint256[] memory leagueState, uint256 idx) private pure returns (uint256) {
        uint256 teamCounter;
        uint256 i;
        for (i = 0 ; i < leagueState.length && teamCounter < idx; i++){
            if (leagueState[i] == DIMENSION_1_END)
                teamCounter++;
        }
        return i;
    }

    function leagueStateGetSkills(uint256[] memory leagueState) public pure returns (uint256[] memory skills) {
        require(isValidLeagueState(leagueState), "invalid league state");
        skills = new uint256[](leagueState.length);
        for (uint256 i = 0 ; i < leagueState.length ; i++){
            if (leagueState[i] == DIMENSION_1_END)
                skills[i] = DIMENSION_1_END;
            else
                skills[i] = getSkills(leagueState[i]);
        }
    }
    
}