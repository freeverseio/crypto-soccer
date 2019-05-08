pragma solidity ^0.5.0;

import "../ERC721/PlayersProps.sol";

/**
 * @title PlayersPropsMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract PlayersPropsMock is PlayersProps {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function setName(uint256 playerId, string memory name) public {
        _setName(playerId, name);
    }

    function setPlayerState(
        uint256 playerId,
        uint16 birth,
        uint16 defence,
        uint16 speed,
        uint16 pass,
        uint16 shoot,
        uint16 endurance,
        uint32 currentTeamId,
        uint8 currentShirtNum,
        uint32 prevLeagueId,
        uint8 prevTeamPosInLeague,
        uint8 prevShirtNumInLeague,
        uint40 lastSaleBlock
    ) public {
        _setGenome(playerId, birth, defence, speed, pass, shoot, endurance);
        _setCurrentHistory(playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock);
    }
}