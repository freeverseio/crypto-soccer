pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

/// @title CryptoTeamsBase represents a team of players.
/// @notice ERC721 compliant.
contract CryptoTeamsBase is ERC721, ERC721Enumerable {
    struct Props {
        string name;
        uint256[] players;
    }

    mapping(uint256 => Props) private _teamProps;

    function getName(uint256 tokenId) public view returns(string){
        require(_exists(tokenId));
        return _teamProps[tokenId].name;
    }

    function _setName(uint256 teamId, string name) internal {
        require(_exists(teamId));
        _teamProps[teamId].name = name;
    }

    function getPlayers(uint256 teamId) public view returns (uint256[]) {
        require(_exists(teamId));
        return _teamProps[teamId].players;
    }

    function _addPlayer(uint256 teamId, uint256 playerId) internal {
        _teamProps[teamId].players.push(playerId);
    }
}

