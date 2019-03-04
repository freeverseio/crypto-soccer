pragma solidity ^0.5.0;

import "./LeaguesBase.sol";

contract LeaguesProof is LeaguesBase {
    struct State {
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the final hashes of the league
        bytes32[] finalTeamStateHashes;
    }

    mapping(uint256 => State) private _states;

    function getInitStateHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _states[id].initStateHash;
    }

    function getFinalTeamStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _states[id].finalTeamStateHashes;
    }

    function _setInitStateHash(uint256 id, bytes32 stateHash) internal {
        require(_exists(id), "unexistent league");
        _states[id].initStateHash = stateHash;
    }

    function _setFinalTeamStateHashes(uint256 id, bytes32[] memory hashes) internal {
        require(_exists(id), "unexistent league");
        _states[id].finalTeamStateHashes = hashes;
    }
}