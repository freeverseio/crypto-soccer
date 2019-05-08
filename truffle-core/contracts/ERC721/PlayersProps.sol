pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

/**
 * @title PlayersProps ERC721 Token Standard with its props
 */
contract PlayersProps is ERC721, ERC721Enumerable {
    struct Props {
        string name;
        uint256 genome;
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
     * @return genome of existing player
     */
    function getGenome(uint256 playerId) public view returns (uint256){
        require(_exists(playerId), "playerId not found");
        return _playerProps[playerId].genome;
    }

    /**
     * @return birth of existing player
     */
    function getBirth(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome);
    }
    
    /**
     * @return defence of existing player
     */
    function getDefence(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14);
    }
    
    /**
     * @return speed of existing player
     */
    function getSpeed(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 2);
    }
    
    /**
     * @return pass of existing player
     */
    function getPass(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 3);
    }
    
    /**
     * @return shoot of existing player
     */
    function getShoot(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 4);
    }
    
    /**
     * @return endurance of existing player
     */
    function getEndurance(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId), "playerId not found");
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 5);
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
     * 5x14bits 
     * skills                  = 5x14 bits
     * monthOfBirthInUnixTime  = 14 bits
     * playerIdx               = 28 bits
     * currentTeamIdx          = 28 bits
     * currentShirtNum         =  4 bits
     * prevLeagueIdx           = 25 bits
     * prevTeamPosInLeague     =  8 bits
     * prevShirtNumInLeague    =  4 bits
     * lastSaleBlocknum        = 35 bits 
     * available               = 40 bits
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
        require(playerId > 0 && playerId < 2**28, "playerId out of bound");
        uint256 genome = birth;
        genome |= uint256(defence) << 14;
        genome |= uint256(speed) << 14 * 2;
        genome |= uint256(pass) << 14 * 3;
        genome |= uint256(shoot) << 14 * 4;
        genome |= uint256(endurance) << 14 * 5;
        _playerProps[playerId].genome = genome;
    }


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
        require(playerId > 0 && playerId < 2**28, "playerId out of bound");
        require(currentTeamId < 2**28, "currentTeamIdx out of bound");
        require(currentShirtNum < 2**4, "currentShirtNum out of bound");
        require(prevLeagueId < 2**25, "prevLeagueIdx out of bound");
        require(prevTeamPosInLeague < 2**8, "prevTeamPosInLeague out of bound");
        require(prevShirtNumInLeague < 2**4, "prevShirtNumInLeague out of bound");
        require(lastSaleBlock < 2**35, "lastSaleBlock out of bound");
        uint256 genome = _playerProps[playerId].genome;
        genome |= uint256(currentTeamId) << 14 * 5 + 28;
        genome |= uint256(currentShirtNum) << 14 * 5 + 28 + 4;
        genome |= uint256(prevLeagueId) << 14 * 5 + 28 + 25;
        genome |= uint256(prevTeamPosInLeague) << 14 * 5 + 28 + 25 + 8;
        genome |= uint256(prevShirtNumInLeague) << 14 * 5 + 28 + 25 + 8 + 4;
        genome |= uint256(lastSaleBlock) << 14 * 5 + 28 + 25 + 8 + 4 + 35;
        _playerProps[playerId].genome = genome;
    }
}
