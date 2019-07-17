pragma solidity ^0.5.0;

import "../assets/Assets.sol";

contract LeaguesBase {
    event LeagueCreated(uint256 leagueId);
    uint8 constant public PLAYERS_PER_TEAM = 25;
    Assets private _assets;

    struct League {
        uint8 nTeams;
       // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        bytes32 usersInitDataHash;
        uint8 nTeamsSigned;
    }

    League[] private _leagues;

    constructor() public {
        _leagues.push(League(0,0,0,0,0));
    }

    function setAssetsContract(address assetsContract) public  {
        _assets = Assets(assetsContract);
    }

    function leaguesCount() public view returns (uint256) {
        return _leagues.length - 1;
    }

    function create(
        uint8 nTeams,
        uint256 initBlock, 
        uint256 step
    ) 
        public 
    {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(nTeams % 2 == 0, "odd teams count");
        require(nTeams > 0, "cannot create leagues with no teams");
        _leagues.push(League(nTeams, initBlock, step, 0, 0));
        emit LeagueCreated(leaguesCount());
    }

    function getUsersInitDataHash(uint256 leagueId) public view returns (bytes32) {
        require(_exists(leagueId), "unexistent league");
        return _leagues[leagueId].usersInitDataHash;
    }

    function getInitBlock(uint256 leagueId) public view returns (uint256) {
        require(_exists(leagueId), "unexistent league");
        return _leagues[leagueId].initBlock;
    }

    function getStep(uint256 leagueId) public view returns (uint256) {
        require(_exists(leagueId), "unexistent league");
        return _leagues[leagueId].step;
    }

    function getNTeams(uint256 leagueId) public view returns (uint256) {
        require(_exists(leagueId), "unexistent league");
        return _leagues[leagueId].nTeams;
    }

    function signTeamInLeague(uint256 leagueId, uint256 teamIdx, uint8[PLAYERS_PER_TEAM] memory teamOrder, uint8 teamTactics) public {
        require(_leagues[leagueId].nTeamsSigned < _leagues[leagueId].nTeams, "league already full");
        _leagues[leagueId].usersInitDataHash = keccak256(abi.encode(
            _leagues[leagueId].usersInitDataHash, 
            teamIdx, 
            teamOrder, 
            teamTactics
        )); 
        _leagues[leagueId].nTeamsSigned++;
    }

    function _exists(uint256 leagueId) internal view returns (bool) {
        return leagueId <= leaguesCount();
    }
}