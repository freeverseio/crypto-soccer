pragma solidity ^0.5.0;

import "../core/LeagueUpdatable.sol";

contract LeagueUpdatableMock is LeagueUpdatable {
    function isUpdated(uint256 id) public view returns (bool) {
        return _isUpdated(id);
    }
}