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

    function getTeamName(uint256 tokenId) public view returns(string){
        require(_exists(tokenId));
        return _teamProps[tokenId].name;
    }

    function setTeamName(uint256 tokenId, string name) public {
        require(_exists(tokenId));
        _teamProps[tokenId].name = name;
    }

    function setTeamPlayersIdx(uint256 tokenId, uint256 playersIdx) public {
        require(_exists(tokenId));
        _teamProps[tokenId].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 tokenId) public view returns (uint256) {
        require(_exists(tokenId));
        return _teamProps[tokenId].playersIdx;
    }

    function mintWithTeamName(address to, uint256 tokenId, string memory name) public onlyMinter {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        require(_nameHashTeam[nameHash] == 0);
        _mint(to, tokenId);
        setTeamName(tokenId, name);
        _nameHashTeam[nameHash] = tokenId;
    }
}

