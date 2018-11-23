# Mappings
    - mapping(bytes32 => address) public teamToOwnerAddr; /// from team hash(name) to the owner's address.
    - mapping(uint => state) public playerIdxToState; 
    - mapping(uint => address) public playerIdxToOwnerAddr; 



# Numerology

## Bits to store block numbers? 

About 1 block is produced every 10s. 
    - about 3M blocks per year: we would need 22 bit
    - about 8K blocks per day: we would need 13 bit
    - if we want the game to last 100 years => 29bit
    - if we want the time between games to be 1 week => 16 bit


## Bits to store monthOfBirthAfterUnixEpoch (i.e. after 1970)

14bit = 1364 years  (ENOUGH)

## Bits to store player skills (they start with an avg of 50 points per skill)

14bit = 32K values per skill!
5 skills => 70bit


## Bits to store playerIdx (i.e. what's the max num of players that can be created?)

28 = 270M players ==> 9 players in a uint256. 


## Bits to store laegueIdx (i.e. what's the max num of leagues that can be created?)

25bit = 33M leagues.


## Bits to store teamIdx (i.e. what's the max num of teams that can be created?)

28bit = 268 teams (allows 9 per uint, 28x9=252)


## Leagues
- Games played by a team in a league: nGamesPerTeam = 2 (nTeams - 1)
- Total games that make a league: nTeams (nTeams-1)    (a square minus de diagonal)

references:
10 teams => 90 games overall, 18 for one team
11 teams => 110 games, 20 for one team
16 teams => 240 games, 30 for one team



## Max games that can be computed in one block

Max gas per block = 8M.
Max gas for a reasonable large Tx = 4M.
Currently, 1 game = 260K. Let's say 300K.

18 games = 5.4M  => so a league of 10 seems like the upper limit if we want to allow challenging in one atomic Tx.

LETS LIMIT TO LEAGUES OF 10 for the time being.


## Bits per game result

if we limit max goals per team per game to 15 => use 4bit per team per game

if nTeams = 10
4 x 2 x 10 x (10-1) = 7.5 * 256 => 3 uints




# Main Structs/Data

## Player state

    - uint state = serialization of:
        - skills: 5 x 14 = 70bit
        - monthOfBirthAfterUnixEpoch: 14bit
        - lastRole in a team: 2bit 
        - currentTeamIdx: 28bit



## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + previousLeagueIdx ( +28 = 224)
                                + currentLeagueIdx ( +28 = 252)



## League

Updaters have the state of all teams at league join, built from genesis.
Challengers can process at most one team through the league. For that, they
need to know the state of all players, and all teams, at start of league.
These are not written anywhere, so they need to be provided as input in the 
challenger's TX, and the contract will check that their hashes are contained
in the previous leagues final hashes.

struct:
    Provided by a user creating the league (and others that join):
    - uint256 init: containing a serialization of (n0, nStep, nTeams, teamIdx_A for 7 teams)
    - uint256 teamIdx_B (space for 9 more teams)

    Provided by first updater:
    - initHash = hash(  hash(stateTeam1) + ... + hash(stateTeamN) )
    - finalHash[10] = [hash(stateTeam1), ..., hash(stateTeamN)]
    - uint results[3]   
    - address updaterAddr = who wrote the update for this league
    - blockNumLastUpdate (28b) + wasUpdatedByBLockchain (12, initially all zero)  
    - MerkleRoot( L )  = computed from L


NOMENCALTURE:

Let L be the collection of league data that needs to be provided so that one can prove that I lied:

    L =  {initHash, finalHash[10], results[3]}

Let MERKLE be the MerkleRoot of the 12 entries (note the hash in the last piece):

    MERKLE = MerkleRoot( {initHash, finalHash[10], hash(results[3]) } )

Note that to proof that one of these entries is part of MERKLE, we need 3 accompanying siblings.


Some data on bits used:

- bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)
- bits for wasUpdatedByBLockchain: 12 (1 bool per finalHash, for initHash, and LeagueHash)

- for results: 4bit (15 goals) x 2 x nTeams x (nTeams-1) 
                            (for nTeams = 16)= 7.5 * 256    = 8 uints
                            (for nTeams = 10)               = 3 uints

- on initHash: just concatenate all playerStates of a teams, and hash. Concatenate all hashes, and hash. The reason why we include hashes instead of concatenating all playerStates of all teams, is that each hash has to be done anyway to prove that it's in the previous-league Merkle tree, and it sounds cheaper.




# Main Functions to control Updaters/Challengers


Updaters keep the entire state. We will disregard for the time being possible 
evolution of players between leagues. 


## isPlayer or isTeam or isLeague in challenging state, or safely updated?

- isPlayerBusyPlayingLeague?
Player.currentTeam -> team.currentLeague -> n + nGames x deltaN < now ?

- isPlayerStateUpdated?
Player.currentTeam -> team.currentLeague -> blockNumLastUpdate != 0 ?

- isPlayerInChallengingPeriod?
Player.currentTeam -> team.currentLeague -> blockNumLastUpdate + challengePeriod < now ?

- isPlayerReadyForNewStuff = isPlayerStateUpdated && !isPlayerInChallengingPeriod

Same for teams or leagues



## WriteLeagueByUpdater(leagueIdx, L)

He provides the collection of data L, as defined above. The contract:

0. checks if league is in challenging period, otherwise halts.
1. Writes initHash, finalHash[10], results[3], challengerUpdaterAddr, MerkleRoot, blockNumLastUpdate

Obviously, does not touch wasUpdatedByBLockchain = 0...0 (16bit)



## ChallengeUpdate(leagueIdx, L, MerkleProofsInit[10][3], allLeaguePlayersStates, selectedTeam)

0.  - require league is in challenging period.
    - require wasUpdatedByBLockchain for selectedTeam == 0.
    - compute MerkleRoot from input data, require it is different from existing MerkleRoot

1. Computes hash(team.playerStates) for all teams, and then uses it to compute initHash.

2. If the bit wasUpdatedByBLockchain for initHash = 1, compare your initHash to the written one. 
    - If different, halt.
    - If equal, jump to Step 6.

3. Prove that proposed computed initHash is correct by proving that each computed hash(team.playerStates) is contained the the previous team's league Merkle tree, by:    
    For each team in league: 
    - goes to team.previousLeagueIdx, uses MerkleProofsInit[team][:] to show the hash was contained in it. 
    - if any of this is not true for a team: halts     

4. Set wasUpdatedByBLockchain for initHash to 1, 

5. If initHash not equal to what is written (you've already proven he was lying):
    - rewrites initHash, finalHash[10], results[3], challengerUpdaterAddr, firstMerkleRoot.
    - set blockNumLastUpdate
    - halts

6. Plays all league (many playGame) for selectTeam, based on the provided allLeaguePlayersStates in this TX call.

7. sets to 1 the wasUpdatedByBLockchain for that team

8. Require that either hash(selectedTeam) or the results of the computed games for selected team are different from those written previously (proving that they previously lied). If equal: 
    - set to 1 the wasUpdatedByBLockchain for that team,... 
    - set blockNumLastUpdate
    - and HALT (basically, you're silly, but we'll use your money to set to 1 that bit)


9. write ONLY the results for the computed games (for selectedTeam). Note that some results may be written twice when new challenger comes, but must coincide.

10. write ONLY the finalHash[selectedTeam]

11. Slash: Transfer the deposit from previous updater/challenger to challengerAddr.

12. rewrite new MerkleRoot.

13. rewrite blockNumLastUpdate.

If rounds arrive to nTeams (max possibility), the MerkleRoot has effectively been computed by the BC.

At any point, when the blockNumLastUpdate is old enough, the state is settled.




