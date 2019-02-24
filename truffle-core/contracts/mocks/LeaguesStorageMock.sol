pragma solidity ^0.5.0;

import "../core/LeaguesStorage.sol";

contract LeaguesStorageMock is LeaguesStorage {
    function setScores(uint256 id, uint256[] memory scores) public {
        _setScores(id, scores);
    }
}