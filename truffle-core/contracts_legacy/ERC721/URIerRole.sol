pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/access/Roles.sol";

contract URIerRole {
    using Roles for Roles.Role;

    event URIerAdded(address indexed account);
    event URIerRemoved(address indexed account);

    Roles.Role private _uriers;

    constructor () internal {
        _addURIer(msg.sender);
    }

    modifier onlyURIer() {
        require(isURIer(msg.sender));
        _;
    }

    function isURIer(address account) public view returns (bool) {
        return _uriers.has(account);
    }

    function addURIer(address account) public onlyURIer {
        _addURIer(account);
    }

    function renounceURIer() public {
        _removeURIer(msg.sender);
    }

    function _addURIer(address account) internal {
        _uriers.add(account);
        emit URIerAdded(account);
    }

    function _removeURIer(address account) internal {
        _uriers.remove(account);
        emit URIerRemoved(account);
    }
}