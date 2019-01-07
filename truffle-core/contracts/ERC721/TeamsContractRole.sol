pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/access/Roles.sol";

contract TeamsContractRole {
    using Roles for Roles.Role;

    event TeamsContractAdded(address indexed account);
    event TeamsContractRemoved(address indexed account);

    Roles.Role private _uriers;

    constructor () internal {
        _addTeamsContract(msg.sender);
    }

    modifier onlyTeamsContract() {
        require(isTeamsContract(msg.sender));
        _;
    }

    function isTeamsContract(address account) public view returns (bool) {
        return _uriers.has(account);
    }

    function addTeamsContract(address account) public onlyTeamsContract {
        _addTeamsContract(account);
    }

    function renounceTeamsContract() public {
        _removeTeamsContract(msg.sender);
    }

    function _addTeamsContract(address account) internal {
        _uriers.add(account);
        emit TeamsContractAdded(account);
    }

    function _removeTeamsContract(address account) internal {
        _uriers.remove(account);
        emit TeamsContractRemoved(account);
    }
}