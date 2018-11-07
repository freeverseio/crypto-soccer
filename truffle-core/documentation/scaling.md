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
Indeed, for 16teams, we still have 256-16x9 = 112bits left, which could store 56 extra games.



# Players

Player is a struct that has:
    - string pname = player name, unique.
    - uint state = serialization of skills. Currently we also store birthMonth, and role. 

We propose that role is taken out and made part of the league struct (as explained in another file). 

Assuming we keep birthMonth, how many states can we keep?

12bit => each state = 5x12= 60 => 4 states, since 12 + 4 x (5x12) = 252 bit ==> still 4 left!

(11 does not allow going to 5)

Let's keep 12bit => 4 states, plus 4 left.


# Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdx = serialization of all idx (positions in player[] array)


# League

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


- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)
- bits per goalAvg: 9 (max -256...256)
- bits per result: 2 (0=undef, 1=home, 2=away, 3=tie)

bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
So if we use the optional teamIdx_B => space for 7 + 9 = total 16 teams per league.




## Writing results

If we choose just 2 bits (0=not played, 1=team 1, 2= team 2, 3 = tie) we can have:
    - uint resultsFirstHalf;
    - uint resultsSecondHalf;
each storing 256/2 = 128 results

The position to write game g of round r is: 
    (recall nRounds = 2 (nTeams-1), nGamesPerRound = nTeams/2)

    pos(r,g) = r*nTeams/2 + g
which has a max of 2 (nTeams-1) nTeams/2 + nTeams/2 -1 = 2 nTeams nTeams/2 = nTeams^2 -1 = N (N-1)    




## Scheduling

We need to sort the games. In every round r=0,...,nTeams-2, there are n=1,...,nTeams/2 games.

