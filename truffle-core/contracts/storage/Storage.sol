pragma solidity >= 0.6.3;

import "./ProxyStorage.sol";
import "../storage/Constants.sol";
import "./Stakers.sol";

/**
* @title Storage common to all project, with setters managed by StorageProxy.
*/
contract Storage is ProxyStorage, Constants{

    uint256[2**12] _slotReserve;
   
    address internal _market;
    address internal _COO;
    address internal _relay;
    address internal _cryptoMktAddr;
    
    mapping(uint256 => uint256) internal _playerIdToState;
    mapping (uint256 => uint256) internal _playerIdToAuctionData;
    mapping (uint256 => bool) internal _playerIdToIsFrozenCrypto;
    mapping (uint256 => uint256) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;
    mapping (uint256 => uint256) internal _playerInTransitToTeam;
    mapping (uint256 => uint8) internal _nPlayersInTransitInTeam;
    mapping (uint256 => bool) internal _teamIdToIsBuyNowForbidden;

    uint256 _maxSumSkillsBuyNowPlayer;
    uint256 _maxSumSkillsBuyNowPlayerMinLapse;
    uint256 _maxSumSkillsBuyNowPlayerProposed;
    uint256 _maxSumSkillsBuyNowPlayerMinLapseProposed;
    uint256 _maxSumSkillsBuyNowPlayerLastUpdate;

    uint256 internal nextVerseTimestamp;
    uint8 internal timeZoneForRound1;
    uint256 internal currentVerse;
    bytes32 internal currentVerseSeed;

    uint256 public gameDeployDay;

    mapping (uint256 => uint256) countryIdToNDivisions;
    mapping (uint256 => uint256) countryIdToNHumanTeams;
    mapping (uint256 => uint256) divisionIdToRound;
    mapping (uint256 => uint256[PLAYERS_PER_TEAM_MAX]) teamIdToPlayerIds;
    mapping (uint256 => address) teamIdToOwner;
    mapping (uint8 => uint256) tzToNCountries;

    uint256 firstVerseTimeStamp;
    uint16 _levelsInOneChallenge;
    uint16 _leafsInLeague;
    uint16 _levelsInLastChallenge;
    uint256 _challengeTime;
    bool _allowChallenges;
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
 
    Stakers internal _stakers;
    
    function isCompany(address addr) public view returns (bool) { return addr == _company; }
    function isSuperUser(address addr) public view returns (bool) { return addr == _superUser; }
    function isRelay(address addr) public view returns (bool) { return addr == _relay; }
    function isCOO(address addr) public view returns (bool)  { return addr == _COO; }
}