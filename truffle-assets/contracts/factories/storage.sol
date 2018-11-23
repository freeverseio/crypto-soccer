pragma solidity ^ 0.4.24;

/*
    Defines all storage structures and mappings
*/

contract Storage {

    /// @dev The main Player struct.
    /// @dev name is a string, unique for every Player
    /// @dev state is a uint256 that serializes age, skills, role.
    /// @dev Each skill is sent as a uint16 for serialization, occupying 20 bits of the state
    /// @dev The order of elements serialized are:
    ///         0-monthOfBirthAfterUnixEpoch; if this goes up to 9999, then the game will run for at least 800 more years.
    ///         1-defense
    ///         2-speed
    ///         3-pass
    ///         4-shoot (for a goalkeeper, this is interpreted as ability to block a shoot)
    ///         5-endurance
    ///         6-role
    struct Player {
        string name;
        uint state;
    }
    /// @dev An array containing the Player struct for all players in existence. 
    /// @dev The ID of each player is actually his index this array.
    Player[] players;


    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
        uint256 playersIdx;
    }

    /// @dev An array containing the Team struct for all teams in existence. 
    /// @dev The ID of each team is actually his index in this array.
    Team[] teams;

    /// @dev The main League struct
    struct League {
        uint[] teamIdxs;
        uint blockFirstGame;
        uint blocksBetweenGames;
        uint resultsFirstHalf;
        uint resultsSecondHalf;
    }
    /// @dev Array containing all leagues created so far
    League[] leagues;


    /// @dev A mapping from hash(playerName) to a Team struct.
    /// @dev Facilitates checking if a playerName already exists.
    mapping(bytes32 => Team) public playerToTeam;

    /// @dev A mapping from team hash(name) to the owner's address.
    /// @dev Facilitates checking if a teamName already exists.
    mapping(bytes32 => address) public teamToOwnerAddr;

    /// @dev Upong deployment of the game, we create the first null player
    /// @dev Choose a silly serialized state (meaningless age, skills, etc)
    /// @dev to differentiate it from 0.
    constructor() public {
        players.push(Player({name: "_", state: uint(-1) }));
    }


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
