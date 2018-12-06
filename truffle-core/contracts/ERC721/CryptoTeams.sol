pragma solidity ^0.4.24;

import "./CryptoTeamsMetadata.sol";

contract CryptoTeams is CryptoTeamsMetadata {
    function addTeam(string memory name, address owner) public onlyMinter {
        uint256 nextTeamId = totalSupply() + 1;
        mintWithName(owner, nextTeamId, name);
    }

    function setPlayersIds(uint256 tokenId, uint256 playersIdx) public onlyMinter {
        _setPlayersIds(tokenId, playersIdx);
    }
}

