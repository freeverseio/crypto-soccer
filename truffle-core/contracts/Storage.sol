pragma solidity >=0.5.12 <0.6.2;

/**
* @title Storage common to all project, with setters managed by StorageProxy.
*/
contract Storage {

    uint8 constant private PLAYERS_PER_TEAM_MAX  = 25;

    uint256[2**8] _slotReserve;
    address internal _storageOwner; // TODO: move to a "proposed new owner" + "accept" instead of stright "set net owner"
    address internal _academyAddr;
       
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => ContractInfo) internal _selectorToContractInfo;
    
    
    mapping(uint256 => uint256) internal _playerIdToState;
    mapping (uint256 => uint256) internal _playerIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToAuctionData;
    mapping (uint256 => uint256) internal _teamIdToRemainingAcqs;

// todo: add NULL_TZ = 0 and use it

    TimeZone[25] public _timeZones;  // timeZone = 0 is a dummy one, without any country. Forbidden to use timeZone[0].
    uint256 public gameDeployDay;
    uint256 public currentRound;

    struct ContractInfo {
        address addr;
        bool requiresPermission;
        bytes32 name;
        bytes4[] selectors;
    }

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