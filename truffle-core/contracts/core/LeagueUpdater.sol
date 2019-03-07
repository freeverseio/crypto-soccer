pragma solidity ^0.5.0;

import "./LeaguesBase.sol";
import "../state/LeagueState.sol";

/// @title an updatable league
/// TODO: change name to LeagueUpdatable
contract LeagueUpdater is LeaguesBase {
    // TODO: add minimum checks
    function updateLeague(
        uint256 id, 
        bytes32 initStateHash,
        bytes32[] memory dayStateHashes,
        uint16[] memory scores
    ) 
        public 
    {
        _setInitStateHash(id, initStateHash);
        _setDayStateHashes(id, dayStateHashes);
        _setScores(id, scores);
    }

    function hashState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTactics(uint256[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(tactics));
    }
}