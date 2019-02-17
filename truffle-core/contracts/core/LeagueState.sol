pragma solidity ^0.4.25;

contract LeagueState {
    function appendTeamToLeagueState(uint256[] memory target, uint256[] memory teamState) public pure returns (uint256[] memory) {
        require(teamState.length != 0, "wrong team state");
        uint256[] memory state = new uint256[](target.length + teamState.length + 1);
        uint256 i;
        for (i = 0; i < target.length ; i++)
            state[i] = target[i];

        for (i = 0 ; i < teamState.length ; i++)
            state[target.length + i] = teamState[i];

        return state;        
    }
}