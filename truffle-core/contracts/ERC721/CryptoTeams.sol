pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";

contract CryptoTeams is ERC721 {
    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
        uint256 playersIdx;
    }

    mapping(uint256 => Team) private teams;
    mapping(bytes32 => uint256) private teamToOwnerAddr;
    uint256 private teamsCount = 1;

    function addTeam(string memory name, address owner) public {
        bytes32 nameHash = keccak256(abi.encodePacked(name));
        require(teamToOwnerAddr[nameHash] == 0);
        teams[teamsCount] = Team({name: name, playersIdx: 0});
        teamToOwnerAddr[nameHash] = teamsCount;
        _mint(owner, teamsCount);
        teamsCount++;
    }

    function getTeamName(uint idx) public view returns(string) { 
        require(_exists(idx));
        return teams[idx].name;
    }

    function getNCreatedTeams() public view returns(uint) {
        return teamsCount - 1;
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        require(_exists(team));
        teams[team].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 team) public returns (uint256) {
        require(_exists(team));
        return teams[team].playersIdx;
    }
}

