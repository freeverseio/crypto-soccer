pragma solidity ^ 0.4.24;

contract CryptoSoccer {
    /// CONSTANTS SECTION
    /// @dev Instead of Enums, we use consts. Enums cannot be casted explitcly!
    /// @dev Instead, consts are truly replaced by their value at compile time. 
    /// @dev So to emulate: enum Role { Keeper, Def, Mid, Att, Subst, Retired }, we do:
    uint8 constant kRoleKeeper = 0;
    uint8 constant kRoleDef = 1; 
    uint8 constant kRoleMid = 2; 
    uint8 constant kRoleAtt = 3; 

    /// @dev Likewise for enum State { Birth, Def, Speed, Pass, Shoot, End, Role }
    uint8 constant kStatBirth = 0; 
    uint8 constant kStatDef = 1; 
    uint8 constant kStatSpeed = 2; 
    uint8 constant kStatPass = 3; 
    uint8 constant kStatShoot = 4; 
    uint8 constant kStatEndur = 5; 
    uint8 constant kStatRole = 6;

    /// @dev Summarize: how many states, and from these, how many are skills: 
    uint8 constant kNumStates = 7; 
    uint8 constant kNumSkills = 5; 

    /// @dev Ennum for globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance]
    uint8 constant kMove2Attack = 0; 
    uint8 constant kCreateShoot = 1; 
    uint8 constant kDefendShoot = 2; 
    uint8 constant kBlockShoot = 3; 
    uint8 constant kEndurance = 4; 

    /// @dev The number of rounds in a given game (18 rounds = 1 event per 5 mins)
    uint8 constant kRoundsPerGame = 18; 

    /// @dev The number of bits used for each rndNumber used to determine each game action
    uint8 constant kBitsPerRndNum = 14; 
    uint16 constant kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 

    /// @dev The amount of bits used per state to serialize them in a uint256 
    uint8 constant kBitsPerState = 14; 

    /// @dev The amount of bits used per state to playerIdx them in a uint256 
    uint8 constant kBitsPerPlayerIdx = 20; 

    /// @dev Max num of players allowed in a team
    uint8 constant kMaxPlayersInTeam = 11;

    /// @dev The amount of bits used game result in a uint256 
    uint8 constant kBitsPerGameResult = 2; 

    /// @dev Vals used to store game results (0=undefined, 1=home wins, 2=away wins, 3=tie)
    uint8 constant kUndef = 0;
    uint8 constant kHomeWins = 1;
    uint8 constant kAwayWins = 2;
    uint8 constant kTie = 3;
}
