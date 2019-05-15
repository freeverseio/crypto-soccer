pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";
import "../game_controller/GameControllerInterface.sol";

/// @title an updatable league
contract LeagueUpdatable is LeaguesScheduler {
    struct Result {
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the day hashes of the league
        bytes32[] dayStateHashes;
        // scores of the league 
        uint16[] scores;
        // updater address
        address updater;
        // update block
        uint256 updateBlock;
    }

    GameControllerInterface private _stakers;

    mapping(uint256 => Result) private _result;

    function setStakersContract(address stakersContract) public  {
        _stakers = GameControllerInterface(stakersContract);
    }

    // TODO: add minimum checks
    function updateLeague(
        uint256 id, 
        bytes32 initStateHash,
        bytes32[] memory dayStateHashes,
        uint16[] memory scores
    ) 
        public 
    {
        require(_exists(id), "invalid league id");
        require(hasFinished(id), "league not finished");
        require(!isUpdated(id), "already updated");
        _result[id].initStateHash = initStateHash;
        _result[id].dayStateHashes = dayStateHashes;
        _result[id].scores = scores;
        _result[id].updater = msg.sender;
        _result[id].updateBlock = block.number;

        if (_stakers != GameControllerInterface(0))
            _stakers.updated(id, 0, msg.sender);
    }

    function resetUpdater(uint256 id) public {
        require(_exists(id), "unexistent league");
        _result[id].updateBlock = 0;

        if (_stakers != GameControllerInterface(0))
            _stakers.challenged(id, msg.sender);
    }

    function getUpdater(uint256 id) external view returns (address) {
        require(_exists(id), "unexistent league");
        return _result[id].updater;
    }

    function getUpdateBlock(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _result[id].updateBlock;
    }

    function getInitStateHash(uint256 id) public view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _result[id].initStateHash;
    }

    function getDayStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _result[id].dayStateHashes;
    }

    function getScores(uint256 id) public view returns (uint16[] memory) {
        require(_exists(id), "unexistent league");
        return _result[id].scores;
    }

    function isUpdated(uint256 id) public view returns (bool) {
        return _result[id].updateBlock != 0;
    }
}