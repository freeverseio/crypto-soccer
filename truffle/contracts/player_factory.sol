pragma solidity ^ 0.4.24;

import "../node_modules/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./storage.sol";
import "./helper_functions.sol";

contract PlayerFactory is Storage, HelperFunctions {
    /// @dev Fired whenever a new player is created
    event PlayerCreation(string playerName, uint playerIdx, uint playerState);
    enum Role { Keeper, Defense, Midfield, Attack, Substitute, Retired }

    // @dev obtain player role from its position and strategy first-second-third (i.e. 4-3-3)
    function getRole(uint idx, uint8 first, uint8 second, uint8 /*third*/) public pure returns(uint8) {
        require (idx < kMaxPlayersInTeam);
        if (idx == 0)
            return uint8(Role.Keeper);
        else if (idx > 0 && idx <= first)
            return uint8(Role.Defense);
        else if (idx > first && idx < first+second+1)
            return uint8(Role.Midfield);
        else
            return uint8(Role.Attack);
    }

    /// @dev The main Player struct.
    /// @param name is a string, unique for every Player
    /// @param state is a uint256 that encodes age and skills. Each is a uint16 between 0-9999.
    ///     Note that there is room for up to 74 digits. The list is:
    ///         0-monthOfBirthAfterUnixEpoch; if this goes up to 9999, then the game will run for at least 800 more years.
    ///         1-defense
    ///         2-speed
    ///         3-pass
    ///         4-shoot (for a goalkeeper, this is interpreted as ability to block a shoot)
    ///         5-endurance
    ///         6-role
    // @param role is: 0=keeper, 1=defence, 2=midfielder, 3=attacker, 4=retired


    /// @dev An internal method that creates a new player and stores it. This
    ///  method doesn't do any checking and should only be called when the
    ///  input data is known to be valid.
    /// @param _playerState It contains all you need to know about a player
    function createPlayerInternal(string _playerName, uint _teamIdx, uint8 _playerNumberInTeam, uint _playerState)
        internal
    {
        // make sure this player name is unique. If so, it has never been assigned to a Team.
        // all teams, when created, have the first player Idx set to uint(-1), to signal
        // that the team has been created.
        bytes32 playerNameHash = keccak256(abi.encodePacked(_playerName));
        require(playerToTeam[playerNameHash].playersIdx == 0);

        // push payer, and update mapping and player count
        uint nCreatedPlayers = players.length;
        players.push(Player({name: _playerName, state: _playerState}));
        playerToTeam[playerNameHash] = teams[_teamIdx];
        teams[_teamIdx].playersIdx = setNumAtIndex(
            nCreatedPlayers,
            teams[_teamIdx].playersIdx,
            _playerNumberInTeam,
            20
        );

        // emit the creation event
        emit PlayerCreation(_playerName, nCreatedPlayers, _playerState);
    }

     function getDefaultPlayerState(
        Team _team,
        uint8 _playerNumberInTeam
     )
     internal
     pure
     returns (uint)
     {
        uint dna = uint(keccak256(abi.encodePacked(
            _team.name,
            _team.userChoice,
            _playerNumberInTeam
            )));
            return computePlayerStateFromRandom(dna, getRole(_playerNumberInTeam,4,3,3), _team.timeOfCreation);
     }

    function computePlayerStateFromRandom(uint longRnd, uint8 playerRole, uint currentTime)
        internal
        pure
        returns(uint)
    {
        // state[0] -> age, state[6] -> role
        // state[1]...state[5] -> skills
        // get random numbers between 0 and 9999:
        uint16[] memory states = decode(7, longRnd, 14);

        // First number is age, in years, at moment of creation can vary between 16 and 35.
        states[0] = 16 + (states[0] % 20);

        // Last number is role
        states[6] = playerRole;

        // Convert age to monthOfBirthAfterUnixEpoch
        // We can optimize by not declaring these as variables, and putting the exact numbers. I leave it this way for
        // clarity, for the time being.
        uint years2secs = 365 * 24 * 3600;
        uint month2secs = 30 * 24 * 3600;
        states[0] = uint16((currentTime - states[0] * years2secs) / month2secs);

        // the next 5 are skills. Adjust skills to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 sk = 1; sk < 6; sk++) {
            states[sk] = states[sk] % 50;
            excess += states[sk];
        }
        // at this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        excess = (250 - excess)/5;
        for (sk = 1; sk < 6; sk++) {
            states[sk] = states[sk] + excess;
        }

        return encode(7, states, 14);
    }

    /// @dev Creates a player where skills are set pseudo-randomly
    /// @param _teamIdx The idx of the team to which this player belongs
    /// @param _userChoice The user enters a team name, then chooses among many possible teams varying this number.
    /// the skills are determined from a random number, which is determined by concatenating the team's name with the
    /// player number. We will not allow two teams with the same name, and hence, same player numbers will not lead
    /// to the same skills. We can optimize this a bit by getting more 4-digit randoms from the long randoms we
    /// generate, so that we need to generate less
    /// @param _playerRole encodes the positions:
    ///         0=keeper, 1=defence, 2=midfield, 3=attack, 4=substitute, 5=retired.
    /// It returns the hash of the player's name
    function createRandomPlayer(
        string _playerName,
        uint _teamIdx,
        uint16 _userChoice,
        uint8 _playerNumberInTeam,
        uint8 _playerRole
        )
        public
    {
        require (_teamIdx < teams.length);
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

    function createPlayer(
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
        public {
        // we should make sure all numbers are below 1e5
        require (_teamIdx < teams.length);
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

    function getNCreatedPlayers() external view returns(uint) { return players.length;}

    function getPlayerState(uint playerIdx) external view returns(uint) {
        return players[playerIdx].state;
    }
}
