pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";

/**
 * @title PlayersProps ERC721 Token Standard with its props
 */
contract PlayersProps is ERC721, ERC721Enumerable {
    struct Props {
        string name;
        uint88 genome;
    }

    // Mapping from player ID to Props
    mapping(uint256 => Props) private _playerProps;

    /**
     * @return name of existing player
     */
    function getName(uint256 playerId) external view returns(string memory) {
        require(_exists(playerId));
        return _playerProps[playerId].name;
    }

    /**
     * @return genome of existing player
     */
    function getGenome(uint256 playerId) public view returns (uint88){
        require(_exists(playerId));
        return _playerProps[playerId].genome;
    }

    /**
     * @return birth of existing player
     */
    function getBirth(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome);
    }
    
    /**
     * @return defence of existing player
     */
    function getDefence(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14);
    }
    
    /**
     * @return speed of existing player
     */
    function getSpeed(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 2);
    }
    
    /**
     * @return pass of existing player
     */
    function getPass(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 3);
    }
    
    /**
     * @return shoot of existing player
     */
    function getShoot(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 4);
    }
    
    /**
     * @return endurance of existing player
     */
    function getEndurance(uint256 playerId) external view returns (uint16) {
        require(_exists(playerId));
        return 0x3fff & uint16(_playerProps[playerId].genome >> 14 * 5);
    }

    /**
     * @dev sets name of existing player
     */
    function _setName(uint256 playerId, string memory name) internal {
        require(_exists(playerId));
        _playerProps[playerId].name = name;
    }

    function _setGenome(
        uint256 playerId,
        uint16 birth,
        uint16 defence,
        uint16 speed,
        uint16 pass,
        uint16 shoot,
        uint16 endurance
    ) internal {
        require(_exists(playerId));
        uint88 genome;
        genome |= birth;
        genome |= uint88(defence) << 14;
        genome |= uint88(speed) << 14 * 2;
        genome |= uint88(pass) << 14 * 3;
        genome |= uint88(shoot) << 14 * 4;
        genome |= uint88(endurance) << 14 * 5;
        _playerProps[playerId].genome = genome;
    }
}
