pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/access/Roles.sol";

contract CoachRole {
    using Roles for Roles.Role;

    event CoachAdded(address indexed account);
    event CoachRemoved(address indexed account);

    Roles.Role private _uriers;

    constructor () internal {
        _addCoach(msg.sender);
    }

    modifier onlyCoach() {
        require(isCoach(msg.sender));
        _;
    }

    function isCoach(address account) public view returns (bool) {
        return _uriers.has(account);
    }

    function addCoach(address account) public onlyCoach {
        _addCoach(account);
    }

    function renounceCoach() public {
        _removeCoach(msg.sender);
    }

    function _addCoach(address account) internal {
        _uriers.add(account);
        emit CoachAdded(account);
    }

    function _removeCoach(address account) internal {
        _uriers.remove(account);
        emit CoachRemoved(account);
    }
}