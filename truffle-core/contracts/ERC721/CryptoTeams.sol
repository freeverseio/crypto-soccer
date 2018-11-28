pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";

contract CryptoTeams is ERC721 {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }

    function _mint(address to, uint256 tokenId) internal {
        require(tokenId != 0);
        super._mint(to, tokenId);
    }

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
        uint256 playersIdx;
    }

    /// @dev An array containing the Team struct for all teams in existence. 
    /// @dev The ID of each team is actually his index in this array.
    mapping(uint256 => Team) internal teams;
    // Team[] private teams;a
    uint256 internal teamsCount = 0;

    /// @dev A mapping from team hash(name) to the owner's address.
    /// @dev Facilitates checking if a teamName already exists.
    mapping(bytes32 => uint256) internal teamToOwnerAddr;

    function addTeam(string memory name, address owner) public {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        // require(getTeamOwner(nameHash) == 0);

        teams[teamsCount] = Team({name: name, playersIdx: 0});
        teamToOwnerAddr[nameHash] = teamsCount;
        _mint(owner, teamsCount);
        teamsCount++;
    }

    function getTeamOwner(bytes32 teamHashName) public view returns(address){
        uint256 teamIdx = teamToOwnerAddr[teamHashName];
        return ownerOf(teamIdx);
    }

    function getTeamName(uint idx) public view returns(string) { 
        require(_exists(idx+1));
        return teams[idx+1].name;
    }
}

