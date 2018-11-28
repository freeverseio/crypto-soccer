pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";

contract CryptoTeams is ERC721 {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
        uint256 playersIdx;
    }

    mapping(uint256 => Team) _TeamProps;

    // function getNCreatedTeams() public view returns(uint) {
    //     return teams.length;
    // }

    // function addTeam(string memory name) public {
    //     uint256 tokenId = getNCreatedTeams() - 1;
    //     //_mint()
    //     teams.push(Team({name: name, playersIdx: 0}));
    // }


}

