pragma solidity ^0.4.24;

import "../ERC721/players.sol";

/**
 * @title PlayerFactoryMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract PlayerFactoryMock is PlayerFactory {
    function getRole_(uint idx, uint8 nDefenders, uint8 nMids) public pure returns(uint8) {
        return super.getRole(idx, nDefenders, nMids);
    }

    function createPlayerInternal_(string _playerName, uint _teamIdx, uint8 _playerNumberInTeam, uint _playerState)
        public
    {
        return super.createPlayerInternal(_playerName, _teamIdx, _playerNumberInTeam, _playerState);
    }

    function computePlayerStateFromRandom_(uint rndSeed, uint8 playerRole, uint currentTime)
        public
        pure
        returns(uint)
    {
        return super.computePlayerStateFromRandom(rndSeed, playerRole, currentTime);
    }

    function createBalancedPlayer_(
        string _playerName,
        uint _teamIdx,
        uint16 _userChoice,
        uint8 _playerNumberInTeam,
        uint8 _playerRole
    )
        public 
    {
        super.createBalancedPlayer(_playerName, _teamIdx, _userChoice, _playerNumberInTeam, _playerRole);
    }

    function createUnbalancedPlayer_(
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
        public 
    {
        return super.createUnbalancedPlayer(
        _playerName,
        _teamIdx,
        _playerNumberInTeam,
        _monthOfBirthAfterUnixEpoch,
        _defense,
        _speed,
        _pass,
        _shoot,
        _endurance,
        _role
        );
    }

    function getNCreatedPlayers_() public view returns(uint) { 
        return super.getNCreatedPlayers(); 
    }

    function getPlayerState_(uint playerIdx) public view returns(uint) {
        return super.getPlayerState(playerIdx);
    }

    function getPlayerName_(uint playerIdx) public view returns(string) {
        return super.getPlayerName(playerIdx);
    }
}