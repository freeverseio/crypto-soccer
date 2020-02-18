pragma solidity >=0.5.12 <0.6.2;

/**
* @title Storage common to all project, with setters managed by StorageProxy.
*/
contract Storage {

    uint8 constant private PLAYERS_PER_TEAM_MAX  = 25;

    uint256[2**16] _slotReserve;
    address internal _storageOwner; // TODO: move to a "proposed new owner" + "accept" instead of stright "set net owner"
    bytes4[] internal _allFunctions;
       
    ContractInfo[] internal _contractIdToInfo;
    mapping (bytes4 => uint256) internal _functionToContractId;

    mapping(uint256 => uint256) internal _playerIdToState;

    TimeZone[25] public _timeZones;  // timeZone = 0 is a dummy one, without any country. Forbidden to use timeZone[0].
    uint256 public gameDeployDay;
    uint256 public currentRound;
    bool public _wasInited;

    struct ContractInfo {
        address addr;
        string name;
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