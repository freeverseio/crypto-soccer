pragma solidity >=0.5.12 <=0.6.3;

import "./ProxyStorage.sol";
import "./Constants.sol";
import "./Stakers.sol";

/**
* @title Storage common to all project, with setters managed by StorageProxy.
*/
contract Storage is ProxyStorage, Constants{

    uint256[2**12] _slotReserve;
   
    address internal _academyAddr;

    mapping(uint256 => uint256) internal _playerIdToState;
    mapping (uint256 => uint256) internal _playerIdToAuctionData;
    mapping (uint256 => bool) internal _playerIdToIsFrozenCrypto;
    mapping (uint256 => uint256) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;
    mapping (uint256 => uint256) internal _playerInTransitToTeam;
    mapping (uint256 => uint8) internal _nPlayersInTransitInTeam;


    uint256 internal nextVerseTimestamp;
    uint8 internal timeZoneForRound1;
    uint256 internal currentVerse;
    bytes32 internal currentVerseSeed;

    uint256 public gameDeployDay;
    uint256 public currentRound;

    mapping (uint256 => uint256) countryIdToNDivisions;
    mapping (uint256 => uint256) countryIdToNHumanTeams;
    mapping (uint256 => uint256) divisionIdToRound;
    mapping (uint256 => uint256[PLAYERS_PER_TEAM_MAX]) teamIdToPlayerIds;
    mapping (uint256 => address) teamIdToOwner;
    mapping (uint8 => uint256) tzToNCountries;



    uint16 _levelsInOneChallenge;
    uint16 _leafsInLeague;
    uint16 _levelsInLastChallenge;
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
    
}