pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

/**
 * @title PlayersProps ERC721 Token Standard with its props
 */
contract PlayersProps is ERC721, ERC721Enumerable {
    struct Props {
        string name;
        uint256 playerState;
    }

    // Mapping from player ID to Props
    mapping(uint256 => Props) private _playerProps;

    /**
     * @return name of existing player
     */
    function getName(uint256 playerId) external view returns(string memory) {
        require(_exists(playerId), "playerId not found");
        return _playerProps[playerId].name;
    }

    /**
     * @return playerState of existing player
     */
    function getPlayerState(uint256 playerId) public view returns (uint256){
        require(_exists(playerId), "playerId not found");
        return _playerProps[playerId].playerState;
    }

    /**
     * @return birth of existing player
     */
    function getBirth(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState);
    }
    
    /**
     * @return defence of existing player
     */
    function getDefence(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState >> 14);
    }
    
    /**
     * @return speed of existing player
     */
    function getSpeed(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState >> 14 * 2);
    }
    
    /**
     * @return pass of existing player
     */
    function getPass(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState >> 14 * 3);
    }
    
    /**
     * @return shoot of existing player
     */
    function getShoot(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState >> 14 * 4);
    }
    
    /**
     * @return endurance of existing player
     */
    function getEndurance(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].playerState >> 14 * 5);
    }

    /**
     * @dev sets name of existing player
     */
    function _setName(uint256 playerId, string memory name) internal {
        require(_exists(playerId), "playerId not found");
        _playerProps[playerId].name = name;
    }

    /**
     * @dev encoding:
     * skills                  = 5x14 bits
     * monthOfBirthInUnixTime  = 14 bits
     * playerIdx               = 28 bits (TODO: add)
     */

    function _setGenome(
        uint256 playerId,
        uint16 birth,
        uint16 defence,
        uint16 speed,
        uint16 pass,
        uint16 shoot,
        uint16 endurance
    ) internal {
        require(_exists(playerId), "playerId not found");
        require(defence < 2**14, "defence out of bound");
        require(speed < 2**14, "defence out of bound");
        require(pass < 2**14, "defence out of bound");
        require(shoot < 2**14, "defence out of bound");
        require(endurance < 2**14, "defence out of bound");
        require(birth < 2**14, "birth out of bound");
        uint256 playerState = birth;
        playerState |= uint256(defence) << 14;
        playerState |= uint256(speed) << 14 * 2;
        playerState |= uint256(pass) << 14 * 3;
        playerState |= uint256(shoot) << 14 * 4;
        playerState |= uint256(endurance) << 14 * 5;
        _playerProps[playerId].playerState = playerState;
    }


    /**
     * @dev encoding:
     * the genome chunk is already 6x14 bits
     * currentTeamIdx          = 28 bits
     * currentShirtNum         =  4 bits
     * prevLeagueIdx           = 25 bits
     * prevTeamPosInLeague     =  8 bits
     * prevShirtNumInLeague    =  4 bits
     * lastSaleBlocknum        = 35 bits 
     * available               = 40 bits
     */

    function _setCurrentHistory(
        uint256 playerId,
        uint32 currentTeamId,
        uint8 currentShirtNum,
        uint32 prevLeagueId,
        uint8 prevTeamPosInLeague,
        uint8 prevShirtNumInLeague,
        uint40 lastSaleBlock
    ) internal {
        require(_exists(playerId), "playerId not found");
        require(currentTeamId < 2**28, "currentTeamIdx out of bound");
        require(currentShirtNum < 2**4, "currentShirtNum out of bound");
        require(prevLeagueId < 2**25, "prevLeagueIdx out of bound");
        require(prevTeamPosInLeague < 2**8, "prevTeamPosInLeague out of bound");
        require(prevShirtNumInLeague < 2**4, "prevShirtNumInLeague out of bound");
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        uint256 playerState = _playerProps[playerId].playerState;
        playerState |= uint256(currentTeamId) << 14 * 6 + 28;
        playerState |= uint256(currentShirtNum) << 14 * 6 + 28 + 4;
        playerState |= uint256(prevLeagueId) << 14 * 6 + 28 + 25;
        playerState |= uint256(prevTeamPosInLeague) << 14 * 6 + 28 + 25 + 8;
        playerState |= uint256(prevShirtNumInLeague) << 14 * 6 + 28 + 25 + 8 + 4;
        playerState |= uint256(lastSaleBlock) << 14 * 6 + 28 + 25 + 8 + 4 + 35;
        _playerProps[playerId].playerState = playerState;
    }

    function _mint(address to, uint256 playerId) internal {
        require(playerId != 0, "id 0 not allowed");
        require(playerId < 2**28, "playerId out of bound");
        super._mint(to, playerId);
    }
}
