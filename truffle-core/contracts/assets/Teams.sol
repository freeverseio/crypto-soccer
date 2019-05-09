pragma solidity ^0.5.0;

import "./Storage.sol";

contract Teams is Storage {
    function createTeam(string memory name) public {
        _addTeam(name);
    }
}