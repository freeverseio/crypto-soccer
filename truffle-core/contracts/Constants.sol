pragma solidity >=0.5.12 <=0.6.3;

/**
 * @title Constants used in the project
 */

contract Constants {
    
    address constant public NULL_ADDR = address(0);
    uint8 constant public PLAYERS_PER_TEAM_INIT = 18;
    uint8 constant public LEAGUES_PER_DIV = 16;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint256 constant public ACADEMY_TEAM = 1;

    // Skills: shoot, speed, pass, defence, endurance
    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;

    // Birth Traits: potential, forwardness, leftishness, aggressiveness
    uint8 constant private IDX_POT = 0;
    uint8 constant private IDX_FWD = 1;
    uint8 constant private IDX_LEF = 2;
    uint8 constant private IDX_AGG = 3;
    // prefPosition idxs: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    uint8 constant public IDX_GK = 0;
    uint8 constant public IDX_D  = 1;
    uint8 constant public IDX_M  = 2;
    uint8 constant public IDX_F  = 3;
    uint8 constant public IDX_MD = 4;
    uint8 constant public IDX_MF = 5;


    uint8 constant public TEAMS_PER_DIVISION = 128; // LEAGUES_PER_DIV * TEAMS_PER_LEAGUE
    uint256 constant public DAYS_PER_ROUND = 16;
    bytes32 constant INIT_ORGMAP_HASH = bytes32(0); // to be computed externally once and placed here

    uint8 constant internal IDX_MSG = 0;
    uint8 constant internal IDX_r   = 1;
    uint8 constant internal IDX_s   = 2;
    // POST_AUCTION_TIME: is how long does the buyer have to pay in fiat, after auction is finished.
    //  ...it includes time to ask for a 2nd-best bidder, or 3rd-best.
    uint256 constant public POST_AUCTION_TIME   = 6 hours; 
    uint256 constant public AUCTION_TIME        = 24 hours; 
    uint256 constant public MAX_VALID_UNTIL     = 30 hours; // the sum of the previous two
    uint256 constant internal VALID_UNTIL_MASK   = 0x3FFFFFFFF; // 2^34-1 (34 bit)
    uint8 constant public MAX_ACQUISITON_CONSTAINTS  = 7;


}