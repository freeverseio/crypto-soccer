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
    mapping (uint256 => AuctionData) internal _playerIdToAuctionData;
    mapping (uint256 => AuctionData) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;

    uint256 internal nextVerseTimestamp;
    uint8 internal timeZoneForRound1;
    uint256 internal currentVerse;
    bytes32 internal currentVerseSeed;

    TimeZone[25] public _timeZones;
    uint256 public gameDeployDay;
    uint256 public currentRound;

    mapping (uint256 => uint256) countryIdToNDivisions;
    mapping (uint256 => uint256) countryIdToNHumanTeams;
    mapping (uint256 => uint256) divisionIdToRound;
    mapping (uint256 => uint256[PLAYERS_PER_TEAM_MAX]) teamIdToPlayerIds;
    mapping (uint256 => address) teamIdToOwner;
    mapping (uint8 => uint256) tzToNCountries;

    struct AuctionData {
        uint128 sellerHiddenPrice;
        uint32 validUntil;
    }

    struct TimeZone {
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