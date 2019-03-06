pragma solidity ^0.5.0;

import "./LeaguesBase.sol";
import "../state/LeagueState.sol";

/// @title an updatable league
/// TODO: change name to LeagueUpdatable
contract LeagueUpdater is LeaguesBase {
    LeagueState private _state;

    constructor(address state) public {
        _state = LeagueState(state);
    } 

    // TODO: add minimum checks
    function updateLeague(
        uint256 id, 
        bytes32 initStateHash,
        bytes32[] memory finalHashes,
        uint16[] memory scores
    ) 
        public 
    {
        _setInitStateHash(id, initStateHash);
        _setFinalTeamStateHashes(id, finalHashes);
        _setScores(id, scores);
    }

    function hashDayState(uint256[] memory leagueState) public view returns (bytes32[] memory) {
        uint256 nTeams = _state.dayStateSize(leagueState);
        bytes32[] memory hashes = new bytes32[](nTeams);
        for (uint256 i = 0; i < nTeams ; i++){
            uint256[] memory teamState = _state.dayStateAt(leagueState, i);
            hashes[i] = keccak256(abi.encode(teamState));
        }
        return hashes;
    }

    function hashInitState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTeamState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTactics(uint256[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(tactics));
    }

}