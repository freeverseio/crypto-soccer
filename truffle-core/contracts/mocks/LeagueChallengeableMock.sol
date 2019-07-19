pragma solidity ^0.5.0;

import "../core/LeagueChallengeable.sol";

contract LeagueChallengeableMock is LeagueChallengeable {
    constructor(address engine, address leagueState) LeagueChallengeable(engine, leagueState) public {
    }

    function updateTacticsToBlockNum(
        uint256[] memory usersInitDataTeamIds,
        uint8[] memory usersInitDataTactics, 
        uint256 blockNum, 
        uint256[] memory usersAlongDataTeamIds,
        uint8[] memory usersAlongDataTactics,
        uint256[] memory usersAlongDataBlocks
    ) 
        public 
        pure 
        returns (uint8[] memory) 
    {
        uint8[] memory tacticsIds = _updateTacticsToBlockNum(usersInitDataTeamIds, usersInitDataTactics, blockNum, usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks);
        uint8[] memory result = new uint8[](tacticsIds.length * 3);
        for (uint256 i = 0 ; i < tacticsIds.length ; i++){
            result[i] = tacticsIds[i];
        }
        return result;
    }
}