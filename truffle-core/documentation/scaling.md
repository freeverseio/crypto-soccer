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
We will limit the maxplayers in a team to 14. Using 28b for each playerIdx, we can do:
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
    - uint256 playersIdxB = serialization of 5 playerIdx + 4 leagueIdx + updateState(2bit)

- in playersIdxB: 256 - 5 x 28 = 116...
- ...so we can fit up to 4 leagueIdx of 28bit each + updateState: 4 x 28 + 2 = 114.

updateState: (0 = no news; 1 = it is in challenge period, 2 = )... TODO FINISH

The 4 leagueIdx, and updateState, start at 0.
When team joins a league, the first entry if filled, etc.
If the 4 entries are filled, the team cannot join another league without being updated.


are set to 0 when an update happens. So we know the leagues that we need
to go through to update a team.

## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + blockNumLastUpdate ( + 28b = 224)
                                + currentLeagueIdx ( +28 = 252)
                                + wasProcessedByBC ( +1 = 253)
    - address lastUpdaterAddr = who wrote the last update


- blockNumLastUpdate (29bit): he block number (29bit) of the last time an updater updated that team in that league (used to determine if challenge period is over or not). If this number is zero, it means that it has not been updated yet. If non-zero, it is used implicitly to know if the challenging time has passed, in case the team wants to join another league.

Rationale: when updating, the updater computes & writes the state of all players in team at the end of each league. It writes the blockNumLastUpdate.

The bool wasProcessedByBC is still 0. This, together witn nonzero blockNumLastUpdate: indicates that it can be challenged.

If a challeger finds it lies, then it re-updates the states of the players via a BC Tx. He sets the bool to 1, meaning that the team is totally ready to join another league, regardless of the value of blockNumLastUpdate. 

==> Note that we only need two playerstates: old (at the begininning of the league, and at the end). 


## League

inputs needed for league creation: 
    - the teamIdx of each team that signs up.
    - the starting block
    - the number of blocks to wait between consecutive games

struct:
    - uint256 init: containing a serialization of (n0, nStep, nTeams, teamIdx_A)
    - (optional) uint256 teamIdx_B (space for 9 more teams)
    - uint256 resultsFirstHalf
    - uint256 resultsSecondHalf
    - uint256 goalAverages
    - uint256 oneOnOneBalance


- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)
- bits per result: 2 (0=undef, 1=home, 2=away, 3=tie)
- bits per goalAvg: 9 (max -256...256)
- bits per oneOnOneBalance: 2 (0=undef, 1=team1, 2=team2, 3=tie)

bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
So if we use the optional teamIdx_B => space for 7 + 9 = total 16 teams per league.

The update of a team through a league always requires updating 4 uint256 of the league struct. 



# Main Functions to control Updaters/Challengers

## ComputeNewState(teamIdx) - can be called in view mode
- returns:
-- new state for every player in the team, after evolving through a league. 

1. Starts with the state of the team before starting the league
2. Computes all playGames for that team in the league (if some games were already updated, skips them)
3. Applies whatever evolution rules we have for each game
4. Returns final state of all players.

Notes: 
- if team is in challenging period, this function checks it, and considers the 'old' state as starting point.
- to know if a game is safely written:
    - first it checks if the entry in 'results' is nonzero. 
    If so, then checks if the opponent's team:
    - has a currentLeagueIdx that is different (which can only happen if challenging period was over)
    - if not, check if blockNumLastUpdate is old enough.


## WriteNewStateByUpdater(teamIdx, newPlayerStates)
- rewrites the 16 uint new player states after having called ComputeNewState in view mode, while swithcing stateBit for them (in the same uint256)
- rewrites 4 uint for league:
    - uint256 resultsFirstHalf
    - uint256 resultsSecondHalf
    - uint256 goalAverages
    - uint256 oneOnOneBalance
- writes blockNumLastUpdate inside uint256 playerIdxsB
- writes lastUpdaterAddr


## ChallengeUpdate(teamIdx, newPlayerStates)
- checks if team is in challenging period, otherwise halts.
- calls ComputeNewState via a Tx
- if newState coincides with what the last updater wrote, halts. Otherwise 
- writes the new player states (no need to change stateBit, as updater already did)
- sets wasProcessedByBC to 1. 
- no need to re-write blockNumLastUpdate, though it's free to do so (all in the same uint)
- no need to re-write lastUpdaterAddr, since it was the BC who did it
- gets the money from lastUpdaterAddr's deposit.



# Costs

Updating a team in a league rewrites:
- (16+4+1) uints = 21 uints
- 1 address

 5000 per word SSTORAGE => 105K  ==> expect 130K gas
 add 4000 for each if we need to access them (!), maybe we don't for most or any




# Extension to 1 update for any number of leagues
## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + blockNumLastUpdate ( + 28b = 224)

    - address lastUpdaterAddr = who wrote the last update

    - uint256 leagues[], where each league has:
            - leagueIdx 
            - wasProcessedByBC 









