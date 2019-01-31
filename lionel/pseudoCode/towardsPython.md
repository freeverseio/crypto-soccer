
userData entried for league, between games

TACTICS

3 bits => 8 possibilitats (443 ....) 
1 bits => pressing or not

4 bits, 2 equips => 1 byte per partit

tactics byte[90]


league.tactics[90]

-------------------

443

join league: (escriure teamIdx) sorted players

player 0: goalie
from more defensive to more attacking

join league: user sends shirtNums[14]

--------------------

league struct:

- teamIdx[10], shirtNums[10][14]  at start of league
- tactics[90] as the league goes on

userDataAtInit = Hash(  teamIdx[10], shirtNums[10][14]  )
userDataAlongLeague = Hash( Hash( Hash( ... tactics))

function updateStrategy(newStrategy) {
    stateHash = keccak256(newStrategy,statehash);
}



user sends new tactic => smart contract H(new, before)

userDataHash = H( userDataAtInit, userDataAtInitHash )


minime => permet guardar el saldo en un block anterior
blockhashes are now stored from solidity (Constantinople) 
geth => gimme Merkle proof of this data in this block

in a smart contract you can now prove smtg existed at some block


https://twitter.com/izqui9/status/1068575567614746625
https://github.com/3box/3box

## TODO




## TODO

- Challanging game       (lionel.md, uses v3)
- Data structure (as in linonel.md)
    - player
    - team
    - league
- Challenging incentives (adria post)

### DataStructure 

```
struct Player {
    skills[5],                  # 5 x 14 = 70bit
    monthOfBirthAfterUnixEpoch, # 14bit
    currentTeamIdx,             # 28bit // the current team of the player
    playerIdx,                  # 28bit // player id
    lastSaleBlocknum            # 29bit
}

struct Team {
    string tname                 # = unique
    playersIdx[14],              # Player.playerIdx
    previousLeagueIdx,           # ( +28 = 224)
    currentLeagueIdx #           # ( +28 = 252)
}

struct League {
    blockNumStart,               # starting block for this league
    blockDelta,                  # blocks between plays
    userInitDataHash,            # contains teamsIdxs[10],shirtNums[10][14] 
    userDuringLeagueDataHash,    # contains tactics[90]
    
    # challange/updater game data
    
    initHash                     # contains Hash(allplayersStates) at init
    finalTeamStatesHash[10],     # contains Hash(allplayersStates) at end
    scores[90],
    updaterAddr,
    blockNumLastUpdate
}
```




## UpdateLeague(leagueIdx, D)

user_U1_create_team
    user provides teamName + int ('Barcelona','112')
            - hash('Barcelona','112') is random seed that determines all players
            - teamIdx is assigned (by appending to teams[])
            - playerIdxs are reserved: teamIdx*14... teamIdx*14 + 13
    
    
user_U2_create_team

user_U1_enrolls_legue_L1

user_U2_enrolls_legue_L1

game1_play

user_U1_changes_tactics_to_443

game2_play()

updater_computes_finalData_and_calls_update_with_lie

challanger_proves_updater_is_lying

updater_computes_finalData_and_calls_update_without_lying

challanging_window_is_closed




user_U1_plays_game_G1

user_U2_plays_game_G1

user_U1_plays_game_G2

user_U2_plays_game_G2




He provides the collection of data D, as defined above. The contract:

0. checks if league is in challenging period, otherwise halts.
1. Writes initHash, finalTeamStatesHash[10], results[3], challengerUpdaterAddr, blockNumLastUpdate



### Challanging incentives

```
updater
    var hashonions : Hash[1000]
    var owner : address
    init()
        hashonions[999]=rand()
        for i=998..0 hashonions[i]=hash(hashonions[i+1])
        smartcontract.register(owner,1000 dai,hashonions[0])
        hashonions.pop() // remove top
        
    when smartcontract.isOkForUPDATE( owner )
        leagueId = rand(1..smartcontract.league())
        if hashonions.peek() % 10 == 0 // first element
            smartcontract.update( leagueId, random_stuff )
        else    
            result = smartcontract.simulatePlay(leagueId)
            smartcontract.update( leagueId, result )
            
    when smartcontract.isOkForREVEAL( owner )
        smartcontract.revealHash( hashonions.peek() )
        hashonions.pop()
    
    when i_not_participated_in_one_month
        smartcontract.burnGas( )
        
    when i_want_to_quit
        smartcontract.queryQuit()
        
    when smartcontract.isOkForQUIT( owner ) 
        smartcontract.quit() // get 1000 dai back
             
challanger
    var owner : address
    init()
        smartcontract.register(owner,1000 dai)
   
    when smartcontract.isOkForCHALLANGE( owner )
        leagueId = smartcontract.getUpdatedAndNotChallangedLeagueId()
        current = smartcontract.getUpdate(leagueId)
        expected = smartcontract.simulatePlay(leagueId)
        if current != expected
            smartcontract.challange( leagueId, expected )
        
    when i_not_participated_in_three_months
        smartcontract.burnGas( )
                
    when i_want_to_quit
        smartcontract.quit() // get 1000 dai back
```






























