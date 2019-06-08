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

    // TODO: remove cause local db will store it
    mapping(uint256 => uint256[]) private _leagueToTeams;

    // TODO: remove cause local db will store it
    mapping(uint256 => uint8[3][]) private _leagueToTactics;

    mapping(uint256 => League) private _leagues;
    uint256 private _leaguesCount;


    function leaguesCount() external view returns (uint256) {
        return _leaguesCount;
    }

    function create(
        uint256 id, 
        uint256 initBlock, 
        uint256 step, 
        uint256[] memory teamIds,
        uint8[3][] memory tactics
    ) 
        public 
    {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(teamIds.length % 2 == 0, "odd teams count");
        require(teamIds.length == tactics.length, "nTeams and nTactics mismatch");
        require(!_exists(id), "league already created");
        uint256 nTeams = teamIds.length;
        bytes32 usersInitDataHash = hashUsersInitData(teamIds, tactics);
        _leagues[id] = League(
            nTeams,
            initBlock, 
            step,
            usersInitDataHash
        );
        _leaguesCount++;
        for (uint256 i=0 ; i<teamIds.length ; i++)
            _leagueToTeams[id].push(teamIds[i]);
        for (uint256 i=0 ; i<tactics.length ; i++)
            _leagueToTactics[id].push(tactics[i]);
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

    function getTeams(uint256 id) external view returns (uint256[] memory) {
        require(_exists(id), "unexistent league");
        return _leagueToTeams[id];
    }

    function getTactics(uint256 id) external view returns (uint8[] memory) {
        require(_exists(id), "unexistent league");
        uint8[] memory tactics = new uint8[](_leagueToTactics[id].length*3);
        for (uint256 i=0 ; i < _leagueToTactics[id].length ; i++) {
            tactics[3*i] = _leagueToTactics[id][i][0];
            tactics[3*i + 1] = _leagueToTactics[id][i][1];
            tactics[3*i + 2] = _leagueToTactics[id][i][2];
        }
        return tactics;
    }

    function hashUsersInitData(uint256[] memory teamIds, uint8[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(teamIds, tactics));
    }

    function _exists(uint256 id) internal view returns (bool) {
        return _leagues[id].initBlock != 0;
    }
}