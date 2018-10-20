/// The main set of constants
k = {
    RoleKeeper : 0,
    RoleDef : 1, 
    RoleMid : 2, 
    RoleAtt : 3, 

    /// Likewise for enum State { Birth, Def, Speed, Pass, Shoot, End, Role }
    StatBirth : 0, 
    StatDef : 1, 
    StatSpeed : 2, 
    StatPass : 3, 
    StatShoot : 4, 
    StatEndur : 5, 
    StatRole : 6,

    /// Summarize: how many states, and from these, how many are skills: 
    NumStates : 7, 
    NumSkills : 5, 

    /// Ennum for globSkills: [0-move2attack, 1-createShoot, 2-defendShoot, 3-blockShoot, 4-currentEndurance]
    Move2Attack : 0, 
    CreateShoot : 1, 
    DefendShoot : 2, 
    BlockShoot : 3, 
    Endurance : 4, 

    /// @dev The number of rounds in a given game (18 rounds = 1 event per 5 mins)
    RoundsPerGame : 18,

    /// The amount of bits used per state to serialize them in a uint256 
    BitsPerState : 14, 

    /// The amount of bits used per state to playerIdx them in a uint256 
    BitsPerPlayerIdx : 20, 

    /// The number of bits used for each rndNumber used to determine each game action
    BitsPerRndNum : 14,

    /// Max num of players allowed in a team
    MaxPlayersInTeam : 11,

    /// Vals used to store game results (0=undefined, 1=home wins, 2=away wins, 3=tie)
    Undef : 0,
    HomeWins : 1,
    AwayWins : 2,
    Tie : 3
}
module.exports = k;
