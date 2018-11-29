pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

contract CryptoTeams is ERC721, ERC721Enumerable {
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
    uint256 private _nextTeamId = 1;

    function _mint(address to, uint256 tokenId, string memory name) internal {
        require(to != address(0));
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        require(_nameHashTeam[nameHash] == 0);
        _teamProps[tokenId] = Props({name: name, playersIdx: 0});
        _nameHashTeam[nameHash] = tokenId;
        _mint(to, tokenId);
    }

    function addTeam(string memory name, address owner) public {
        _mint(owner, _nextTeamId, name);
        _nextTeamId++;
    }

    function getTeamName(uint idx) public view returns(string) { 
        require(_exists(idx));
        return _teamProps[idx].name;
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        require(_exists(team));
        _teamProps[team].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 team) public view returns (uint256) {
        require(_exists(team));
        return _teamProps[team].playersIdx;
    }
}

