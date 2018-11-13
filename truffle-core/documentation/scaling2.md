# Numerology

## CKitties stats: see https://etherscan.io/token/0x06012c8cf97bead5deae237070f9587f8e7a266d

Exisiting cats: 1,160,408, 
Different owners: 68698 addresses
Total transfers: 3,068,744

## Bits to store block numbers? 

About 1 block is produced every 10s. 
    - about 3M blocks per year: we would need 22 bit
    - about 8K blocks per day: we would need 13 bit
    - if we want the game to last 100 years => 29bit
    - if we want the time between games to be 1 week => 16 bit


## Bits to store monthOfBirthAfterUnixEpoch (i.e. after 1970)

10bit = 85 years
12bit = 341 years  
13bit = 682 years  (ENOUGH)

## Bits to store player state (they start with an avg of 50 points per skill)

10bit = 1024
12bit = 4098
13bit = 16386


## Bits to store playerIdx (i.e. what's the max num of players that can be created?)

Currently: 
    20bit = 1M players ==> 12 players in a uint256. 
We should use more:
    28 = 270M players ==> 9 players in a uint256. 

## Bits to store laegueIdx (i.e. what's the max num of leagues that can be created?)

Limiting this allows to keep track of many leagues played without updating.
25bit = 33M leagues.


## Bits to store teamIdx (i.e. what's the max num of teams that can be created?)

25bit = 33.5M teams (good number since we can squeeze 10 in a uint256, and use the remaining 6 bits to say how many teams are to be read 2^6 = 64)

26bit = 67M teams (allows 9 teams per uint256)
27bit = 134M teams (allows 9, 27x9=234)
28bit = 268 teams (allows 9, 28x9=252)

## Leagues
- Games played by a team in a league: nGamesPerTeam = 2 (nTeams - 1)
- Total games that make a league: nTeams (nTeams-1)    (a square minus de diagonal)

10 teams => 90 games overall, 18 for one team
11 teams => 110 games, 20 for one team
16 teams => 240 games, 30 for one team

## Max games that can be computed in one block

Max gas per block = 8M.
Max gas for a reasonable large Tx = 4M.
Currently, 1 game = 260K. Let it be 300K.

18 games = 5.4M  => so a league of 10 seems like the upper limit if we want to allow challenging in one atomic Tx.

If we allow challenging in 2 rounds, then it's fine to use leagues of 16 teams.

LETS LIMIT TO LEAGUES OF 10 for the time being.

## Bits per game result

max goals per team per game = 15 (4bit)

if nTeams = 10
4 x 2 x 10 x (10-1) = 7.5 * 256 => 3 uints

if nTeams = 16
4 x 2 x 16 x (16-1) = 7.5 * 256 => 8 uints



# Main Structs

## Players

Player is a struct that has:
    - string pname = player name, unique.
    - uint state = serialization of skills(70bit) 
                    + currentTeamIdx 
                    + prevTeamIdx


Currently we also store birthMonth, and role as skills (revisar). We propose that role is taken out and made part of the league struct (as explained in another file). 



## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + previousLeagueIdx ( +28 = 252)
                                + currentLeagueIdx ( +28 = 252)



## League

Updaters have the state of all teams at league join. 
Challengers can process at most one team through the league. For that, they
need to know that state of all players, and all teams at start of league.
These are not written anywhere, so they need to be provided as input by the 
challenger, and check that their hashes correspond to the previous leagues final hashes.

struct:
    Provided by a user creating the league (and others that join):
    - uint256 init: containing a serialization of (n0, nStep, nTeams, teamIdx_A for 7 teams)
    - uint256 teamIdx_B (space for 9 more teams)

    Provided by first updater:
    - address firstUpdaterAddr = who wrote the very first update for this league
    - firstMerkleRoot = MerkleRoot( initHash, finalHash[16], hash(results) )
    - blockNumLastUpdate (28b) + wasUpdatedByBLockchain (18b, initially all zero)  

    Provided by challengers:
    - address challengerAddr = who wrote the last challenge
    - challengerMerkleRoot 
    - initHash = hash(  hash(stateTeam1) + ... + hash(stateTeamN) )
    - finalHash[16] = [hash(stateTeam1), ..., hash(stateTeamN)]
    - uint results[8]   

Some data on bits used:

- bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)
- bits for wasUpdatedByBLockchain: 18 (1 bool per finalHash, for initHash, and LeagueHash)

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




## WriteLeagueByUpdater(leagueIdx, MerkleRoot)

1. writes the firstMerkleRoot, firstUpdaterAddr, blockNumLastUpdate.

Obviously, does not touch wasUpdatedByBLockchain = 0...0 (16bit)



## ChallengeUpdateForFirstTime(leagueIdx, initHash, finalHash[16], results[8])

0. checks if league is in challenging period, otherwise halts.
1. checks if league was challenged for first time, otherwise halts.
2. computes challengerMerkleRoot from input data. If equal to firstMerkleRoot, halts.
2. Writes initHash, finalHash[16], results[8], challengerUpdaterAddr, firstMerkleRoot.
3. Re-writes blockNumLastUpdate.

Still, wasUpdatedByBLockchain = 0...0 (16bit)




## ChallengeUpdate(leagueIdx, finalHash[16], results[8], MerkleProofsInit[16][5], allLeaguePlayersStates, selectedTeam, challengerLeagueHash)

0.  - require league is in challenging period.
    - require league is in challenge > 1 (e.g. firstMerkleRoot != 0), and that 
    wasUpdatedByBLockchain for selectedTeam == 0.

1. Computes hash(team.playerStates) for all teams, and initHash.

2. If the bit wasUpdatedByBLockchain for initHash = 1, compare your initHash to the written one. 
    - If different, halt.
    - If equal, jump to Step 6.

3. Prove that proposed computed initHash is correct by proving that each computed hash(team.playerStates) is contained the the previous team's league Merkle tree, by:    
    For each team in league: 
    - goes to team.previousLeagueIdx, uses MerkleProofsInit[team][:] to show the hash was contained in it. 
    - if any of this is not true for a team: halts     

4. Set wasUpdatedByBLockchain for initHash to 1, 

5. If initHash not equal to what is written (you've already proven he was lying):
    - rewrites initHash, finalHash[16], results[8], challengerUpdaterAddr, firstMerkleRoot.
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

11 Transfer the deposit from previous updater (not from firstUpdater) to challengerAddr.

12. require challengerMerkleRoot != firstMerkleRoot, otherwise halt: (firstUpdater did not lie, his deposit will not be taken)

13. write new challengerMerkleRoot.

14. set blockNumLastUpdate

If rounds arrive to nTeams (max possibility), the MerkleRoot is computed by the BC.

At any point, when the blockNumLastUpdate is old enough:

- either challengerMerkleRoot = firstMerkleRoot => first updater's deposit is released.
- otherwise: last challenger gets first updater's deposit too.

In all cases, the challenger gets (at least) the deposit from all previous challengers!





