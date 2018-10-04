pragma solidity ^ 0.4.24;

// TODO: import "../node_modules/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./storage.sol";
import "./helpers.sol";

/*
    Contract to manage player creation
*/

contract PlayerFactory is Storage, HelperFunctions {

    /// @dev Event fired whenever a new player is created
    event PlayerCreation(string playerName, uint playerIdx, uint playerState);

    enum Role { Keeper, Defense, Midfield, Attack, Substitute, Retired }

    /// @dev Returns player role given his pos in the team, and a selected strategy 
    /// @dev Strategy (e.g. 4-4-3) is specified by the first 2 nums. 
    /// @dev The 3rd number is not needed (it always equals 10 - nDefenders - nMids)
    function getRole(uint idx, uint8 nDefenders, uint8 nMids) internal pure returns(uint8) {
        require (idx < kMaxPlayersInTeam, "Player pos in team larger than 11!");
        if (idx == 0)
            return uint8(Role.Keeper);
        else if (idx > 0 && idx <= nDefenders)
            return uint8(Role.Defense);
        else if (idx > nDefenders && idx < nDefenders+nMids+1)
            return uint8(Role.Midfield);
        else
            return uint8(Role.Attack);
    }

    /// @dev An internal method that creates a new player and stores it. This
    /// @dev method doesn't do any checking and should only be called when the
    /// @dev input data is known to be valid.
    function createPlayerInternal(string _playerName, uint _teamIdx, uint8 _playerNumberInTeam, uint _playerState)
        internal
    {
        /// @dev First, make sure this player name is unique. If so, it has never been assigned to a Team.
        /// @dev A team is created if it has a not-null owner addr.
        bytes32 playerNameHash = keccak256(abi.encodePacked(_playerName));
        bytes32 teamNameHash = keccak256(abi.encodePacked(playerToTeam[playerNameHash].name));
        require(teamToOwnerAddr[teamNameHash] == 0, "Player already exists with this name");

        /// @dev Update player count
        uint nCreatedPlayers = players.length;

        /// @dev Push playert
        players.push(Player({name: _playerName, state: _playerState}));

        /// @dev Update mapping
        playerToTeam[playerNameHash] = teams[_teamIdx];

        /// @dev Update inverse relation (from teams to playerIdx)
        teams[_teamIdx].playersIdx = setNumAtIndex(
            nCreatedPlayers,
            teams[_teamIdx].playersIdx,
            _playerNumberInTeam,
            20
        );

        /// @dev Emit the creation event
        emit PlayerCreation(_playerName, nCreatedPlayers, _playerState);
    }

    /// @dev Main interface to create a player by users. We receive a random number,
    /// @dev computed elsewhere (e.g. from hash(name+userChoice+dorsal)) and create 
    /// @dev a balanced player whose skills add up to 250.
    function computePlayerStateFromRandom(uint rndSeed, uint8 playerRole, uint currentTime)
        internal
        pure
        returns(uint)
    {
        /// @dev Get random numbers between 0 and 9999 and assign them to states, where:
        /// @dev state[0] -> age, state[6] -> role
        /// @dev state[1]...state[5] -> skills
        uint16[] memory states = decode(7, rndSeed, 14);

        /// @dev Last number is role, as provided from outside. Just store it.
        states[6] = playerRole;

        /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
        states[0] = 16 + (states[0] % 20);

        /// @dev Convert age to monthOfBirthAfterUnixEpoch.
        /// @dev TODO: We can optimize by not declaring these as variables, and putting the exact numbers. 
        /// @dev I leave it this way for clarity, for the time being.
        uint years2secs = 365 * 24 * 3600;
        uint month2secs = 30 * 24 * 3600;
        states[0] = uint16((currentTime - states[0] * years2secs) / month2secs);

        /// @dev The next 5 are states skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 sk = 1; sk < 6; sk++) {
            states[sk] = states[sk] % 50;
            excess += states[sk];
        }
        /// @dev At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        excess = (250 - excess)/5;
        for (sk = 1; sk < 6; sk++) {
            states[sk] = states[sk] + excess;
        }

        return serialize(7, states, 14);
    }

    /// @dev Creates a player where skills are set pseudo-randomly assigned
    /// @param _teamIdx The idx of the team to which this player belongs
    /// @param _userChoice The user enters a team name, then chooses among many possible teams varying this number.
    /// the skills are determined from a random number, which is determined by concatenating the team's name with the
    /// player number. We will not allow two teams with the same name, and hence, same player numbers will not lead
    /// to the same skills. We can optimize this a bit by getting more 4-digit randoms from the long randoms we
    /// generate, so that we need to generate less.
    /// @param _playerRole serializes the positions:
    ///         0=keeper, 1=defence, 2=midfield, 3=attack, 4=substitute, 5=retired.
    /// @dev Returns the hash of the player's name
    function createBalancedPlayer(
        string _playerName,
        uint _teamIdx,
        uint16 _userChoice,
        uint8 _playerNumberInTeam,
        uint8 _playerRole
    )
        internal 
    {
        require (_teamIdx < teams.length, "Trying to assign a player to a team not created yet");
        uint dna = uint(keccak256(abi.encodePacked(
            teams[_teamIdx].name,
            _userChoice,
            _playerNumberInTeam
        )));
        createPlayerInternal(
            _playerName,
            _teamIdx,
            _playerNumberInTeam,
            computePlayerStateFromRandom(dna, _playerRole, now)
        );
    }

    /// @dev Creates a player where skills are set outside the blockchain, and hence, can be arbitrary
    /// @dev To be used, eventually, to generate promo players or super players. 
    function createUnbalancedPlayer(
        string _playerName,
        uint _teamIdx,
        uint8 _playerNumberInTeam,
        uint _monthOfBirthAfterUnixEpoch,
        uint _defense,
        uint _speed,
        uint _pass,
        uint _shoot,
        uint _endurance,
        uint _role
    )
        internal 
    {
        /// @dev TODO: we should make sure all numbers are below 2^14-1
        require (_teamIdx < teams.length, "Trying to assign a player to a team not created yet");
        uint bits = 14;
        uint state = _monthOfBirthAfterUnixEpoch +
                     (_defense     << bits) +
                     (_speed       << (bits*2)) +
                     (_pass        << (bits*3)) +
                     (_shoot       << (bits*4)) +
                     (_endurance   << (bits*5)) +
                     (_role        << (bits*6));

        createPlayerInternal(_playerName, _teamIdx, _playerNumberInTeam, state);
    }


/* 
    @dev Section with functions only for external/testing use.
*/
    function getNCreatedPlayers() internal view returns(uint) { return players.length; }

    function getPlayerState(uint playerIdx) internal view returns(uint) {
        return players[playerIdx].state;
    }

    function getPlayerName(uint playerIdx) internal view returns(string) {
        return players[playerIdx].name;
    }
}
