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

    /// The amount of bits used per state to serialize them in a uint256 
    BitsPerState : 14, 

    /// The amount of bits used per state to playerIdx them in a uint256 
    BitsPerPlayerIdx : 20, 

    /// Max num of players allowed in a team
    MaxPlayersInTeam : 11,
}
module.exports = k;
