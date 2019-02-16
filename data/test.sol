pragma solidity ^0.5.0;

contract SoccerSim {
    
    uint256 constant LEAGUE_COUNT = 10;

    event LeagueChallangeAvailable(uint256 leagueNo, bytes32 value);
    event LeagueChallangeSucessfull(uint256 leagueNo);
    
    mapping(uint256=>uint256) leagues;

    uint256 starts;

    constructor() public {
        next();
    }
    
    function next() public {
        starts = block.number;
        for (uint256 i = 0;i < legueCount();i++) {
            leagues[i]==0;
        }
    }

    function legueCount() public pure returns(uint256) {
        return LEAGUE_COUNT;   
    }

    function canLeagueBeUpdated(uint256 _leagueNo) public view returns(bool) {
        return leagues[_leagueNo]==0 && starts > block.number && starts <  block.number + 10;
    }

    function update(uint256 _leagueNo, bytes32 _value) external {
        leagues[_leagueNo] = 1;
        emit LeagueChallangeAvailable(_leagueNo,_value);
    }
    
    function challange(uint256 _leagueNo, bytes32 _value) external {
        _value;
        leagues[_leagueNo] = 0;
        emit LeagueChallangeSucessfull(_leagueNo);
    }
    
}