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

## Bits per game result

If we just care about who won => 2 bits (0=undef, 1=home, 2=away, 3=tie)

2 bit = who wins => 128 results in uint256

## Leagues goal average

If we assume that the maxGoals per game per team = 8, then we need

10 teams => 18 for one team => -144...144 => 9bit
11 teams => 20 for one team => ... 9bit
16 teams => 30 for one team => ... 9bit too!

So we can keep the goal avg of 28 teams in one uint :-)
Indeed, for 16teams, we still have 256-16x9 = 112bits left.

## Leagues 1 against 1 balance

To decide a winner of a league in case of tie of points, 
we need to look at the 2 respective games between teams that tie.
These are a total of nTeams x (nTeams-1)/2 values = half the number of games in a league.
We need 2 bits for each (0=undef, 1=team1, 2=team2, 3=tie)

10 teams => 45 values
11 teams => 55 values
16 teams => 120 values => 240 bit => fits in 1 uint


Note that when we update a team for the entire league, we compute the 2 relevant games
for every-other team, so it's enough to store the final value.

## If we wanted to store all results, we would need:

- if maxGoals = 14 (4 bits) => 1 game = 8bit =>
    => 10 teams = 90 x 8 / 256 = 3 uint256 
    => 16 teams = 240 x 8 / 256 = 8 uint256  (too much)


# Main Structs

## Players

Player is a struct that has:
    - string pname = player name, unique.
    - uint state = state0 + state1 + stateBit
        - state0/state1:serialization of skills. Currently we also store birthMonth, and role. 
        - stateBit: indicates which of the 2 states is the last updated state

We propose that role is taken out and made part of the league struct (as explained in another file). 

Assuming we keep birthMonth, how many states can we keep?

12bit => each state = 5x12= 60 => 4 states, since 12 + 4 x (5x12) = 252 bit ==> still 4 left!

(11 does not allow going to 5)

Let's keep 12bit => 4 states, plus 4 left.


## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + previousLeagueIdx ( +28 = 252)
                                + currentLeagueIdx ( +28 = 252)


- blockNumLastUpdate (29bit): he block number (29bit) of the last time an updater updated that team in that league (used to determine if challenge period is over or not). If this number is zero, it means that it has not been updated yet. If non-zero, it is used implicitly to know if the challenging time has passed, in case the team wants to join another league.

Rationale: when updating, the updater computes & writes the state of all players in team at the end of each league. It writes the blockNumLastUpdate.

The bool wasProcessedByBC is still 0. This, together witn nonzero blockNumLastUpdate: indicates that it can be challenged.

If a challeger finds it lies, then it re-updates the states of the players via a BC Tx. He sets the bool to 1, meaning that the team is totally ready to join another league, regardless of the value of blockNumLastUpdate. 

==> Note that we only need two playerstates: old (at the begininning of the league, and at the end). 


## League

Miners have the state of all teams at league join. 
Challengers can process at most one team through the league. For that, they
need to know that state of all players, and all teams at start of league.
These are not written anywhere, so they need to be provided as input by the 
challenger, and check that their hashes correspond to the previous leagues final hashes.

So the input of the challenger is nTeams x 14 uints, or since actually the state is 60bit,
and 4 fit in a 256, then we need nTeams x nPlayers/4 (e.g. 16 x 16/4 = 64 uints!).

One possibility is that these first re-write the player states, and the playGame works as usual. 

inputs needed for league creation: 
    - the teamIdx of each team that signs up.
    - the starting block
    - the number of blocks to wait between consecutive games

struct:
    - uint256 init: containing a serialization of (n0, nStep, nTeams, teamIdx_A for 7 teams)
    - uint256 teamIdx_B (space for 9 more teams)
    - address firstUpdaterAddr= who wrote the very first update for this league
    - address challengerUpdaterAddr = who wrote the last update
    - blockNumLastUpdate (28b) + wasUpdatedByBLockchain (18b)  
    - firstLeagueHash = MerkleRoot( initHash, finalHash[16], hash(results) )
    - challengerLeagueHash = same, written by last challenger
    - initHash = hash(  stateTeam1 + ... + stateTeamN )
    - finalHash[16] = [hash(stateTeam1), ..., hash(stateTeamN)]
    - uint results[8]

- bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)
- bits for wasUpdatedByBLockchain: 18 (1 bool per finalHash, for initHash, and LeagueHash)
- for results: 4bit (15 goals) x 2 x nTeams x (nTeams-1) = 7.5 * 256 => 8 uints



# Main Functions to control Updaters/Challengers

Updaters keep the entire state. We will disregard for the time being possible 
evolution of players between leagues. 

## isPlayer or isTeam or isLeague in challenging state, or safely updated?

these functions just look at player's current team, looks at current league, and checks if blockNumLastUpdate is:
    - 0 => not updated yet
    - less than needed => in challenge state
    - old enough => safely updated.



## ComputeLeague(leagueIdx) - can be called in view mode
- returns:
-- the LeagueHash, initHash, finalHash
-- the results of all games
-- the states of all players in all teams of the league


1. Computes all playGames for that team in the league (if some games were already updated, skips them) Doubt: how to pass the initial state of teams to playGame without actually writing it to each team.
2. Applies whatever evolution rules we have for each game
3. Computes all results, states, and finalHash.




## WriteLeagueByUpdater(leagueIdx, MerkleRoot)

1. writes the firstLeagueHash, firstUpdaterAddr, blockNumLastUpdate.

Obviously, wasUpdatedByBLockchain = 0...0 (16bit)


## ChallengeUpdateForFirstTime(leagueIdx, initHash, finalHash[16], results[8])
0. checks if league is in challenging period, otherwise halts.
1. checks if league is challenged for first time, otherwise halts.
2. computes challengerLeagueHash from input data.
2. Just writes initHash, finalHash[16], and results[8]. 
3. Writes challengerUpdaterAddr, blockNumLastUpdate. Stops.

Still, wasUpdatedByBLockchain = 0...0 (16bit)


## ChallengeUpdate(leagueIdx, MerkleProofsInit[16][5], allLeaguePlayersStates, selectedTeam, challengerLeagueHash)
0. checks if league is in challenging period, otherwise halts.
1. checks if league is in challenge > 1 and that selectedTeam wasUpdatedByBLockchain, otherwise halts.
2. If wasUpdatedByBLockchain for initHash = 1, skip to 6.
3. For each team in league: 
    - computes hash, 
    - goes to previousLeagueIdx, uses MerkleProofsInit to show the hash is correct,
        if any is not: halts
4. Computes initHash, sets wasUpdatedByBLockchain for initHash to 1, and set blockNumLastUpdate
5. If initHash not equal to what is written: rewrite, halt (no need to prove more)
6. Plays all league (many playGame) for selectTeam 
    - either by writing allLeaguePlayersStates first (costly)
    - ideally, by just using the call data
7. writes ONLY the results computed, and 
    - note that some results will be written twice when new challenger comes, but must coincide
8. writes ONLY the finalHash[selectedTeam]
9. sets to 1 the wasUpdatedByBLockchain for that team
10. If BOTH results and hash were correct: halt
11. Get the money from previous updater (not from firstUpdater)
12. if challengerLeagueHash coincides with firstLeagueHash, halt.
13. write new challengerLeagueHash.


If rounds arrive to nTeams (max possibility), the LeagueHash is computed automatically.

At any point, when the blockNumLastUpdate is old enough:

- either challengerLeagueHash = firstLeagueHash => first updater's deposit is released.
- otherwise: last challenger gets first updater's deposit too.

In all cases, the challenger gets (at least) the deposit from all previous challengers!






doubt if the first challenger actually provides a good input that has the same firstLeagueHash as the first updater (he's silly)... what happens? write something for this?



