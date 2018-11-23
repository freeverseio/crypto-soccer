Recall the main structs:

## Players

- string pname = player name, unique.
- uint state = serialization of skills + currentTeamIdx

state = function( teamName, userChoice, positionInTeam)

players[]

## Teams

- string tname = unique
- uint256 playersIdxA = serialization of 9 playerIdx (positions in player[] array)
- uint256 playersIdxB = serialization of up to 7 playerIdx + previousLeagueIdx +currentLeagueIdx 

## Mappings

- mapping(bytes32 => Team) playerToTeam; /// from hash(playerName) to a Team struct.
- mapping(bytes32 => address) public teamToOwnerAddr; /// from team hash(name) to the owner's address.

## Uniqueness

Currently we require:
- require( playerToTeam[pNameHash].playersIdx == 0);
- require( teamToOwner[tNameHash]==0 )

We propose to leave the playerName totally out of the game.

See ERC721 spec here: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md

Note that the (optional) function in ERC721 is:

    function tokenURI(uint256 _tokenId) external view returns (string);

which points to a website/resource that returns a JSON struct:
{
    "title": "Asset Metadata",
    "type": "object",
    "properties": {
        "name": {
            "type": "string",
            "description": "Identifies the asset to which this NFT represents",
        },
        "description": {
            "type": "string",
            "description": "Describes the asset to which this NFT represents",
        },
        "image": {
            "type": "string",
            "description": "A URI pointing to a resource with mime type image/* representing the asset to which this NFT represents. Consider making any images at a width between 320 and 1080 pixels and aspect ratio between 1.91:1 and 4:5 inclusive.",
        }
    }
}

Alessandro proposes a nice pattern:

    string private _tokenCID;
    constructor( string CID) public { _tokenCID = CID; }
    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        uint256 state = getState(tokenId);
        string memory stateString = uint2str(state);
        return strConcat(_tokenCID, "/?state=", stateString);
    }

So that it points to a server that GENERATES name, description, and the image, all on the fly, given the state of the player.


# Proposal

Get rid of playerName, and hence get rid of mapping(bytes32 => Team) playerToTeam;

Split all players into:
	A.- players created via 'team creation'
	B.- players created uniquel (promoplayers, superplayers, etc)

Let's discuss only the subset A. 

We propose that the creation of a team automatically generates nPlayersPerTeam players... VIRTUALLY. They are not really written anywhere, buy their skills are uniquely determine by the team created:

	state = function( teamName, userChoice, positionInTeam)

So there is an invertible relation  (teamIdx, positionInTeam)  <--> playerIdx

	playerIdx = teamIdx * numPlayersPerTeam + positionInTeam

Such playerIdx is always fixed, never changed.



Usecase workflow detailed:

1. Imagine a user joins and provides:  teamName, userChoice. 

2. We write this:
	- create a team struct with only the concat(tname+userChoice) in it. 
	- append to teams[]
	- write teamToOwnerAddr (as byproduct, forbids 2 equal teamNames)

Note that we leave team.playerIdxs as zero. 

Note that, at this point, both updaters or challengers can compute the playerStates. They don't need to access any written player struct at all. For each player in the team, they look at the team.playersIdxA, see that ==0, and hence, deduce this playerState needs to be computed from 
	- hash(tname + userChoice + 0,...10)


3. When joining a league, as always, we add their teamIdx and write 'currentLeagueIdx'


4. Updaters/Challengers, instead of checking that the hash coincides with prevLeagueMerkleRoot, they seem the 0 at the corresponding player in team.playerState, and compute it on the fly.

5. The team will then have a final state, hidden in the final MerkleRoot. It can continue playing leagues.



Summary:
	- teamCreation: write 2 uints (teamName+choice) + teamToOwnerAddr
	- an append to teams[]


# Selling/buying

Any operation that requires knowing the supersafe state of a player (e.g. sell/buy) requires someone to finally write the state. There are 2 main cases (the first one is possibly a rare one).

We will need a mapping:

- mapping(uint => Player) playerIdx2Player: a mapping from playerIdx given at birth of a team, to the entry of the vector of actually created playerIdxs. 

Entries in this mapping are never changed.


1. If player has never evolved or joined a league (maybe the value is because the name or skills are cool). 
	- TODO: how do we verify this? 

The user himself can create the TX (could also be done via asking an updater, as with case 2).

	- TX put for sale: update(player 7, teamIdx=34),
		- writes the playerState=uint(skills + currentTeamIdx) by computing it on the fly from teamName+Userchoice+...
		- writes this mapping(uint => Player) playerIdx2Player:


	- TX: the sell/buy TX:
		- updates the new team's playerIdx
		- update the entry for thisplayer in the old team playerIdx.



2. If player has evolved or playerd leagues. An updater MUST first update the state:

	- TX "put for sale": emit even "please update playerIdx!" 

	- TX update by updater: provides state and the TX checks it is correct:

		- if playerIdx2Player[playerIdx] == 0, 
			- it computes it, it finds team via formula, and computes state from team data
			- otherwise, it computes team from playerState, and checks as always, that hash if contained in prevMerkleRoot (updater provides MerkleProof for this player only)

		- writes the playerState=uint(skills + currentTeamIdx) by providing his last computed state.


	- TX: the sell/buy TX: (as in case 1 above)





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











