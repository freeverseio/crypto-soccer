# Explanation about mappings in this game.

Notation:  pname = playerName,  tname = teamName, idx = index

# Players
Player is a struct that has:
    - strin pname = player name, unique.
    - uint state = serialization of age and skills.

players[] is a vector of Players.

the genesis player (at pos 0) is has name: _ and state uint(-1)

# Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdx = serialization of all idx (positions in player[] array)

teams[] is a vector of Teams.

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

inputs: an array of teamIdx with nTeams elements, or a serialization of it. It could contain nTeams explicitly in the last value.

struct:
    - uint teamsIdxs: the serialization of teamIdx
    - n0: starting block
    - nStep: separation between blocks: n0, n0 + nStep, n0 + 2 nStep,...

Given Nteams, there are nGamesPerTeam = 2 (nTeams - 1) games to be played by each team.
In total, there are nTotalGames = nTeams (nTeams-1)

Proof: every team plays (n-1) times as local => n (n-1)
Another proof: there are 2 (n-1) rounds, each with n/2 games.

For 10 teams => 90 games. For 11 => 110 games. For 20 teams => 380 games. 

If we only keep the result of the games, we could add this:

    - uint games: serialization of game results

Say we use 4 bits for the score of each team (0,...,15). We could have 256/4=64 scores => 32 games.
Say we use 3 bits for the score of each team (0,...,7). We could have 256/3=85.3 scores.

If we only store 'who won', we need 2 bits: 0=not played, 1=team 1, 2= team 2, 3 = tie

We can then store 256/2 = 128 games.

We need to sort them. In every round r=0,...,n-2 there are n=1,...,nTeams/2 games. 

Use Round-Robin algorithm for tournament scheduling (and modify from Wikpedia so that they increase)

Let N = nTeams.

    r=0
        0 N-1 N-2   ,...,N/2+1
        1 2   3     ,...,N/2

    r=1
        0 1 N-1 ... N/2+2
        2 3 4   ... N/2+1

    r=N-2
        0   N-2  N-3 ... N/2
        N-1 1    2   ... N/2-1

the number appearing in any position always increases by 1, and jumps from n-1 to 1, not zero. Define:
    P(x)=x if x < n, x-(n-1) otherwise

r: 
    g0          = ( 0,         P(1+r) )     = (0, 1), (0,2),  ...,(0,N-1)
    g1          = ( P(N-1+r),  P(2+r) )     = (N-1, 2), (1,3),... (N-2, 1) 
    g2          = ( P(N-2+r),  P(3+r) )     = (N-2, 3), (N-1,4),...(N-3, 2) 

    gn          = ( P(N-n+r),  P(n+1+r) )

    g_{N/2-1}   = ( P(N/2+1+r),P(N/2+r))    = (N/2+1,N/2),...


so if g=1,...,N (N-1) is the index that goes through every game in the league, then, the relation with the round and game in that round is:

If g <= N (N-1)/2: We are in the first leg of the league.

    g = r * N/2 + n,   n = 0,...,N/2-1;   r = 0,...,N-1

whose inverse is:

    r = floor( g / (N/2) ) 
    n = g - r * N/2

from which we obtain the two teams that play there: 
    
    P(N-n+r) vs P(n+1+r), unless n=0, in which case the first team is team 0.

If g > N (N-1)/2: We are in the second leg of the league. Just change g <- g mod N (N-1)/2, apply the previous formulas, and just reverse the final order of teams.

















# ideas

- A league could be filled with teams with all-players-with-all-skills-identical = Difficulty. 
-- if Difficulty = 250/5, then it's an avg team. League creators could set D above, if needed.
-- there's no need to actually 'create' those players nor store them anywhere.


