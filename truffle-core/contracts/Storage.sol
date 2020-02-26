pragma solidity >=0.5.12 <=0.6.3;

import "./ProxyStorage.sol";
import "./Constants.sol";

/**
* @title Storage common to all project, with setters managed by StorageProxy.
*/
contract Storage is ProxyStorage, Constants{

    uint256[2**8] _slotReserve;
   
    address internal _academyAddr;
    
    mapping(uint256 => uint256) internal _playerIdToState;
    mapping (uint256 => uint256) internal _playerIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;

    uint256 internal nextVerseTimestamp;
    uint8 internal timeZoneForRound1;
    uint256 internal currentVerse;
    bytes32 internal currentVerseSeed;


// todo: add NULL_TZ = 0 and use it

    TimeZone[25] public _timeZones;  // timeZone = 0 is a dummy one, without any country. Forbidden to use timeZone[0].
    uint256 public gameDeployDay;
    uint256 public currentRound;

    struct Team {
        uint256[PLAYERS_PER_TEAM_MAX] playerIds; 
        address owner;
    }

    struct Country {
        uint256 nDivisions;
        uint8 nDivisionsToAddNextRound;
        mapping (uint256 => uint256) divisonIdxToRound;
        mapping (uint256 => Team) teamIdxInCountryToTeam;
        uint256 nHumanTeams;
    }

    struct TimeZone {
        Country[] countries;
        uint8 nCountriesToAdd;
        bytes32[2] orgMapHash;
        bytes32[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        bytes32 scoresRoot;
        uint8 updateCycleIdx;
        uint256 lastActionsSubmissionTime;
        uint256 lastUpdateTime;
        bytes32 actionsRoot;
        uint256 lastMarketClosureBlockNum;
    }    


}