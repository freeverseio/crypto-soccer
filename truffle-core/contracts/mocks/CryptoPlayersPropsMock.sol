pragma solidity ^0.4.24;

import "../ERC721/PlayersProps.sol";

/**
 * @title PlayersPropsMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract PlayersPropsMock is PlayersProps {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function setName(uint256 playerId, string name) public {
        _setName(playerId, name);
    }

    function setGenome(
        uint256 playerId,
        uint16 birth,
        uint16 defence,
        uint16 speed,
        uint16 pass,
        uint16 shoot,
        uint16 endurance
    ) public {
        _setGenome(playerId, birth, defence, speed, pass, shoot, endurance);
    }
}