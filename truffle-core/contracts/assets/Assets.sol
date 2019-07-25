pragma solidity >=0.4.21 <0.6.0;

import "../state/PlayerState.sol";
import "./AssetsBase.sol";

/// teamId == 0 is invalid and represents the null team
/// TODO: fix the playerPos <=> playerShirt doubt
contract Assets is AssetsBase {

    PlayerState internal _playerState;

    constructor() public {
        _leagues.push(League(0,0,0,0,0));
    }

    function setStatesContract(address statesAddr) public {
        _playerState = PlayerState(statesAddr);
    }

    function botTeamIdToDNA(uint256 teamId) public pure returns (bytes32) {
        (uint256 leagueId, uint8 posInLeague) = botTeamIdToLeagueIdAndPos(teamId); 
        return keccak256(abi.encode(leagueId, posInLeague));
    }

    // teamId = 1, 1, ... (the first team is valid)
    // leagueId = 1, 2, ...  (the first league is dummy)
    function botTeamIdToLeagueId(uint256 teamId) public pure returns (uint256) {
        return 1 + (teamId-1) / TEAMS_PER_LEAGUE;
    }

    function botTeamIdToLeagueIdAndPos(uint256 teamId) public pure returns (uint256, uint8) {
        uint256 leagueId = botTeamIdToLeagueId(teamId);
        uint256 posInLeague = (teamId - 1)  - TEAMS_PER_LEAGUE * (leagueId - 1);
        require(posInLeague < TEAMS_PER_LEAGUE, "Overflow in team to league assignment");
        return (leagueId, uint8(posInLeague));
    }

    function botTeamIdFromLeagueIdAndPos(uint256 leagueId,  uint8 posInLeague) internal pure returns (uint256) {
        return 1 + (leagueId -1) * TEAMS_PER_LEAGUE + posInLeague;
    }

    function isTeamWritten(uint256 teamId) public view returns (bool) {
        return _teamIdToTeam[teamId].creationBlocknum != 0; 
    }

    function isBotTeam(uint256 teamId) public view returns (bool) {
        return _teamExists(teamId) && !isTeamWritten(teamId);
    }

    function botTeamIdToTimeCreation(uint256 teamId) public view returns (uint256) {
        return _leagues[botTeamIdToLeagueId(teamId)].initBlock;
    }

    function getTeamCreationBlocknum(uint256 teamId) public view returns (uint256) {
        require(_teamExists(teamId), "invalid team id to get creation timestamp");
        if (isTeamWritten(teamId)) {
            return _teamIdToTeam[teamId].creationBlocknum;
        } else {
            return botTeamIdToTimeCreation(teamId);
        }
    }

    function getCurrentLeagueId(uint256 teamId) public view returns (uint256) {
        require(_teamExists(teamId), "invalid team id");
        if (isTeamWritten(teamId)) {
            return _teamIdToTeam[teamId].currentLeagueId;
        } else {
            return botTeamIdToLeagueId(teamId); 
        }
    }

    /// get the current and previous team league and position in league
    function getTeamCurrentHistory(uint256 teamId) external view returns (
            uint256 currentLeagueId,
            uint8 posInCurrentLeague,
            uint256 prevLeagueId,
            uint8 posInPrevLeague
        )
    {
        require(_teamExists(teamId), "invalid team id");
        if (isTeamWritten(teamId)) {
            return(
                _teamIdToTeam[teamId].currentLeagueId,
                _teamIdToTeam[teamId].posInCurrentLeague,
                _teamIdToTeam[teamId].prevLeagueId,
                _teamIdToTeam[teamId].posInPrevLeague
            );
        } else {
            (uint256 leagueId, uint8 posInLeague) = botTeamIdToLeagueIdAndPos(teamId); 
            return(
                leagueId,
                posInLeague,
                0,
                0
            );
        }
    }


    /// @dev Transfers a team to a new owner. 
    /// @dev This function should be called only when the transfer is legit, as checked elsewhere.
    function transferTeam(uint256 teamId, address newOwner) public {
        _teamExists(teamId);
        require(newOwner != address(0), "meaningless adress");
        if (isTeamWritten(teamId)) {
            require(newOwner != getTeamOwner(teamId), "unable to transfer between the same user");
            _teamIdToTeam[teamId].teamOwner = newOwner;
        } else {
            uint256[PLAYERS_PER_TEAM] memory playerIds; 
            (uint256 leagueId, uint8 posInLeague) = botTeamIdToLeagueIdAndPos(teamId); 
            _teamIdToTeam[teamId] = Team(
                botTeamIdToDNA(teamId), 
                leagueId, 
                posInLeague, 
                0, 
                0, 
                playerIds, 
                botTeamIdToTimeCreation(teamId),
                newOwner
            );
        }
    }

    // TODO: exchange fails on playerId0 & playerId1 of the same team
    function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) public {
        // TODO: check ownership address
        require(_playerExists(playerId0) && _playerExists(playerId1), "unexistent playerId");
        uint256 state0 = getPlayerState(playerId0);
        uint256 state1 = getPlayerState(playerId1);
        uint256 newState0 = state0;
        uint256 teamId0 = _playerState.getCurrentTeamId(state0);
        uint256 teamId1 = _playerState.getCurrentTeamId(state1);
        require(!isBotTeam(teamId0) && !isBotTeam(teamId1), "Players in BotTeams cannot be traded");
        uint256 playerShirt0 = _playerState.getCurrentShirtNum(state0);
        uint256 playerShirt1 = _playerState.getCurrentShirtNum(state1);
        newState0 = _playerState.setCurrentTeamId(newState0, _playerState.getCurrentTeamId(state1));
        newState0 = _playerState.setCurrentShirtNum(newState0, _playerState.getCurrentShirtNum(state1));
        state1 = _playerState.setCurrentTeamId(state1,_playerState.getCurrentTeamId(state0));
        state1 = _playerState.setCurrentShirtNum(state1,_playerState.getCurrentShirtNum(state0));
        newState0 = _playerState.setLastSaleBlock(newState0, block.number);
        state1 = _playerState.setLastSaleBlock(state1, block.number);

        _teamIdToTeam[teamId0].playerIds[playerShirt0] = playerId1;
        _teamIdToTeam[teamId1].playerIds[playerShirt1] = playerId0;

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

    // function createTeam(string memory name, address owner) public {
    //     bytes32 nameHash = keccak256(abi.encode(name));
    //     require(_teamNameHashToOwner[nameHash] == address(0), "team already exists");
    //     _teamNameHashToOwner[nameHash] = owner;
    //     uint256[PLAYERS_PER_TEAM] memory playerIds;
    //     teams.push(Team(nameHash, 0, 0, 0, 0, playerIds, block.timestamp));
    //     uint256 id = teams.length - 1;
    //     emit TeamCreated(id);
    // }

    // @dev When signing a team to a new league, it updates current and prev data
    // @dev It does not check legitimacy of this step. Should be done before calling this.
    function _updateTeamHistory(
        uint256 teamId,
        uint256 leagueId,
        uint8 posInLeague
    )
    public
    {
        _teamIdToTeam[teamId].prevLeagueId = _teamIdToTeam[teamId].currentLeagueId;
        _teamIdToTeam[teamId].posInPrevLeague = _teamIdToTeam[teamId].posInCurrentLeague;
        _teamIdToTeam[teamId].currentLeagueId = leagueId;
        _teamIdToTeam[teamId].posInCurrentLeague = posInLeague;
    }

    // TODO: exception when not existent team
    function getTeamOwner(uint256 teamId) public view returns (address) {
        if (isTeamWritten(teamId)){ 
            return _teamIdToTeam[teamId].teamOwner;
        } else {
            return FREEVERSE;
        }
        
    }

    function countTeams() public view returns (uint256){
        return countLeagues() * TEAMS_PER_LEAGUE;
    }

    function getTeamDNA(uint256 teamId) public view returns (bytes32) {
        require(_teamExists(teamId), "invalid team id when getting DNA");
        if (isTeamWritten(teamId)){ 
            return _teamIdToTeam[teamId].dna;
        } else {
            return botTeamIdToDNA(teamId);
        }
    }

    function getTeamPlayerIds(uint256 teamId) public view returns (uint256[PLAYERS_PER_TEAM] memory playerIds) {
        require(_teamExists(teamId), "invalid team id when getTeamPlayerIds");
        if (isTeamWritten(teamId)){ 
            for (uint8 pos = 0 ; pos < PLAYERS_PER_TEAM ; pos++){
                if (_teamIdToTeam[teamId].playerIds[pos] == 0) // virtual player
                    playerIds[pos] = teamIdAndPosToVirtualPlayerId(teamId, pos);
                else
                    playerIds[pos] = _teamIdToTeam[teamId].playerIds[pos];
            }
        } else { // bot team => all players are necessarily virtual
            for (uint8 pos = 0 ; pos < PLAYERS_PER_TEAM ; pos++){
                playerIds[pos] = teamIdAndPosToVirtualPlayerId(teamId, pos);
            }            
        }
    }

    function getPlayerState(uint256 playerId) public view returns (uint256) {
        require(_playerExists(playerId), "unexistent player");
        if (_isPlayerVirtual(playerId))
            return generateVirtualPlayerState(playerId);
        else
            return _playerIdToState[playerId];
    }

    function teamIdAndPosToVirtualPlayerId(uint256 teamId, uint8 posInTeam) public view returns (uint256) {
        require(_teamExists(teamId), "unexistent team");
        require(posInTeam < PLAYERS_PER_TEAM, "invalid player pos");
        return PLAYERS_PER_TEAM * (teamId - 1) + 1 + posInTeam;
    }

    function generateVirtualPlayerState(uint256 playerId) public view returns (uint256) {
            require(_playerExists(playerId), "playerId invalid");
            uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM;
            uint256 posInTeam = playerId - PLAYERS_PER_TEAM * (teamId - 1) - 1;
            uint256 seed = _computeSeed(getTeamDNA(teamId), posInTeam);
            uint16[5] memory skills = _computeSkills(seed);
            uint16 birth = _computeBirth(seed, getTeamCreationBlocknum(teamId));
            return _playerState.playerStateCreate(
                skills[0], // defence,
                skills[1], // speed,
                skills[2], // pass,
                skills[3], // shoot,
                skills[4], // endurance,
                birth, // monthOfBirthInUnixTime,
                playerId,
                teamId,
                posInTeam, // currentShirtNum,
                0, // prevLeagueId,
                0, // prevTeamPosInLeague,
                0, // prevShirtNumInLeague,
                0 // lastSaleBloc
            );
    }

    function _setPlayerState(uint256 state) internal {
        uint256 playerId = _playerState.getPlayerId(state);
        require(_playerExists(playerId), "unexistent player");
        uint256 teamId = _playerState.getCurrentTeamId(state);
        require(_teamExists(teamId), "unexistent team");
        uint256 shirtNumber = _playerState.getCurrentShirtNum(state);
        require(shirtNumber < PLAYERS_PER_TEAM, "invalid shirt number");
        shirtNumber = _playerState.getPrevShirtNumInLeague(state);
        require(shirtNumber < PLAYERS_PER_TEAM, "invalid shirt number");
        uint256 saleBlock = _playerState.getLastSaleBlock(state);
        require(saleBlock != 0 && saleBlock <= block.number, "invalid sale block");
        _playerIdToState[playerId] = state;
    }

    function _teamExists(uint256 teamId) internal view returns (bool) {
        return botTeamIdToLeagueId(teamId) <= countLeagues();
    }

    function _playerExists(uint256 playerId) internal view returns (bool) {
        if (playerId == 0) return false;
        if (_playerIdToState[playerId] != 0) return true;
        uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM;
        return teamId <= countTeams();
    }

    function _isPlayerVirtual(uint256 playerId) internal view returns (bool) {
        require(_playerExists(playerId), "unexistent player");
        return _playerIdToState[playerId] == 0;
    }

    /// Compute a random age between 16 and 35
    /// @param rnd is a random number used as seed of the skills
    /// @param currentTime in seconds since unix epoch
    /// @return monthOfBirth in monthUnixTime
    function _computeBirth(uint256 rnd, uint256 currentTime) internal pure returns (uint16) {
        rnd >>= BITS_PER_SKILL*NUM_SKILLS;
        uint16 seed = uint16(rnd & SKILL_MASK);
        /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
        uint16 age = 16 + (seed % 20);

        /// @dev Convert age to monthOfBirthAfterUnixEpoch.
        /// @dev I leave it this way for clarity, for the time being.
        uint years2secs = 365 * 24 * 3600; // TODO: make it a constant
        uint month2secs = 30 * 24 * 3600; // TODO: make it a constant

        return uint16((currentTime - age * years2secs) / month2secs);
    }

    /// Compute the pseudorandom skills, sum of the skills is 250
    /// @param rnd is a random number used as seed of the skills
    /// @return 5 skills
    function _computeSkills(uint256 rnd) internal pure returns (uint16[NUM_SKILLS] memory) {
        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(rnd & SKILL_MASK);
            rnd >>= BITS_PER_SKILL;
        }

        /// The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < 5; i++) {
            skills[i] = skills[i] % 50;
            excess += skills[i];
        }

        /// At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        uint16 delta = (250 - excess) / 5;
        for (uint8 i = 0; i < 5; i++)
            skills[i] = skills[i] + delta;

        uint16 remainder = (250 - excess) % 5;
        for (uint8 i = 0 ; i < remainder ; i++)
            skills[i]++;

        return skills;
    }

    /// @return seed
    function _computeSeed(bytes32 teamNameHash, uint256 posInTeam) internal pure returns (uint256) {
        return uint256(keccak256(abi.encode(teamNameHash, posInTeam)));
    }

    /// @return hashed arg casted to uint256
    function _intHash(string memory arg) internal pure returns (uint256) {
        return uint256(keccak256(abi.encode(arg)));
    }
    
        function countLeagues() public view returns (uint256) {
        return _leagues.length - 1;
    }
    
    /// LEAGUE FUNCTIONS
    function createLeague(
        uint256 initBlock, 
        uint256 step
    ) 
        public 
    {
        require(initBlock > block.number, "invalid init block");
        require(initBlock < block.number + MAX_INITBLOCK_DELAY, "cannot create a league too far in future");
        require(step > 0, "invalid block step");
        _leagues.push(League(TEAMS_PER_LEAGUE, initBlock, step, 0, 0));
        emit LeagueCreated(countLeagues());
    }

    function leagueIdAndPosToTeamId(uint256 leagueId, uint8 posInLeague) public pure returns (uint256) {
        return (leagueId -1) * TEAMS_PER_LEAGUE + posInLeague;
    }

    function getUsersInitDataHash(uint256 leagueId) public view returns (bytes32) {
        require(_leagueExists(leagueId), "unexistent league");
        return _leagues[leagueId].usersInitDataHash;
    }

    function getInitBlock(uint256 leagueId) public view returns (uint256) {
        require(_leagueExists(leagueId), "unexistent league");
        return _leagues[leagueId].initBlock;
    }

    function getStep(uint256 leagueId) public view returns (uint256) {
        require(_leagueExists(leagueId), "unexistent league");
        return _leagues[leagueId].step;
    }

    function getNTeams(uint256 leagueId) public view returns (uint256) {
        require(_leagueExists(leagueId), "unexistent league");
        return _leagues[leagueId].nTeams;
    }

    // @dev the callee should first verify that the teams are not already involved in un-verified leagues
    function _signTeamInLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague, uint8[PLAYERS_PER_TEAM] memory teamOrder, uint8 teamTactics) internal {
        require(!isBotTeam(teamId), "BotTeams cannot sign a new league");
        require(getCurrentLeagueId(teamId) != leagueId, "cannot sign to a league twice");
        require(isBotTeam(botTeamIdFromLeagueIdAndPos(leagueId, posInLeague)), "this place in the league is already occupied");
        _updateTeamHistory(teamId, leagueId, _leagues[leagueId].nTeamsSigned);
        _leagues[leagueId].usersInitDataHash = keccak256(abi.encode(
            _leagues[leagueId].usersInitDataHash, 
            teamId, 
            teamOrder, 
            teamTactics
        )); 
    }

    function _leagueExists(uint256 leagueId) internal view returns (bool) {
        return leagueId <= countLeagues();
    }
    
    
}
