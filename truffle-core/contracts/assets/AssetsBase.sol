pragma solidity >=0.4.21 <0.6.0;

import "../state/PlayerState.sol";

/// teamId == 0 is invalid and represents the null team
/// TODO: fix the playerPos <=> playerShirt doubt
contract AssetsBase {
    event TeamCreated (uint256 id);
    event LeagueCreated(uint256 leagueId);

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        bytes32 dna;
        uint256 currentLeagueId;
        uint8 posInCurrentLeague;
        uint256 prevLeagueId;
        uint8 posInPrevLeague;
        uint256[PLAYERS_PER_TEAM] playerIds;
        uint256 creationTimestamp; // timestamp as seconds since unix epoch
        address teamOwner;
    }

    struct League {
        uint8 nTeams;
       // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        bytes32 usersInitDataHash;
        uint8 nTeamsSigned;
    }

    uint8 constant public TEAMS_PER_LEAGUE = 10;
    uint8 constant public PLAYERS_PER_TEAM = 25;
    uint8 constant internal BITS_PER_SKILL = 14;
    uint16 constant internal SKILL_MASK = 0x3fff;
    uint8 constant public NUM_SKILLS = 5;
    address constant public FREEVERSE = address(1);

    mapping(uint256 => uint256) internal _playerIdToState;
    mapping(uint256 => Team) internal _teamIdToTeam;
    mapping(bytes32 => address) internal _teamGivenNameHashToTeamId;

    League[] internal _leagues;
    PlayerState internal _playerState;

}
