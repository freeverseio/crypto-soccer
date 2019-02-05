pragma solidity ^0.4.24;

import "../ERC721/TeamsProps.sol";

/**
 * @title TeamsPropsMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract TeamsPropsMock is TeamsProps {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
    
    function setName(uint256 teamId, string name) public {
        _setName(teamId, name);
    }

    function addPlayer(uint256 teamId, uint256 playerId) public {
        _addPlayer(teamId, playerId);
    }
}