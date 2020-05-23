pragma solidity >= 0.6.3;

import "./ProxyStorage.sol";
import "../storage/Constants.sol";
import "./Stakers.sol";

/**
 @title Storage for all assets.
 @author Freeverse.io, www.freeverse.io
 @dev Treat with great care. All additions must go below the last line
 @dev of the previous version 
*/

contract Storage is ProxyStorage, Constants{

    uint256[2**12] _slotReserve;

    /// Roles
    address public _market;
    address public _COO;
    address public _relay;
    address public _cryptoMktAddr;
    
    /// Assets Storage
    uint256 public gameDeployDay;
    mapping(uint256 => uint256) internal _playerIdToState;
    mapping (uint256 => uint256) internal _playerIdToAuctionData;
    mapping (uint256 => bool) internal _playerIdToIsFrozenCrypto;
    mapping (uint256 => uint256) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;
    mapping (uint256 => uint256) internal _playerInTransitToTeam;
    mapping (uint256 => uint8) internal _nPlayersInTransitInTeam;
    mapping (uint256 => bool) internal _teamIdToIsBuyNowForbidden;
    mapping (uint256 => uint256) internal countryIdToNDivisions;
    mapping (uint256 => uint256) internal countryIdToNHumanTeams;
    mapping (uint256 => uint256) internal divisionIdToRound;
    mapping (uint256 => uint256[PLAYERS_PER_TEAM_MAX]) internal teamIdToPlayerIds;
    mapping (uint256 => address) internal teamIdToOwner;
    mapping (uint8 => uint256) internal tzToNCountries;

    /// Used for restricting skills of players offered via BuyNow pattern
    uint256 internal _maxSumSkillsBuyNowPlayer;
    uint256 internal _maxSumSkillsBuyNowPlayerMinLapse;
    uint256 internal _maxSumSkillsBuyNowPlayerProposed;
    uint256 internal _maxSumSkillsBuyNowPlayerMinLapseProposed;
    uint256 internal _maxSumSkillsBuyNowPlayerLastUpdate;

    /// Storage required by Updates/Challenges games
    /// Time is often measured in verses
    /// A new verse starts with a userActions submission, 
    /// one every 15min
    uint256 internal nextVerseTimestamp;
    uint8 internal timeZoneForRound1;
    uint256 internal currentVerse;
    bytes32 internal currentVerseSeed;

    uint256 public firstVerseTimeStamp;
    uint16 internal _levelsInOneChallenge;
    uint16 internal _leafsInLeague;
    uint16 internal _levelsInLastChallenge;
    uint256 internal _challengeTime;
    bool internal _allowChallenges;
    mapping (uint256 => bytes32[2]) _actionsRoot;
    mapping (uint256 => bytes32[2]) _activeTeamsPerCountryRoot;
    mapping (uint256 => bytes32[2]) _orgMapRoot;
    mapping (uint256 => uint8[2]) _levelVerifiableByBC;
    mapping (uint256 => bytes32[MAX_CHALLENGE_LEVELS][2]) _roots;
    mapping (uint256 => uint8[2]) _challengeLevel;
    mapping (uint256 => uint8) _newestOrgMapIdx;
    mapping (uint256 => uint8) _newestRootsIdx;
    mapping (uint256 => uint256) _lastActionsSubmissionTime;
    mapping (uint256 => uint256) _lastUpdateTime;
 
    /// The update/challenge game needs to call the external
    /// Stakers contract, that only manages deposits, slashes, etc.
    Stakers public _stakers;
}