pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "../factories/teams.sol";

contract CryptoTeams is ERC721 {
    TeamFactory private _teamFactory;

    constructor(address teamFactory) public {
        _teamFactory = TeamFactory(teamFactory);
    }

    function getTeamFactory() external view returns (address) {
        return _teamFactory;
    }
}

