pragma solidity >=0.4.21 <0.6.0;

import "./Storage.sol";

contract Players is Storage {
    constructor(address playerState) public Storage(playerState) {
    }

    /// Get the skills of a player
    function getPlayerSkills(uint256 playerId) external view returns (uint16[NUM_SKILLS] memory) {
        require(_playerExists(playerId), "unexistent player");
        return _playerState.getSkillsVec(getPlayerState(playerId));
    }

 

    /// @return hashed arg casted to uint256
    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(arg)));
    }
  }