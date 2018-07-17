pragma solidity ^ 0.4.24;

/// @dev Defines all storage structures and mappings for cryptosoccer
// TODO: this contract should be protected by Ownable
contract Storage {

    struct Player {
        string name;
        uint state;
    }

    uint8 constant kMaxPlayersInTeam = 11;

    struct Team {
        string name;
        uint256 playersIdx;
        uint256 timeOfCreation;
        uint16 userChoice;
    }

    /// @dev An array containing the Team struct for all teams in existence. The ID
    ///  of each team is actually an index into this array.
    Team[] teams;

    /// @dev An array containing the Player struct for all players in existence. The ID
    ///  of each player is actually an index into this array.
    Player[] players;

    /// @dev A mapping from hash(playerName) to a Team struct.
    mapping(bytes32 => Team) public playerToTeam;

    /// @dev A mapping from team hash(name) to the owner's address.
    mapping(bytes32 => address) public teamToOwnerAddr;

    constructor() public {
        // create the deault player at pos 0.
        players.push(Player({name: '_', state: uint(-1)}));
    }
}
