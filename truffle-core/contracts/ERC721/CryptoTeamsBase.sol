pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract CryptoTeamsBase is ERC721, ERC721Enumerable, MinterRole {
    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Props {
        string name;
        uint256 playersIdx;
    }

    mapping(uint256 => Props) private _teamProps;
    mapping(bytes32 => uint256) private _nameHashTeam;

    function _getTeamName(uint256 tokenId) internal view returns(string){
        require(_exists(tokenId));
        return _teamProps[tokenId].name;
    }

    function _setTeamPlayersIdx(uint256 tokenId, uint256 playersIdx) internal {
        require(_exists(tokenId));
        _teamProps[tokenId].playersIdx = playersIdx;
    }

    function _getTeamPlayersIdx(uint256 tokenId) internal view returns (uint256) {
        require(_exists(tokenId));
        return _teamProps[tokenId].playersIdx;
    }

    function _mint(address to, uint256 tokenId, string memory name) internal {
        require(to != address(0));
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        require(_nameHashTeam[nameHash] == 0);
        _teamProps[tokenId] = Props({name: name, playersIdx: 0});
        _nameHashTeam[nameHash] = tokenId;
        _mint(to, tokenId);
    }

    function _mint(address to, uint256 tokenId) internal onlyMinter {
        super._mint(to, tokenId);
    }
}