Use a Round-Robin algorithm for tournament scheduling (modify Wikpedia's so that all entries increase instead of decrease).

Let us, for shortness, use N = nTeams.

    r=0
        0 N-1 N-2   ,...,N/2+1
        1 2   3     ,...,N/2

    r=1
        0 1 N-1 ... N/2+2
        2 3 4   ... N/2+1
.
.
.

    r=N-2
        0   N-2  N-3 ... N/2
        N-1 1    2   ... N/2-1

With our choice, the numbers appearing in any position always increase by 1, and jump from n-1 to 1, not zero. Define:

    P(x) = { x if x < n; x-(n-1) otherwise }

Then, for a given round r:

    game(0,r)          = ( 0,         P(1+r) )     = (0, 1), (0,2),  ..., (0, N-1)
    game(1,r)          = ( P(N-1+r),  P(2+r) )     = (N-1, 2), (1,3), ..., (N-2, 1)
    game(2,r)          = ( P(N-2+r),  P(3+r) )     = (N-2, 3), (N-1,4), ..., (N-3, 2)
    ...
    game(N/2-1,r)      = ( P(N/2+1+r),P(N/2+r))    = (N/2+1,N/2),...

So the generic formula, for a given round r and game n, is:

    game(n,r)          = ( P(N-n+r),  P(n+1+r) )



so if g=1,...,N (N-1) is the index that goes through every game in the league, then, the relation with the round and game in that round is:

If g <= N (N-1)/2: We are in the first leg of the league.

    g(r,n) = r * N/2 + n,   n = 0,...,N/2-1;   r = 0,...,N-1

whose inverse is:

    r = floor( g / (N/2) )
    n = g - r * N/2

from which we obtain the two teams that played at game g:

    - find r,n given g,
    - teams: P(N-n+r) vs P(n+1+r), unless n=0, in which case the first team is team 0.

If g > N (N-1)/2: We are in the second leg of the league. Just change g by { g mod N (N-1)/2 }, apply the previous formulas, and just reverse the final order of teams.


## Finding games for a given team

An important function we'll need to implement is one that, given a team, it finds the games in which they play.

For t=0, we have n=0 for all r. So the games are g(0,r).
Let's look at t>0. Since:

    game(r,n)          = ( P(N-n+r),  P(n+1+r) )


We ask: for fixed t = 1,...,N-1, (we exclude the first team!!) when is it the case that either

    P(N-n+r) = t, or P(n+1+r) = t ?

There must be exactly one solution to each of those, since the 1st answers when is the team at home, and the second, when is the team away.

We parametrize the solutions by the round 'r':

    0 <= r <= N-2,

so that given t and r, in which game 'n' of that round do they appear? Note that we need answers for which:

    Restrictions:    0 <= n < N/2 (up to N/2-1)

Since P(x) = t has two solutions: x=t, x=t+N-1, then we have, for fixed t,r:

    A1:    n = N + r - t
    A2:    n = 1 + r - t

    B1:    n = -1 -r + t ( note n(B1) = -b(A2) )
    B2:    n = N -2 -r + t

It can be seen that the restrictions imply:

    A1: t in [t-N, t-N/2-1]
    B1: t in [t-N/2, t-1]
    A2: t in [t-1, t+N/2-2]
    B2: t in [t+N/2-1, t+N-2]

They all form a continuous non-intersecting set, except for the transition A2-B1, which intersects at r=t-1, for which n=0 in both. This is impossible for A2, since the home team at n=0 is always t=0. So:

    A2: t in [t, t+N/2-2]

So for a given t, determine the 3 points (t-N/2-1, t-1, t+N/2-2) and use A1,B1,A2,B2 to find n as r goes from 0 through those points. Finally, use g(r,n) above to find the game.

We shall call this procedure finding g(t,r) => "find the game for team t at round r"




# Updating player skills for Oracles

Instead of keeping all player states + timestamps, just keep 2 states in the same uint:

    - playerState = serialize(state 1, state 0, bool current state)

Easy since, recall, we only need 70 bits for a state.

For example: at start, current state = 0, which points to the initial created player state. 

    - playerState = serialize( rubbish state, state 0, 0)

An Oracle re-writes:

    - playerState = serialize(new state 1, state 0, 1)

During the challenging period, the challenger has access to state 0, and can compute new state 1.

This applies to players. The proposal is to always update an entire team together. This implies that in a league, we don't need to know about that team anymore (all games with that team have been processed), and hence, the team can move on to new states and new leagues.

This also means that during the challenging period the team cannot SELL/BUY players (it can SELL/BUY the entire team), so that the previous playerIdxs of the team is the one at league start.

playerIdxs0 --> ORACLE + CHALLENGE + pause time --> new player states --> BUY/SELL --> new playerIdxs1

We must be able to answer: isPlayerVerified? or equivalently, is it in pause time?

This can be done by only keeping the block num of the last oracle write. There is plenty of room in player state. Although that would involve writing the same number in the 11 player state. It is more a 'team' property. 

We need to add it to the team struct:

    - uint lastBlockUpdate

Too bad that we basically 24-28 bits for it. Maybe the space can be reused for other things.


# Changing strategy (443, 541, etc)

We shall make further use of the space left in playerState ( 256 - 2 x 70 + 1 = 115 bit). If we limit to a max of 20 games per league, this leaves space for 5 bit chunks. So a player role may now be:
    0 - undefined
    1 - keeper
    2 - defence
    3 - mid
    4 - attack
    5... others (e.g. retired)

So a when a player joins a league all such 20 chunks are set to undefined except for the last one. 
Say it is a midfielder. So the 20 chunks are 0...00002
If a user updates a player's role, say, at game five, to attack. It sets the chunks to 0....00322222.
So, the Oracle/Challenger knows the role in every game. As always, this remains untouched until a new league starts, which is beyond the challenging period.




## Details of algorithm

Note that players can be restricted to be traded only when NOT in a league. 

Inputs: 
    - playerIdx 


Step 1: find current team for that player via mapping(bytes32 => Team) playerToTeam;

Step 2: find league for that team

    - we need a mapping(bytes32) => League
        - doubt: to League or to leagueIdx? Is it redundant with league.teamIdx?
        - the latter, I don't think so, since we need to scan the teams in a league...

Step 3: Goes through all games where the team must play:
    - Read the players skills at start of league
    - If game has not been processed yet:
        - find players' role in that game
        - play game
        - write result (unless called in view mode)
    - Calculate deam deltas

Step 4: write final player states (unless called in view mode)











## Initial numbers (without time-stamping)
    - adding a state to 1 player = 47882
    - adding a state to 11 players = 310582,
        - note that 11 x 47882 = 526K, 
        - so this is more: 11 x 27882 + 20K = 326K.

    => updating all team is ~ the same as playing 1 game





## Costs for us of a typical user

User starts on our site with one of the many existing-teams-pool.
Can we have such a pool?

- all view functions are free, so we can: query team skills for selling or playing a game
- problem: as soon as they change a name, or tactics we're doomed.
    - plus, UX is bad if you can't personalize a but the team.

Say we pay for the team creation --> currently 11x 

creates empty team: 78109
adding 1 player: 131477
adding 11 players: 1.446.247
creating a team: 1.541.112 = 1446247 + 78109 + 20.000
If we did that atomically in 1 go: 1541112 - 12 * 20000 = 1.301.112 (saves 120M)

we are talking at 2.5 USD for 1.3M gas...

Alternative: play free, only allow change names, etc. after 1st league, first pay.

==> Creating a team is about 2.5 USD, playing more leagues is about 50 cents per team.




## Can we break computations in steps?

To simplify, stay within 1 league. Say it has maaany teams, so updating 1 team is too costly.

Recall our paradigm:

playerIdxs0 --> ORACLE + CHALLENGE + pause time --> new player states --> BUY/SELL --> new playerIdxs1

Now, imagine that the ORACLE result is too costly. The update is for N=40 games. The max is 30.

updateTeamForGames(40) =  updateTeamForGames(10) · updateTeamForGames(30)

so his update(40) first does the update(30), computes and writes hash(state(30)), then computes update(10) and writes:
    - final state
    - hash(state(30))

Any challenger can provide the state(30) whos hash coincides, and prove he's wrong.


## Can we nest 2 leagues without updating state?

- If leagues require validation, nope.
- If leagues don't require...

    









# ideas

- A league could be filled with teams with all-players-with-all-skills-identical = Difficulty.
    - if Difficulty = 250/5, then it's an avg team. League creators could set D above, if needed.
    - there's no need to actually 'create' those players nor store them anywhere.


