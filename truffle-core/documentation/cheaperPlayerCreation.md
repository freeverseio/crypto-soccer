Recall the main structs:

## Players

- string pname = player name, unique.
- uint state = serialization of skills + currentTeamIdx + prevTeamIdx

skills = function( teamName, userChoice, positionInTeam)

## Teams

- string tname = unique
- uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
- uint256 playersIdxB = serialization of up to 7 playerIdx + previousLeagueIdx +currentLeagueIdx 

## Mappings

- mapping(bytes32 => Team) playerToTeam; /// from hash(playerName) to a Team struct.
- mapping(bytes32 => address) public teamToOwnerAddr; /// from team hash(name) to the owner's address.

## Uniqueness

- require( playerToTeam[pNameHash].playersIdx == 0);
- require( teamToOwner[tNameHash]==0 )

------------

I think that we can get rid completely of playerToTeam, since we already use currentTeamIdx. 
Although, how do we check that name does not exist...? If we need a mapping, maybe map it to something else. Or maybe use it only for players which were named by user.

What if all starting players are assigned implicitly to consecutive teams.

0...10 => team1, ...

So for the first league, we do not write playerToTeam.

1. Imagine a user joins and provides:  tname, userChoice. 

2. We add to teams[end+1] a team whose only initialized struct is the tname+userChoice concatenated or not. We also write teamToOwnerAddr. As byproduct, forbids 2 equal teamNames.

At this point, both updaters or challenger can compute (not read) the playerStates. They don't even need to access any written player struct at all. For each player in the team, they look at the team.playersIdxA, see that ==0, and hence, deduce this playerState needs to be computed from 
	- hash(tname + userChoice + 0,...10)

3. (optional) 
	Imagine someone wants to sell a player now. User (or updater!) sends a 
	- TX: update(player 7, teamIdx=34, name="Toni"),
		- writes the playerStruct
		- writes the mapping playerToTeam (forbidding existing names)

	Then, the sell/buy TX:
		- updates the new team's playerIdx
		- update the entry in the old team playerIdx (all of them) the same way we would. Note that we don't update the others, so that we know that those players 


4. When joining a league, as always, we add their teamIdx and write 'currentLeagueIdx'


5. Updaters/Challengers, instead of checking that the hash coincides with prevLeagueMerkleRoot, they need to compute it.

6. The team will then have a final state, hidden in the final MerkleRoot. It can continue playing leagues.



Summary:
	- teamCreation: write 2 uints (teamName+choice) + teamToOwnerAddr
	- playerCreation (needed for selling): write 3 uints (name, state, playerToTeam)


-------------

Further savings: if we don't allow playerName to be defined by user, we can save 2 uints per player (lot of money). Names could be generated automatically by using, again, hash(tname + userChoice + 0,...10), and retaining only, say, 100bits, then storing it in the remaining part of playerState.
























