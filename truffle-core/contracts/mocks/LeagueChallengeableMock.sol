pragma solidity ^0.5.0;

import "../core/LeagueChallengeable.sol";

contract LeagueChallengeableMock is LeagueChallengeable {
    constructor(address engine, address leagueState) LeagueChallengeable(engine, leagueState) public {
    }

    function updateTacticsToBlockNum(
        uint256[] memory usersInitDataTeamIds,
        uint8[3][] memory usersInitDataTactics, 
        uint256 blockNum, 
        uint256[] memory usersAlongDataTeamIds,
        uint8[3][] memory usersAlongDataTactics,
        uint256[] memory usersAlongDataBlocks
    ) 
        public 
        pure 
        returns (uint8[] memory) 
    {
        uint8[3][] memory tactics = _updateTacticsToBlockNum(usersInitDataTeamIds, usersInitDataTactics, blockNum, usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks);
        uint8[] memory result = new uint8[](tactics.length * 3);
        for (uint256 i = 0 ; i < tactics.length ; i++){
            result[3*i] = tactics[i][0];
            result[3*i+1] = tactics[i][1];
            result[3*i+2] = tactics[i][2];
        }
        return result;
    }
}