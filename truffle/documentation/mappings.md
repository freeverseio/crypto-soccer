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






