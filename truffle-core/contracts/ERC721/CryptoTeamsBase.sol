pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";
import "./CryptoPlayersBase.sol";

/// @title CryptoTeamsBase represents a team of players.
/// @notice ERC721 compliant.
contract CryptoTeamsBase is ERC721, ERC721Enumerable, MinterRole {
    struct Props {
        string name;
        uint256[] players;
    }

    mapping(uint256 => Props) private _teamProps;

    function addPlayer(uint256 teamId, uint256 playerId) public {
        _teamProps[teamId].players.push(playerId);
    }

    function getName(uint256 tokenId) public view returns(string){
        require(_exists(tokenId));
        return _teamProps[tokenId].name;
    }

    function getPlayers(uint256 teamId) public view returns (uint256[]) {
        require(_exists(teamId));
        return _teamProps[teamId].players;
    }

    function mintWithName(address to, string memory name) public onlyMinter {
        uint256 teamId = calculateId(name);
        require(!_exists(teamId));
        _mint(to, teamId);
        _teamProps[teamId].name = name;
    }

    function calculateId(string name) public pure returns (uint256) {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(nameHash);
        return id;
    }
}

