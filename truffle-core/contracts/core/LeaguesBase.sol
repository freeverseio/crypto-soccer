pragma solidity ^0.5.0;

contract LeaguesBase {
    event LeagueCreated(uint256 id);

    struct League {
        uint256 nTeams;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        bytes32 usersInitDataHash;
    }

    League[] private _leagues;

    constructor() public {
        _leagues.push(League(0,0,0,0));
    }

    function leaguesCount() external view returns (uint256) {
        return _leagues.length - 1;
    }

    function create(
        uint256 initBlock, 
        uint256 step,
        uint8 nTeams
    ) 
        public 
    {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(nTeams % 2 == 0, "odd teams count");
        _leagues.push(League(nTeams, initBlock, step, 0));
        uint256 id = _leagues.length - 1;
        emit LeagueCreated(id);
    }

    function getUsersInitDataHash(uint256 id) public view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].usersInitDataHash;
    }

    function getInitBlock(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initBlock;
    }

    function getStep(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].step;
    }

    function getNTeams(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].nTeams;
    }

    function hashUsersInitData(uint256[] memory teamIds, uint8[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(teamIds, tactics));
    }

    function _exists(uint256 id) internal view returns (bool) {
        return _leagues[id].initBlock != 0;
    }
}