pragma solidity >=0.4.21 <0.6.0;

import "./Storage.sol";

contract Players is Storage {
    constructor(address playerState) public Storage(playerState) {
    }

    function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) public {
        // TODO: check ownership address
        require(_playerExists(playerId0) && _playerExists(playerId1), "unexistent playerId");
        uint256 state0 = getPlayerState(playerId0);
        uint256 state1 = getPlayerState(playerId1);
//         _playerState.playerStateCreate([
//                 _playerState.getDefence(state0), // defence,
//                 _playerState.getSpeed(state0), // defence,
//                 _playerState.getPass(state0), // defence,
//                 _playerState.getShoot(state0), // defence,
//                 _playerState.getEndurance(state0), // defence,
//                 _playerState.getMonthOfBirthInUnixTime(state0), // defence,
//                 playerId0, // defence,
//                 _playerState.getCurrentTeamId(state1), // defence,
//                 _playerState.getCurrentShirtNum(state1), // defence,
//                 _playerState.getPrevLeagueId(state0), // defence,
//                 _playerState.getPrevTeamPosInLeague(state0), // defence,
//                 _playerState.getPrevShirtNumInLeague(state0), // defence,
//                 block.number ]);
// _playerState.playerStateCreate(
//                 _playerState.getDefence(state0), // defence,
//                 _playerState.getSpeed(state0), // defence,
//                 _playerState.getPass(state0), // defence,
//                 _playerState.getShoot(state0), // defence,
//                 _playerState.getEndurance(state0), // defence,
//                 _playerState.getMonthOfBirthInUnixTime(state0), // defence,
//                 playerId0, // defence,
//                 _playerState.getCurrentTeamId(state1), // defence,
//                 _playerState.getCurrentShirtNum(state1), // defence,
//                 _playerState.getPrevLeagueId(state0), // defence,
//                 _playerState.getPrevTeamPosInLeague(state0), // defence,
//                 _playerState.getPrevShirtNumInLeague(state0), // defence,
//                 block.number);
    }

    function getPlayerSkills(uint256 playerId) external view returns (uint16[NUM_SKILLS] memory) {
        require(_playerExists(playerId), "unexistent player");
        return _playerState.getSkillsVec(getPlayerState(playerId));
    }

    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(arg)));
    }
  }