pragma solidity ^0.5.0;

import "./TeamsProps.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract TeamsMintable is TeamsProps, MinterRole {
    function mint(address to, string memory name) public onlyMinter {
        uint256 teamId = calculateId(name);
        require(!_exists(teamId));
        _mint(to, teamId);
        _setName(teamId, name);
    }

    function getTeamId(string memory name) public view returns (uint256) {
        uint256 id = calculateId(name);
        require(_exists(id));
        return id;
    }

    function calculateId(string memory name) public pure returns (uint256) {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(nameHash);
        return id;
    }
}

