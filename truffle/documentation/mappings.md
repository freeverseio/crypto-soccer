# Explanation about mappings in this game.

Notation:  pname = playerName,  tname = teamName, idx = index

# Players
Player is a struct that has:
    - pname = player name, unique.
    - dna = hash of (tname+userNumber+numberInsideTeam), also unique. Determines the player's skills at moment of creation.

players[] is a vector of Players.

We could keep these two mapping only:

        playerToTeam:   hash(pname) -> Team

We would ensure uniqueness of pname and dna by simply:
        require( playerToTeam[hash(pname)]==0 ) 
    ...before creation.
We would not store anywhere the playerIdx, which is, the position in the player[] vector. 
We would still be able to query for all existing players, by simply going through all entries in the player vector.


# Teams
Team is a struct that has:
    - tname = unique
    - vector of player idx (positions in the players[] array)

teams[] is a vector of Teams.

We could keep these mappings only:
    - teamToOwner:  hash(tname) -> ownerAddr

We would ensure uniqueness of tname simply:
        require( teamToOwner[tname]==0 ) 

We would access all players in a teams by:
    dnaToPlayer[ team.dnas[3] ]



# Gas estimates
 ✓ creates an empty team, checks that nTeams moves from 0 to 1 (194335 gas)
 ✓ adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1 (254011 gas)
  
 Now: using playersIdx:
✓ creates an empty team, checks that nTeams moves from 0 to 1 (127866 gas)
✓ adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1 (186448 gas)

Now: without nCreatedPlayers:
✓ creates an empty team, checks that nTeams moves from 0 to 1 (107637 gas)
✓ adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1 (166026 gas)

Storing player's role in playerState
✓ creates an empty team, checks that nTeams moves from 0 to 1 (107637 gas)
✓ adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1 (146937 gas)
✓ plays a game using a transation, not a call, to compute gas cost (259410 gas)

playgame:
  without endurance:    261232
  with endurance:       272633



with uint16: plays a game using a transation, not a call, to compute gas cost (263144 gas)
with uint: 261824

