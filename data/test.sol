pragma solidity ^0.5.0;

contract SoccerSim {
    
    uint256 constant UPDATABLE = 1;
    uint256 constant UPDATED = 2;
    uint256 constant CHALLANGABLE = 3;
    uint256 constant CHALLANGED = 4;

    uint256 constant LEAGUE_COUNT = 5;

    event LeagueChallangeAvailable(uint256 leagueNo, bytes32 value);
    event LeagueChallangeSucessfull(uint256 leagueNo);
    
    mapping(uint256=>uint256) leagues;

    constructor() public {
        next();
    }
    
    function next() public {
        for (uint256 i = 0;i < legueCount();i++) {
            leagues[i]=UPDATABLE;
        }
    }

    function legueCount() public pure returns(uint256) {
        return LEAGUE_COUNT;   
    }

    function canLeagueBeUpdated(uint256 _leagueNo) public view returns(bool) {
        return leagues[_leagueNo]==UPDATABLE || leagues[_leagueNo]==CHALLANGED;
    }

    function update(uint256 _leagueNo, bytes32 _value) external {
        require(canLeagueBeUpdated(_leagueNo));
        if (leagues[_leagueNo]==UPDATABLE) {
            if (_leagueNo == 2) {
                leagues[_leagueNo]=CHALLANGABLE;
            } else {
                leagues[_leagueNo]=UPDATED;
            }
        } else if (leagues[_leagueNo]==CHALLANGED) {
            leagues[_leagueNo]=UPDATED;
        }
        emit LeagueChallangeAvailable(_leagueNo,_value);
    }
    
    
    function canLeagueBeChallanged(uint256 _leagueNo) public view returns(bool) {
        return leagues[_leagueNo]==CHALLANGABLE;
    }
    
    function challange(uint256 _leagueNo, bytes32 _value) external {
        require(canLeagueBeChallanged(_leagueNo));
        leagues[_leagueNo] = CHALLANGED;
        emit LeagueChallangeSucessfull(_leagueNo);
    }
    
}