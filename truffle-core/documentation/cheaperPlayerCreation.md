Recall the main structs:

## Players

- string pname = player name, unique.
- uint state = serialization of skills + currentTeamIdx

skills = function( teamName, userChoice, positionInTeam)

players[]

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





1. Imagine a user joins and provides:  teamName, userChoice. 

2. We add to teams[end+1] a team whose only initialized struct is the tname+userChoice concatenated or not. We also write teamToOwnerAddr. As byproduct, forbids 2 equal teamNames.

At this point, both updaters or challenger can compute (not read) the playerStates. They don't even need to access any written player struct at all. For each player in the team, they look at the team.playersIdxA, see that ==0, and hence, deduce this playerState needs to be computed from 
	- hash(tname + userChoice + 0,...10)


Never write players (!!!)

playerIdx always fixed at birth = teanNumer * numPlayersPerTeam + positionInTeam

When selling:

**** mapping playerIdx -> team.  (NO CAL)

When quering team for a given playerIdx:
	- exists in mapping? 
		- if No => formula
		- if yes => mapping





4. When joining a league, as always, we add their teamIdx and write 'currentLeagueIdx'


5. Updaters/Challengers, instead of checking that the hash coincides with prevLeagueMerkleRoot, they need to compute it.

6. The team will then have a final state, hidden in the final MerkleRoot. It can continue playing leagues.



Summary:
	- teamCreation: write 2 uints (teamName+choice) + teamToOwnerAddr
	- playerCreation (needed for selling): write 3 uints (name, state, playerToTeam)


-------------

Further savings: if we don't allow playerName to be defined by user, we can save 2 uints per player (lot of money). Names could be generated automatically by using, again, hash(tname + userChoice + 0,...10), and retaining only, say, 100bits, then storing it in the remaining part of playerState.










----

SELLING

Imagine someone wants to sell a player now. There are 2 main cases. 

1. If player has never evolved or joined a league. Maybe because the name or skills are cool. 
The user himself can create the TX (also an updaterm of course).

	- TX put for sale: update(player 7, teamIdx=34),
		- writes the playerState=uint(skills + currentTeamIdx) by computing it on the fly from teamName+Userchoice+...

	- TX: the sell/buy TX:
		- updates the new team's playerIdx
		- update the entry in the old team playerIdx (all of them) the same way we would. Note that we don't update the others, so that we know that those players 



2. If player has evolved or playerd leagues. An updater MUST first update the state:

	- TX put for sale:  => emit even "please update me!" 

	- TX update by updater (provides state and the TX checks that the hash is correct)
		- writes the playerState=uint(skills + currentTeamIdx) by providing his last computed state.

	- TX: the sell/buy TX:
		- updates the new team's playerIdx
		- update the entry in the old team playerIdx (all of them) the same way we would. Note that we don't update the others, so that we know that those players 



TODOs:
- finish this file (player Creation) with ...
- update scaling2 now that we don't have playerNames!!!! no mapping....
- structs, logica become updater/challenger, freeze while participating....
- refactoring playGame 'in a txt file' now that it will use data from outside.
- introduce tactics (442) between games (see scaling.md)
- API for: think about evolvePlayer or evolveTeam
	- when playing a game
	- when time passes:
		- age effect?
		- decide training to focus on skills: share 10 points among the 5 skills, or a max of 2-3.
		- or maybe playing a game gives
			- instant reward (skills++) together with extra training points.
















