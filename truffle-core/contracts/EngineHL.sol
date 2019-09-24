pragma solidity ^0.5.0;

import "./Encoding.sol";
import "./Engine.sol";

contract EngineHL is Encoding{
    
    uint8 public constant ROUNDS_PER_MATCH = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)
    uint8 private constant BITS_PER_RND     = 36;   // Number of bits allowed for random numbers inside match decisisons
    uint256 public constant MAX_RND         = 68719476735; // Max random number allowed inside match decisions: 2^36-1
    uint256 public constant MAX_PENALTY     = 10000; // Idx used to identify normal player acting as GK, or viceversa.
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    uint256 private constant TENTHOUSAND    = uint256(10000); 
    uint256 private constant TENTHOUSAND_SQ = uint256(100000000); 
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0

    Engine private _engine;

    function playMatch(
        address engineAddr,
        uint256 seed,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        bool[2] memory matchBools  // [is2ndHalf, isHomeStadium]
    )
        public
        pure
        returns (uint8[2] memory teamGoals) 
    {
        uint8[11][2] memory lineups;
        uint8[9][2] memory playersPerZone;
        bool[10][2] memory extraAttack;
        (lineups[0], extraAttack[0], playersPerZone[0]) = getLineUpAndPlayerPerZone(tactics[0]);
        (lineups[1], extraAttack[1], playersPerZone[1]) = getLineUpAndPlayerPerZone(tactics[1]);

        uint16[5][11][2] memory playerSkills;
        uint8[11][2] memory forwardness;
        uint8[11][2] memory leftishness;
        
        (playerSkills[0], forwardness[0], leftishness[0]) = getLinedUpSkills(states[0], lineups[0]);
        (playerSkills[1], forwardness[1], leftishness[1]) = getLinedUpSkills(states[1], lineups[1]);
        
        return Engine(engineAddr).playMatch(seed, playerSkills, extraAttack, forwardness, leftishness, playersPerZone, matchBools);
    }



    
    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    function getLineUpAndPlayerPerZone(uint256 tactics) 
        public 
        pure 
        returns (uint8[11] memory lineup, bool[10] memory extraAttack, uint8[9] memory playersPerZone) 
    {
        uint8 tacticsId;
        (lineup, extraAttack, tacticsId) = decodeTactics(tactics);
        return (lineup, extraAttack, getPlayersPerZone(tacticsId));
    }
    
    // TODO: can this be expressed as
    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    function getPlayersPerZone(uint8 tacticsId) internal pure returns (uint8[9] memory) {
        require(tacticsId < 4, "we currently support only 4 different tactics");
        if (tacticsId == 0) return [1,2,1,1,2,1,0,2,0];  // 0 = 442
        if (tacticsId == 1) return [1,3,1,1,2,1,0,1,0];  // 0 = 541
        if (tacticsId == 2) return [1,2,1,1,1,1,1,1,1];  // 0 = 433
        if (tacticsId == 3) return [1,2,1,1,3,1,0,1,0];  // 0 = 451
    }


    function getLinedUpSkills(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState, 
        uint8[11] memory lineup
    )
        public
        pure
        returns (
            uint16[5][11] memory linedUpSkills,
            uint8[11] memory forwardness,
            uint8[11] memory leftishness
        )
    {
        for (uint8 p = 0; p < 11; p++) {
            require(lineup[p] < PLAYERS_PER_TEAM_MAX, "lineup val too large");
            require(getPlayerIdFromSkills(teamState[lineup[p]]) != FREE_PLAYER_ID, "cannot line up an empty player");            
            linedUpSkills[p] = getSkillsVec(teamState[lineup[p]]);
            forwardness[p] = getForwardness(teamState[lineup[p]]);
            leftishness[p] = getLeftishness(teamState[lineup[p]]);
        }
    }

}

