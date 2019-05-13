pragma solidity >=0.4.21 <0.6.0;

import "./Storage.sol";

contract Players is Storage {


    constructor(address playerState) public Storage(playerState) {
    }

    function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) public {
        // TODO: check ownership address
        
    }



    function getPlayerSkills(uint256 playerId) external view returns (uint16[NUM_SKILLS] memory) {
        require(_playerExists(playerId), "unexistent player");
        if (_isVirtual(playerId)) {
            uint256 teamId = getPlayerTeam(playerId);
            uint256 posInTeam = getPlayerPosInTeam(playerId);
            string memory teamName = getTeamName(teamId);
            uint256 seed = uint256(keccak256(abi.encodePacked(teamName, posInTeam)));
            return _computeSkills(seed);
        }
        return _playerState.getSkillsVec(getPlayerState(playerId));
    }

    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(arg)));
    }
  }