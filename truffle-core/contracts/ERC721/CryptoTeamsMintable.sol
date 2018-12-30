pragma solidity ^0.4.24;

import "./CryptoTeamsStorage.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract CryptoTeamsMintable is CryptoTeamsStorage, MinterRole {
    function mintWithName(address to, string name) public onlyMinter {
        uint256 teamId = calculateId(name);
        require(!_exists(teamId));
        _mint(to, teamId);
        _setName(teamId, name);
    }

    function getTeamId(string name) public view returns (uint256) {
        uint256 id = calculateId(name);
        require(_exists(id));
        return id;
    }

    function calculateId(string name) public pure returns (uint256) {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(nameHash);
        return id;
    }
}

