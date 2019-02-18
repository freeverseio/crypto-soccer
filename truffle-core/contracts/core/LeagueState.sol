pragma solidity ^0.4.25;

contract LeagueState {
    uint256 constant public DIVIDER = 0;

    function appendTeamToLeagueState(uint256[] memory leagueState, uint256[] memory teamState) public pure returns (uint256[] memory) {
        require(isValid(leagueState), "invalid league state");
        require(isValid(teamState), "invalid team state");

        if(leagueState.length == 0)
            return teamState;
        if(teamState.length == 0)
            return leagueState;

        uint256[] memory state = new uint256[](leagueState.length + teamState.length + 1);
        uint256 i;
        for (i = 0; i < leagueState.length ; i++)
            state[i] = leagueState[i];

        for (i = 0 ; i < teamState.length ; i++)
            state[leagueState.length + 1 + i] = teamState[i];

        return state;        
    }

    function isValid(uint256[] memory state) public pure returns (bool) {
        if (state.length == 0)
            return true;
        if (state[0] == DIVIDER)
            return false;
        if (state[state.length-1] == DIVIDER)
            return false;
        return true;
    }
}