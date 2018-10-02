# Relation between block and time

About 1 block is produced every 10s. 
    - about 3M blocks per year: we would need 22 bit
    - about 8K blocks per day: we would need 13 bit
    - if we want the game to last 100 years => 29bit
    - if we want the time between games to be 1 week => 16 bit



# Explanation about mappings in this game.

Notation:  pname = playerName,  tname = teamName, idx = index

# Players
Player is a struct that has:
    - string pname = player name, unique.
    - uint state = serialization of age and skills.

We currently use 14bit for each skill (max skill value = 16K). We don't need to put age nor role in the state, necessarily. So skills take about nSkills x 14 = 70 (or 84 if we include age and role). There's plenty of space left. We'll use it later to store 2 state in the same uint256.

    - players[] is a vector of Players.

the genesis player (at pos 0) has name: _ and state uint(-1)

# Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdx = serialization of all idx (positions in player[] array)


    - teams[] is a vector of Teams.

# Mappings

    - mapping(bytes32 => Team) playerToTeam; /// from hash(playerName) to a Team struct.
    - mapping(bytes32 => address) public teamToOwnerAddr; /// from team hash(name) to the owner's address.

# Uniqueness

We ensure uniqueness of pname and tname by simply:

    - require( playerToTeam[pNameHash].playersIdx == 0);
    - require( teamToOwner[tNameHash]==0 )


# Quering over all players in a team

An example: in playGame, all players in a team are accessed via:

    input to playGame = teamIdx => getSkill(teamIdx, p in 0,...10)
        =>  gets team from teams[teamIdx],
            and playerIdx in array from getNumAtIndex(team.playersIdx, p, 20);


# League

## Struct

inputs: 
    - an array of teamIdx with nTeams elements, or a serialization of it. It could contain nTeams explicitly in the last value; 
    - the starting block
    - the number of blocks to wait between consecutive games

struct:
    - uint teamsIdxs: the serialization of teamIdx
    - serialization of (n0, nStep): 
        - starting block of first match to be played
        - nStep: separation between blocks: n0, n0 + nStep, n0 + 2 nStep,...

In a league of max 20 teams, they could have a max of 12 bit, not enough.
In a league of max 10 teams, 25 bit per team => 33M teams.

We need about 29 bit for block number, and about 17 for nStep.
This leaves about 209 bit for other teams if we needed it. The nTeams may be limited by number of games below.


    - leagues[] is an array of leagues.


Given nTeams, there are nGamesPerTeam = 2 (nTeams - 1) games to be played by each team.
In total, there are nTotalGames = nTeams (nTeams-1) in a league.

Proof: every team plays (n-1) times as local => n (n-1)
Another proof: there are 2 (n-1) rounds, each with n/2 games.

For 10 teams => 90 games. For 11 => 110 games. For 20 teams => 380 games.

If we only keep the result of the games, we could add this:

    - uint games: serialization of game results

Say we use 4 bits for the score of each team (0,...,15). We could have 256/4=64 scores => 32 games.
Say we use 3 bits for the score of each team (0,...,7). We could have 256/3=85.3 scores.

If we only store 'who won', we need 2 bits: 0=not played, 1=team 1, 2= team 2, 3 = tie

We can then store 256/2 = 128 games.

We could even just use 1 bit: "has it been processed?" but maybe it's a pain, given that one needs to look back at the states before...?


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
    - adding a state to 11 players = 310582  (note that 11 x 47882 = 526K)
    => updating all team is ~ the same as playing 1 game











# ideas

- A league could be filled with teams with all-players-with-all-skills-identical = Difficulty.
    - if Difficulty = 250/5, then it's an avg team. League creators could set D above, if needed.
    - there's no need to actually 'create' those players nor store them anywhere.


