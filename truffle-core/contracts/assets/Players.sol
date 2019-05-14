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

    function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) public {
        // TODO: check ownership address
        require(_playerExists(playerId0) && _playerExists(playerId1), "unexistent playerId");
        uint256 state0 = getPlayerState(playerId0);
        uint256 state1 = getPlayerState(playerId1);
        uint256 newState0 = state0;
        newState0 = _playerState.setCurrentTeamId(newState0, _playerState.getCurrentTeamId(state1));
        newState0 = _playerState.setCurrentShirtNum(newState0, _playerState.getCurrentShirtNum(state1));
        state1 = _playerState.setCurrentTeamId(state1,_playerState.getCurrentTeamId(state0));
        state1 = _playerState.setCurrentShirtNum(state1,_playerState.getCurrentShirtNum(state0));
        newState0 = _playerState.setLastSaleBlock(newState0, block.number);
        state1 = _playerState.setLastSaleBlock(state1, block.number);

        // TODO
        // if getBlockNumForLastLeagueOfTeam(teamIdx1, ST) > state1.getLastSaleBlocknum():
        //     state1.prevLeagueIdx = ST.teams[teamIdx1].currentLeagueIdx
        //     state1.prevTeamPosInLeague = ST.teams[teamIdx1].teamPosInCurrentLeague

        // if getBlockNumForLastLeagueOfTeam(teamIdx2, ST) > state2.getLastSaleBlocknum():
        //     state2.prevLeagueIdx = ST.teams[teamIdx2].currentLeagueIdx
        //     state2.prevTeamPosInLeague = ST.teams[teamIdx2].teamPosInCurrentLeague

        _setPlayerState(newState0);
        _setPlayerState(state1);
    }

    /// @return hashed arg casted to uint256
    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(arg)));
    }
  }