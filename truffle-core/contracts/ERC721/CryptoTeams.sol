pragma solidity ^0.4.24;

import "./CryptoTeamsBase.sol";
import "./CryptoTeamsMetadata.sol";

contract CryptoTeams is CryptoTeamsBase, CryptoTeamsMetadata {
    function addTeam(string memory name, address owner) public {
        uint256 nextTeamId = totalSupply() + 1;
        mint(owner, nextTeamId, name);
    }
}

