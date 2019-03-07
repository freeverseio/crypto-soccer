pragma solidity ^0.5.0;

import "./LeaguesBase.sol";

/// @title an updatable league
contract LeagueUpdatable is LeaguesBase {
    struct Result {
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the day hashes of the league
        bytes32[] dayStateHashes;
        // hash of tactics
        bytes32[] tacticHashes;
        // scores of the league 
        uint16[] scores;
        // updater address
        address updater;
        // update block
        uint256 updateBlock;
    }

    mapping(uint256 => Result) private _result;

    // TODO: add minimum checks
    function updateLeague(
        uint256 id, 
        bytes32 initStateHash,
        bytes32[] memory dayStateHashes,
        bytes32[] memory tacticHashes,
        uint16[] memory scores
    ) 
        public 
    {
        require(_exists(id), "invalid league id");
        _result[id].initStateHash = initStateHash;
        _result[id].dayStateHashes = dayStateHashes;
        _result[id].tacticHashes = tacticHashes;
        _result[id].scores = scores;
        _result[id].updater = msg.sender;
        _result[id].updateBlock = block.number;
    }

    function getUpdater(uint256 id) external view returns (address) {
        require(_exists(id), "unexistent league");
        return _result[id].updater;
    }

    function getUpdateBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _result[id].updateBlock;
    }

    function hashState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTactics(uint256[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(tactics));
    }

    function getInitStateHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _result[id].initStateHash;
    }

    function getDayStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _result[id].dayStateHashes;
    }


    function getTacticHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _result[id].tacticHashes;
    }

    function getScores(uint256 id) external view returns (uint16[] memory) {
        require(_exists(id), "unexistent league");
        return _result[id].scores;
    }
}