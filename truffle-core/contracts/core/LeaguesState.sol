pragma solidity ^0.5.0;

import "./LeaguesBase.sol";

contract LeaguesState is LeaguesBase {
    uint256 constant public DIVIDER = 0;

    struct State {
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the final hashes of the league
        bytes32[] finalTeamStateHashes;
    }

    mapping(uint256 => State) private _states;

    function getFinalTeamStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _states[id].finalTeamStateHashes;
    }

    function _setFinalTeamStateHashes(uint256 id, bytes32[] memory hashes) internal {
        require(_exists(id), "unexistent league");
        _states[id].finalTeamStateHashes = hashes;
    }

    function _setInitStateHash(uint256 id, bytes32 stateHash) internal {
        require(_exists(id), "unexistent league");
        _states[id].initStateHash = stateHash;
    }

    function getInitStateHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _states[id].initStateHash;
    }

    /**
     * @dev encoding:
     *        defence:   0xff00000000
     *        speed:     0x00ff000000
     *        pass:      0x0000ff0000
     *        shoot:     0x000000ff00
     *        endurance: 0x00000000ff
     */
    function playerStateCreate(
        uint8 defence,
        uint8 speed,
        uint8 pass,
        uint8 shoot,
        uint8 endurance 
    )
        public 
        pure
        returns (uint256 state)
    {
        state |= uint256(defence) << 8 * 4;
        state |= uint256(speed) << 8 * 3;
        state |= uint256(pass) << 8 * 2;
        state |= uint256(shoot) << 8;
        state |= endurance;
    }

    function teamStateCreate() public pure returns (uint256[] memory state){
    }

    /// @dev append a player state to team state
    function teamStateAppend(uint256[] memory teamState, uint256 playerState) public pure returns (uint256[] memory state) {
        state = new uint256[](teamState.length + 1);
        for (uint256 i = 0 ; i < teamState.length ; i++)
            state[i] = teamState[i];
        state[state.length-1] = playerState;
    }

    function leagueStateCreate() public pure returns (uint256[] memory state) {
    }

    function leagueStateAppend(uint256[] memory leagueState, uint256[] memory teamState) public pure returns (uint256[] memory state) {
        require(isValidTeamState(teamState), "invalid team state");
        require(teamState.length != 0, "empty team not allowed");

        if (leagueState.length == 0)
            return teamState;

        state = new uint256[](leagueState.length + teamState.length + 1);
        for (uint256 i = 0 ; i < leagueState.length ; i++)
            state[i] = leagueState[i];
        state[leagueState.length] = DIVIDER;
        for (uint256 i = 0 ; i < teamState.length ; i++) 
            state[leagueState.length + 1 + i] = teamState[i];
    }

    function isValidTeamState(uint256[] memory state) public pure returns (bool) {
        for (uint256 i = 0 ; i < state.length ; i++)
            if (state[i] == DIVIDER)
                return false;
        return true;
    }

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